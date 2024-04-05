package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type funcType func() map[string]string

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/mariadb", mariadbHandler)
	r.HandleFunc("/postgres", postgresHandler)
	r.HandleFunc("/redis", redisHandler)
	r.HandleFunc("/solr", solrHandler)
	r.HandleFunc("/mongo", mongoHandler)
	r.HandleFunc("/opensearch", opensearchHandler)
	r.HandleFunc("/storage", persistentStorageHandler)
	r.HandleFunc("/", handleReq)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":3000", handler(r)))
}

func handler(m *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			driver := strings.ReplaceAll(r.URL.Path, "/", "")
			service := r.URL.Query().Get("service")
			incompatibleError := fmt.Sprintf("%s is not a compatible driver with service: %s", driver, service)
			if err := recover(); err != nil {
				http.Error(w, incompatibleError, http.StatusInternalServerError)
			}
		}()
		handler := http.Handler(m)
		handler.ServeHTTP(w, r)
	}
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	var funcToCall []funcType
	for _, conFunc := range funcToCall {
		fmt.Fprintf(w, dbConnectorPairs(conFunc(), ""))
	}
}

func dbConnectorPairs(m map[string]string, connectorHost string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "\"%s=%s\"\n", key, value)
	}
	host := fmt.Sprintf(`"SERVICE_HOST=%s"`, connectorHost)
	connectorOutput := host + "\n" + b.String()
	return connectorOutput
}

func connectorKeyValues(values []string) string {
	b := new(bytes.Buffer)
	for _, value := range values {
		if value != "" {
			v := strings.SplitN(value, ":", 2)
			fmt.Fprintf(b, "\"%s=%s\"\n", v[0], v[1])
		}
	}
	return b.String()
}

func cleanRoute(basePath string) (string, string) {
	cleanRoute := strings.ReplaceAll(basePath, "/", "")
	localService := strings.ReplaceAll(cleanRoute, "10.", "10-")
	replaceHyphen := strings.ReplaceAll(localService, "-", "_")
	lagoonService := strings.ToUpper(replaceHyphen)
	return localService, lagoonService
}
