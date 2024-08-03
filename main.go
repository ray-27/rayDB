package main

import (
	"log"
	"net/http"

	"github.com/ray-27/rayDB.git/assets"
	"github.com/ray-27/rayDB.git/server"
)

func handlers() {
	http.HandleFunc("/", server.Home_Handler)
}

func main() {

	assets.RayDB_logo(7)
	handlers()

	log.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
