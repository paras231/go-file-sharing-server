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

// stream a video

func streamVideo(w http.ResponseWriter, req *http.Request) {

	file, err := os.Open("imgs/pricevideo.mp4")
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	// Set the content type to video/mp4
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, file)
	if err != nil {
		fmt.Println("Error streaming video:", err)
		return
	}
	flusher.Flush()
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/download", streamLargeFileHandler)
	http.HandleFunc("/stream-video", streamVideo)
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
