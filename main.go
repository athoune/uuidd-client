package main

import (
	"fmt"
	"os"

	"github.com/athoune/uuidd-client/uuidd"
	"github.com/google/uuid"
)

func main() {
	path := os.Getenv("UUID_PATH")
	if path == "" {
		path = "/run/uuidd/request"
	}
	client := uuidd.New(path)

	u, err := client.TimeUUID()
	if err != nil {
		panic(err)
	}
	fmt.Println(u)

	err = client.BulkTimeUUID(3, func(uu uuid.UUID) error {
		fmt.Println(uu)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
