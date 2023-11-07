package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

var (
	opensearchHost string
)

func opensearchHandler(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	opensearchHost = service
	opensearchConnectionStr := fmt.Sprintf("http://%s:9200", opensearchHost)

	fmt.Fprintf(w, opensearchConnector(opensearchConnectionStr))
}

func cleanOpensearchOutput(sr *opensearchapi.Response) string {
	var mp map[string]interface{}
	err := json.NewDecoder(sr.Body).Decode(&mp)
	if err != nil {
		fmt.Print(err)
	}

	var values []string
	for k, v := range mp {
		if strings.HasPrefix(k, "hits") {
			valStr := fmt.Sprint(v)
			values = append(values, valStr)
		}
	}

	var matches []string
	r := regexp.MustCompile(`LAGOON_\w*:(?:[^\[\]\)\(\s])*`)
	for _, str := range values {
		matches = r.FindAllString(str, -1)
	}

	keyVals := connectorKeyValues(matches)
	host := fmt.Sprintf(`"SERVICE_HOST=%s"`, opensearchHost)
	opensearchOutput := host + "\n" + keyVals
	return opensearchOutput
}

func createOpensearchIndexDocument(client *opensearch.Client) {
	settings := strings.NewReader(`{
		'settings': {
			'index': {
				'number_of_shards': 1,
				'number_of_replicas': 0
				}
			}
		}`)

	_ = opensearchapi.IndicesCreateRequest{
		Index: "opensearch-test",
		Body:  settings,
	}

	m := make(map[string]string)
	for _, e := range os.Environ() {
		if i := strings.Index(e, "="); i >= 0 {
			m[e[:i]] = e[i+1:]
		}
	}

	jsonD, _ := json.Marshal(m)
	reqDoc := strings.NewReader(string(jsonD))

	docId := "1"
	req := opensearchapi.IndexRequest{
		Index:      "opensearch-test",
		DocumentID: docId,
		Body:       reqDoc,
	}
	_, err := req.Do(context.Background(), client)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(1 * time.Second)
}

func opensearchConnector(connectionString string) string {
	client, _ := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{connectionString},
	})

	createOpensearchIndexDocument(client)

	content := strings.NewReader(`{
		"size": 20,
		"query": {
			"match_all": {}
		}
	}`)

	searchReq := &opensearchapi.SearchRequest{
		Index: []string{"opensearch-test"},
		Body:  content,
	}

	searchResponse, _ := searchReq.Do(context.Background(), client)
	openSearchResults := cleanOpensearchOutput(searchResponse)
	return openSearchResults
}
