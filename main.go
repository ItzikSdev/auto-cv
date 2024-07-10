package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/dslipak/pdf"
	"github.com/nguyenthenguyen/docx"
)

type FileInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
	Text string `json:"text"`
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile(r.RequestURI[1:])
	if file == nil {
		http.Error(w, "error", http.StatusNotFound)
	}
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileType := filepath.Ext(handler.Filename)
	if fileType == ".pdf" {
		content, err := pdf.NewReader(file, handler.Size)
		if err != nil {
			panic(err)
		}
		plainText, err := content.GetPlainText()
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
		}
		text, err := io.ReadAll(plainText)
		if err != nil {
			http.Error(w, "error", http.StatusInternalServerError)
		}
		f := FileInfo{
			Name: handler.Filename,
			Type: fileType[1:],
			Size: handler.Size,
			Text: string(text),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(f); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
		}
		return
	}
	if fileType == ".docx" {
		content, err := docx.ReadDocxFromMemory(file, handler.Size)
		if err != nil {
			panic(err)
		}
		f := FileInfo{
			Name: handler.Filename,
			Type: fileType[1:],
			Size: handler.Size,
			Text: content.Editable().GetContent(),
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(f); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "unsupported file type", http.StatusBadRequest)
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}

func main() {
	setupRoutes()
}
