package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var (
	mariaUser          = os.Getenv("MARIADB_USER")
	mariaPassword      = os.Getenv("MARIADB_PASSWORD")
	mariaDB            = os.Getenv("MARIADB_DATABASE")
	mariaHost          = "mariadb"
	mariaPort          = 3306
	mariaDriver        = "mysql"
	mariaConnectionStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mariaUser, mariaPassword, mariaHost, mariaPort, mariaDB)

	postgresUser          = os.Getenv("POSTGRES_USER")
	postgresPassword      = os.Getenv("POSTGRES_PASSWORD")
	postgresDB            = os.Getenv("POSTGRES_DB")
	postgresHost          = "postgres"
	postgresSSL           = "disable"
	postgresConnectionStr = fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=%s host=%s",
		postgresUser, postgresPassword, postgresDB, postgresSSL, postgresHost)
)

func main() {

	handler := http.HandlerFunc(handleReq)
	http.Handle("/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp[mariaHost] = dbConnector(mariaDriver, mariaConnectionStr)
	resp[postgresHost] = dbConnector(postgresHost, postgresConnectionStr)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func dbConnector(driver string, connStr string) string {
	db, err := sql.Open(driver, connStr)
	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	var version string
	err = db.QueryRow("select version()").Scan(&version)
	if err != nil {
		log.Print(err)
	}
	return version
}
