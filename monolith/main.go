package main

import (
	"log"
	"net/http"

	"github.com/mhausenblas/imgn/monolith/handlers"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./ui")))
	http.HandleFunc("/upload", handlers.UploadFile)
	http.HandleFunc("/gallery", handlers.ListFiles)
	log.Println("Running")
	http.ListenAndServe(":8080", nil)
}
