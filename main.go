package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	machinery "github.com/uselagoon/machinery/utils/variables"
	"gopkg.in/yaml.v3"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
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

// getServices Parses docker-compose.yml and returns a list of services
func getServices() []Service {
	type BuildConfig struct {
		Context    string `yaml:"context,omitempty"`
		Dockerfile string `yaml:"dockerfile,omitempty"`
	}

	type ServiceConfig struct {
		Image   string            `yaml:"image,omitempty"`
		Labels  map[string]string `yaml:"labels,omitempty"`
		Build   BuildConfig       `yaml:"build,omitempty"`
		Ports   []string          `yaml:"ports,omitempty"`
		Volumes []string          `yaml:"volumes,omitempty"`
	}

	type DockerCompose struct {
		Services map[string]ServiceConfig `yaml:"services"`
		Volumes  map[string]interface{}   `yaml:"volumes,omitempty"`
	}

	compose := &DockerCompose{}
	data, err := os.ReadFile("docker-compose.yml")
	if err != nil {
		log.Fatalf("Error reading docker-compose.yml: %v", err)
	}

	err = yaml.Unmarshal(data, compose)
	if err != nil {
		log.Fatalf("Error parsing docker-compose.yml: %v", err)
	}

	var serviceList []Service
	typeRegexp := regexp.MustCompile(`^\w*`)
	for service := range compose.Services {
		if service != "web" {
			serviceType := typeRegexp.FindString(service)
			serviceList = append(serviceList, Service{Type: serviceType, Name: service})
		}
	}
	storage := machinery.GetEnv("STORAGE_LOCATION", "")
	if storage != "" {
		serviceList = append(serviceList, Service{Type: "storage", Name: storage})
	}

	return serviceList
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
