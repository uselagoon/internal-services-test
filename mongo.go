package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	machineryEnvVars "github.com/uselagoon/machinery/utils/variables"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	mongoConnectionStr string
	mongoHost          string
	database           string
)

func mongoHandler(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	localService, lagoonService := cleanRoute(service)
	mongoUser := machineryEnvVars.GetEnv(fmt.Sprintf("%s_USERNAME", lagoonService), "lagoon")
	mongoPassword := machineryEnvVars.GetEnv(fmt.Sprintf("%s_PASSWORD", lagoonService), "lagoon")
	mongoHost = machineryEnvVars.GetEnv(fmt.Sprintf("%s_HOST", lagoonService), localService)
	mongoPort := machineryEnvVars.GetEnv(fmt.Sprintf("%s_PORT", lagoonService), "27017")
	mongoDatabase := machineryEnvVars.GetEnv(fmt.Sprintf("%s_DATABASE", lagoonService), "lagoon")

	if mongoHost != localService {
		mongoConnectionStr = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", mongoUser, mongoPassword, mongoHost, mongoPort, mongoDatabase)
		database = mongoDatabase
	} else {
		mongoConnectionStr = fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)
		database = mongoDatabase
	}
	log.Print(fmt.Sprintf("Using %s as the connstring", mongoConnectionStr))

	fmt.Fprintf(w, "%s", mongoConnector(mongoConnectionStr, database))
}

func cleanMongoOutput(docs []bson.M) string {
	valStr := fmt.Sprint(docs)
	r := regexp.MustCompile(`"Key":"(LAGOON_\w+)","value":"([^"]+)"`)
	matches := r.FindAllStringSubmatch(valStr, -1)
	var mongoResults []string
	for _, str := range matches {
		mongoResults = append(mongoResults, fmt.Sprintf("%s %s", str[1], str[2]))
	}
	b := new(bytes.Buffer)
	for _, value := range mongoResults {
		v := strings.SplitN(value, " ", 2)
		fmt.Fprintf(b, "\"%s=%s\"\n", v[0], v[1])
	}
	host := fmt.Sprintf(`"SERVICE_HOST=%s"`, mongoHost)
	mongoOutput := host + "\n" + b.String()
	return mongoOutput
}

func mongoConnector(connectionString string, database string) string {
	client, err := mongo.Connect(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Print(err)
	}

	envCollection := client.Database(database).Collection("env-vars")

	deleteFilter := bson.D{{}}
	_, err = envCollection.DeleteMany(context.TODO(), deleteFilter)
	if err != nil {
		log.Print(err)
	}

	environmentVariables := []interface{}{}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		bsonData := bson.D{{Key: "Key", Value: pair[0]}, {Key: "value", Value: pair[1]}}
		environmentVariables = append(environmentVariables, bsonData)
		if err != nil {
			log.Print(err)
		}
	}

	_, err = envCollection.InsertMany(context.TODO(), environmentVariables)
	if err != nil {
		log.Print(err)
	}
	filter := bson.D{{Key: "Key", Value: bson.Regex{Pattern: "LAGOON", Options: ""}}}
	cursor, _ := envCollection.Find(context.TODO(), filter, options.Find().SetProjection(bson.M{"_id": 0}))
	var docs []bson.M
	if err = cursor.All(context.TODO(), &docs); err != nil {
		log.Print(err)
	}
	mongoOutput := cleanMongoOutput(docs)
	return mongoOutput
}
