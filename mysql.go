package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	mysqlVersion       string
	mysqlConnectionStr string
)

func mysqlHandler(w http.ResponseWriter, r *http.Request) {
	mysqlPath := r.URL.Path
	localService, lagoonService := cleanRoute(mysqlPath)
	mysqlUser := getEnv(fmt.Sprintf("%s_USERNAME", lagoonService), "lagoon")
	mysqlPassword := getEnv(fmt.Sprintf("%s_PASSWORD", lagoonService), "lagoon")
	mysqlHost := getEnv(fmt.Sprintf("%s_HOST", lagoonService), localService)
	mysqlPort := getEnv(fmt.Sprintf("%s_PORT", lagoonService), "3306")
	mysqlDatabase := getEnv(fmt.Sprintf("%s_DATABASE", lagoonService), "lagoon")

	mysqlConnectionStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
	log.Print(fmt.Sprintf("Using %s as the connstring", mysqlConnectionStr))

	fmt.Fprintf(w, dbConnectorPairs(mysqlConnector(mysqlConnectionStr), mysqlVersion))
}

func mysqlConnector(connectionString string) map[string]string {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	createTable := "CREATE TABLE IF NOT EXISTS env(env_key text, env_value text)"
	_, err = db.Exec(createTable)
	if err != nil {
		log.Print(err)
	}

	query := "INSERT INTO env(env_key, env_value) VALUES (?, ?)"

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		_, err := db.Exec(query, pair[0], pair[1])
		if err != nil {
			log.Print(err)
		}
	}

	q := "LAGOON_%"
	rows, err := db.Query(`SELECT * FROM env where env_key LIKE ?`, q)
	if err != nil {
		log.Print(err)
	}

	db.QueryRow("SELECT VERSION()").Scan(&mysqlVersion)

	defer rows.Close()
	results := make(map[string]string)
	for rows.Next() {
		var envKey, envValue string
		_ = rows.Scan(&envKey, &envValue)
		results[envKey] = envValue
	}

	return results
}
