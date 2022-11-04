package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

type funcType func() map[string]string

func main() {

	handler := http.HandlerFunc(handleReq)
	mariaHandler := http.HandlerFunc(mariaHandler)
	postgresHandler := http.HandlerFunc(postgresHandler)
	solrHandler := http.HandlerFunc(solrHandler)
	http.Handle("/", handler)
	http.Handle("/mariadb", mariaHandler)
	http.Handle("/postgres", postgresHandler)
	http.Handle("/solr", solrHandler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	var funcToCall []funcType
	for _, conFunc := range funcToCall {
		fmt.Fprintf(w, createKeyValuePairs(conFunc()))
	}
}

func createKeyValuePairs(f map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range f {
		fmt.Fprintf(b, "\"%s=%s\"\n", key, value)
	}
	return b.String()
}
