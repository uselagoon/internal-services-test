package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type funcType func() map[string]string

var funcToCall []funcType

func main() {

	handler := http.HandlerFunc(handleReq)
	mariaHandler := http.HandlerFunc(mariaHandler)
	postgresHandler := http.HandlerFunc(postgresHandler)
	http.Handle("/", handler)
	http.Handle("/mariadb", mariaHandler)
	http.Handle("/postgres", postgresHandler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	funcToCall = append(funcToCall, mariaDBConnector, postgresDBConnector)
	for _, conFunc := range funcToCall {
		resp := conFunc()
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
	}
}
