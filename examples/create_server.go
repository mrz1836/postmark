package main

import (
	"context"
	"github.com/mrz1836/postmark"
	"log"
)

func main() {
	client := postmark.NewClient("[SERVER-TOKEN]", "[ACCOUNT-TOKEN]")

	server := postmark.Server{
		Name:  "Edit Server Name",
		Color: "Red",
	}
	server, err := client.CreateServer(context.TODO(), server)
	if err != nil {
		panic(err)
	}

	log.Println(server.ID)
	log.Println(server.Name)
	log.Println(server.Color)
}
