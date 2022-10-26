package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vanng822/go-solr/solr"
)

var (
	service           = os.Getenv("SOLR_HOST")
	solrConnectionStr = fmt.Sprintf("http://%s:8983/solr", service)
)

func solrConnector() {
	si, err := solr.NewSolrInterface(solrConnectionStr, "drupal")
	if err != nil {
		log.Print(err)
	} else {
		log.Print(solrConnectionStr)
	}
	query := solr.NewQuery()
	query.Q("*:*")
	s := si.Search(query)
	r, err := s.Result(nil)
	if err != nil {
		log.Print("Error: ", err)
	} else {
		fmt.Println(r.Results.Docs)
	}
}
