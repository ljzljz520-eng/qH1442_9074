package main

import (
	"log"
	"net/http"
	"os"

	"aftercare/internal/aftercare"
	"aftercare/internal/api"
)

func main() {
	addr := os.Getenv("AFTERCARE_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	repository := aftercare.NewMemoryRepository()
	service := aftercare.NewService(repository)
	handler := api.NewHandler(service)

	log.Printf("aftercare service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
