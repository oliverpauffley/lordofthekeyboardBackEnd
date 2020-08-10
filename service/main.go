package main

import (
	"fmt"
	"keyboard/service/app"
	"keyboard/service/quotes"
	"net/http"
	"os"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("must set 'API_KEY'")
		os.Exit(1)
	}

	server := app.Server{
		Router:      http.NewServeMux(),
		QuoteSource: quotes.NewAPIClient("https://the-one-api.herokuapp.com/v1", apiKey),
	}

	server.Routes()

	http.ListenAndServe(":8090", server.Router)
}
