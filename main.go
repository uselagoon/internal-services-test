package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

type LandingPageData struct {
	Services []Service
}

type Service struct {
	Type string
	Name string
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/mariadb", mariadbHandler)
	r.HandleFunc("/postgres", postgresHandler)
	r.HandleFunc("/redis", redisHandler)
	r.HandleFunc("/solr", solrHandler)
	r.HandleFunc("/mongo", mongoHandler)
	r.HandleFunc("/opensearch", opensearchHandler)
	r.HandleFunc("/storage", persistentStorageHandler)
	r.HandleFunc("/mysql", mariadbHandler)
	r.HandleFunc("/", handleReq)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":3000", handler(r)))
}

func handler(m *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		driver := strings.ReplaceAll(r.URL.Path, "/", "")
		service := r.URL.Query().Get("service")
		incompatibleError := fmt.Sprintf("%s is not a compatible driver with service: %s", driver, service)
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, incompatibleError, http.StatusInternalServerError)
			}
		}()
		timeoutHandler := http.TimeoutHandler(m, 3*time.Second, incompatibleError)
		timeoutHandler.ServeHTTP(w, r)
	}
}

func getServices() []Service {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()

	containers, err := c.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		fmt.Println(err)
	}

	var serviceList []Service
	serviceRegexp := regexp.MustCompile(`/internal-services-test_(.*?)_1`)
	typeRegexp := regexp.MustCompile(`(\w*)-.*`)
	for _, ctr := range containers {
		for _, name := range ctr.Names {
			if strings.Contains(name, "internal-services-test") && !strings.Contains(name, "commons") {
				serviceName := serviceRegexp.FindStringSubmatch(name)
				serviceType := typeRegexp.FindStringSubmatch(serviceName[1])
				serviceList = append(serviceList, Service{Type: serviceType[1], Name: serviceName[1]})
			}
		}
	}

	return serviceList
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("landingTemplate.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serviceList := getServices()
	data := LandingPageData{
		Services: serviceList,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	localService := strings.ReplaceAll(cleanRoute, ".", "-")
	replaceHyphen := strings.ReplaceAll(localService, "-", "_")
	lagoonService := strings.ToUpper(replaceHyphen)
	return localService, lagoonService
}
