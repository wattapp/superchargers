package web

import (
	"fmt"
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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/assets", "public/assets")
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
