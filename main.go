package main

import (
	handler "github.com/gouef/github-lib-template/api"
	"log"
	"net/http"
)

func main() {
	// Tento handler bude obsluhovat všechny požadavky na /api/get-action
	http.HandleFunc("/api/get-action", handler.GetAction)
	http.HandleFunc("/api/get-go-dependency", handler.GetGoDependency)
	http.HandleFunc("/", handler.GetAction)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
