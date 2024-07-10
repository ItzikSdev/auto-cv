package main

import (
	"net/http"
)

func setupRoutes() {
	http.HandleFunc("/upload", UploadFile)
	http.ListenAndServe(":8080", nil)
}

func main() {
	setupRoutes()
}
