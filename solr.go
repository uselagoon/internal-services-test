package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/vanng822/go-solr/solr"
)

var (
	service           = os.Getenv("SOLR_HOST")
	solrConnectionStr = fmt.Sprintf("http://%s:8983/solr", service)
)

func solrHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, convertSolrDoc(solrConnector()))
}

func convertSolrDoc(f []solr.Document) string {
	b := new(bytes.Buffer)
	for _, doc := range f {
		fmt.Fprintf(b, "\"%s\"\n", doc)
	}
	v := b.String()
	v1 := strings.ReplaceAll(v, "[", "")
	v2 := strings.ReplaceAll(v1, "]", "")
	v3 := strings.ReplaceAll(v2, "map", "")
	return v3
}

func solrConnector() []solr.Document {
	si, err := solr.NewSolrInterface(solrConnectionStr, "core")
	if err != nil {
		log.Print(err)
	} else {
		log.Print(solrConnectionStr)
	}
	si.DeleteAll()
	d := solr.Document{}
	for _, e := range os.Environ() {

		pair := strings.SplitN(e, "=", 2)
		d.Set(pair[0], pair[1])
		if err != nil {
			panic(err.Error())
		}
	}
	arrayD := []solr.Document{}
	arrayD = append(arrayD, d)
	si.Add(arrayD, 1, nil)
	si.Commit()
	query := solr.NewQuery()
	query.Q("*:*")
	query.FieldList("LAGOON_*")
	s := si.Search(query)
	r, err := s.Result(nil)
	if err != nil {
		log.Print("Error: ", err)
	}
	return r.Results.Docs
}
