package main

import (
	"log"
	"net/http"
	"os"

	"github.com/adelgado0723/processor"
)

func main() {
	client := processor.NewAuthenticationClient(
		http.DefaultClient, "https", "us-street.api.smartystreets.com",
		"19e816ad-66c0-d816-5e84-f625ecf14765", "UJV7R6qfdQZ5pRO7epX9")

	pipeline := processor.NewPipeline(os.Stdin, os.Stdout, client, 8)

	if err := pipeline.Process(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
