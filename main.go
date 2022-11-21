package main

import (
	"log"
	"net/http"

	"github.com/SalviCF/authorization-server/router"
)

func main() {

	// Init custom Router
	rtr := router.InitRouter()

	// Serve contents
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(rtr.Serve)))
}
