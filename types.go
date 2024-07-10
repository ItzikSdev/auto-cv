package main

type FileInfo struct {
	Name string  `json:"name"`
	Type DocType `json:"type"`
	Size int64   `json:"size"`
	Text string  `json:"text"`
}

type DocType struct {
	Docx string `json:"docx"`
	PDF  string `json:"pdf"`
}
