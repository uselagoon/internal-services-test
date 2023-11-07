package main

import (
	"fmt"
	machineryEnvVars "github.com/uselagoon/machinery/utils/variables"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func persistentStorageHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	log.Print(fmt.Sprintf("Writing to %s", path))
	fmt.Fprintf(w, persistentStorageConnector(path))
}

func persistentStorageConnector(route string) string {
	if route != machineryEnvVars.GetEnv("STORAGE_LOCATION", "") {
		return "Storage location is not defined - ensure format matches '/storage?path=[path]'"
	}
	path := route + "/storage.txt"
	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	for _, e := range os.Environ() {
		if strings.Contains(e, "LAGOON_") {
			_, err := f.WriteString(strconv.Quote(e) + "\n")
			if err != nil {
				log.Print(err)
			}
			e := f.Sync()
			if e != nil {
				log.Print(e)
			}
		}
	}

	fileBuffer, err := os.ReadFile(path)
	var results = string(fileBuffer)
	storagePath := fmt.Sprintf(`"STORAGE_PATH=%s"`, path)
	storageResults := storagePath + "\n" + results
	return storageResults
}
