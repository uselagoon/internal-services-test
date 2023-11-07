package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"os"
)

func persistentStorageHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	fmt.Fprintf(w, persistentStorageConnector(path))
}

func persistentStorageConnector(route string) string {
	if route != os.Getenv("STORAGE_LOCATION") {
		return "Storage location is not defined - ensure format matches '/storage?path'"
	}
	path := route + "/storage.txt"
	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	for _, e := range os.Environ() {
		_, err := f.WriteString(strconv.Quote(e) + "\n")
		if err != nil {
			log.Print(err)
		}
		e := f.Sync()
		if e != nil {
			log.Print(e)
		}
	}

	fileBuffer, err := os.ReadFile(path)
	var results = string(fileBuffer)
	return results
}
