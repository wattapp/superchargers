package main

import (
	"github.com/wattapp/superchargers/pkg/database"
	"github.com/wattapp/superchargers/pkg/location"
	"github.com/wattapp/superchargers/pkg/web"
)

func main() {
	_, err := database.Connect()
	if err != nil {
		panic(err)
	}

	go location.Update()

	err = web.Run()
	if err != nil {
		panic(err)
	}
}
