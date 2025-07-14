package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	store, err := NewStore(dsn)
	if err != nil {
		log.Fatal(err)
	}
	server := &Server{store: store, hub: NewHub()}

	log.Println("Go backend listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", server.routes()))
}
