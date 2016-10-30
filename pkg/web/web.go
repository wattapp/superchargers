package web

import (
	"fmt"
	"log"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

var Schema graphql.Schema

func Run() {
	var err error
	Schema, err = BuildSchema()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	// Configure middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := handler.New(&handler.Config{
		Schema: &Schema,
		Pretty: true,
	})

	e.Any("/graphql", standard.WrapHandler(h))

	// Run the server
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))

	err = e.Run(standard.New(addr))
	if err != nil {
		log.Fatal(err)
	}
}
