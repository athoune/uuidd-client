package main

import (
	"fmt"
	"os"

	"github.com/athoune/uuidd-client/uuidd"
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
}
