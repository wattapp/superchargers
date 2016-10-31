package main

import (
	"fmt"

	"github.com/wattapp/superchargers/pkg/database"
	"github.com/wattapp/superchargers/pkg/location"
)

func main() {
	fmt.Println("Starting to update all locations...")
	_, err := database.Connect()
	if err != nil {
		panic(err)
	}

	var added, updated int
	added, updated, err = location.Sync()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Added: %d, Updated: %d\n", added, updated)
}
