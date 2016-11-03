package web

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

var Schema graphql.Schema

func Run() error {
	var err error
	Schema, err = BuildSchema()
	if err != nil {
		return err
	}

	e := echo.New()

	// Configure middleware
	// Heroku has its own logging
	if os.Getenv("DYNO") == "" {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "time=${time_rfc3339} method=${method} path=${path} host=${host} status=${status} bytes_in=${bytes_in} bytes_out=${bytes_out}\n",
		}))
	}
	e.Use(middleware.Recover())
	e.Static("/", "public")
	e.Get("/.well-known/acme-challenge/:challenge", letsEncrypt)
	e.File("/graphiql", "public/graphiql.html")

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
	challenge := os.Getenv("LETS_ENCRYPT_CHALLENGE")
	param := c.Param("challenge")
	if param == challenge {
		return c.String(http.StatusOK, os.Getenv("LETS_ENCRYPT_KEY"))
	}

	return errors.New("Let's Encrypt challenge did not match")
}
