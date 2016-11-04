package web

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/wattapp/superchargers/pkg/metrics"
)

var (
	Schema           graphql.Schema
	isDyno           = os.Getenv("DYNO") != ""
	encryptChallenge = os.Getenv("LETS_ENCRYPT_CHALLENGE")
	encryptKey       = os.Getenv("LETS_ENCRYPT_KEY")
)

func Run() error {
	var err error
	Schema, err = BuildSchema()
	if err != nil {
		return err
	}

	e := echo.New()
	e.Pre(redirectHTTPS)
	e.Use(recordMetrics)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339} method=${method} path=${path} host=${host} status=${status} bytes_in=${bytes_in} bytes_out=${bytes_out}\n",
		Skipper: func(c echo.Context) bool {
			return !strings.HasPrefix(c.Request().Host(), "localhost")
		},
	}))
	e.Use(middleware.Recover())
	e.Static("/", "public")
	e.Get("/.well-known/acme-challenge/:challenge", letsEncrypt)
	e.File("/graphiql", "public/graphiql.html")
	e.File("/faq", "public/faq.html")

	h := handler.New(&handler.Config{
		Schema: &Schema,
		Pretty: true,
	})

	e.Any("/graphql", standard.WrapHandler(h))

	// Run the server
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))

	err = e.Run(standard.New(addr))
	return err
}

func letsEncrypt(c echo.Context) error {
	param := c.Param("challenge")
	if param == encryptChallenge {
		return c.String(http.StatusOK, encryptKey)
	}

	return errors.New("Let's Encrypt challenge did not match")
}

// Heroku specific HTTPS redirect
func redirectHTTPS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !isDyno {
			return next(c)
		}

		req := c.Request()
		host := req.Host()
		uri := req.URI()
		proto := req.Header().Get("X-Forwarded-Proto")
		if proto != "https" {
			metrics.Incr("http_to_https")
			return c.Redirect(http.StatusMovedPermanently, "https://"+host+uri)
		}
		return next(c)
	}
}

func recordMetrics(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		start := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}
		path := req.URL().Path()
		if path == "" {
			path = "/"
		}

		tags := map[string]string{
			"path":   path,
			"host":   req.Host(),
			"method": req.Method(),
			"status": strconv.Itoa(res.Status()),
		}

		fields := map[string]interface{}{
			"bytes_out": res.Size(),
			"took":      time.Since(start).Seconds(),
			"bytes_in":  req.ContentLength(),
		}

		err := metrics.Write("http_request", tags, fields)
		if err != nil {
			// We don't want to return server errors on failure to record metrics
			fmt.Println(err)
		}
		return nil
	}
}
