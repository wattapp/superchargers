package main

import (
	"time"

	"github.com/wattapp/superchargers/pkg/database"
	"github.com/wattapp/superchargers/pkg/metrics"
	"github.com/wattapp/superchargers/pkg/web"
)

func main() {
	err := metrics.Connect()
	if err != nil {
		panic(err)
	}

	go metrics.Stats(time.Minute)

	_, err = database.Connect()
	if err != nil {
		panic(err)
	}

	err = web.Run()
	if err != nil {
		panic(err)
	}
}
