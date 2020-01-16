package main

import (
	"io"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Create("./newFile")
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(file, r.Body)
	if err != nil {
		panic(err)
	}
	w.Write([]byte("upload success"))
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":5050", nil)
}