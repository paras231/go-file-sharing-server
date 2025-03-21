package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		fmt.Fprintf(w, "%v: %v\n", name, headers)
	}
}

// streaming a large file

func streamLargeFileHandler(w http.ResponseWriter, req *http.Request) {
	file, err := os.Open("imgs/test.JPG")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=largefile.mp4")

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error streaming file", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/download", streamLargeFileHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
