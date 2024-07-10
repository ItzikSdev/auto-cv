package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dslipak/pdf"
	"github.com/nguyenthenguyen/docx"
)

var words = []string{
	"docker",
	"node.js",
	"cloudwatch",
	"blabla",
	"lambda",
	"ci/cd",
	"pipeline",
}
var exiestWord = []bool{}

func extractTextFromFile(f FileInfo) {
	for i := 0; i < len(words); i++ {
		lowerCase := strings.ToLower(f.Text)
		if strings.Contains(lowerCase, words[i]) {
			exiestWord = append(exiestWord, true)
			log.Printf("Find word %v in CV", words[i])
		} else {
			log.Printf("Word %v in missing", words[i])
		}
	}
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
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
			Type: DocType{PDF: fileType[1:]},
			Size: handler.Size,
			Text: string(text),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(f); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
		}
		extractTextFromFile(f)
		return
	}
	if fileType == ".docx" {
		content, err := docx.ReadDocxFromMemory(file, handler.Size)
		if err != nil {
			panic(err)
		}
		f := FileInfo{
			Name: handler.Filename,
			Type: DocType{Docx: fileType[1:]},
			Size: handler.Size,
			Text: content.Editable().GetContent(),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(f); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
		}
		extractTextFromFile(f)
		return
	}

	http.Error(w, "unsupported file type", http.StatusBadRequest)
}
