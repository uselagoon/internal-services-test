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

type funcType func() string

var funcToCall []funcType

func main() {

	handler := http.HandlerFunc(handleReq)
	http.Handle("/", handler)
	http.Handle("/mariadb", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	funcToCall = append(funcToCall, mariaDBConnector, postgresDBConnector)
	for _, conFunc := range funcToCall {
		resp := make(map[string]string)
		resp["Connection"] = conFunc()
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
	}
}

func mariaDBConnector() string {
	db, err := sql.Open(mariaDriver, mariaConnectionStr)
	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	createTable := "CREATE TABLE IF NOT EXISTS env(environment text)"
	createResult, err := db.Exec(createTable)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println(createResult)
	}

	query := "INSERT INTO env(environment) VALUES (?)"
	insertResult, err := db.Exec(query, mariaConnectionStr)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println(insertResult)
	}

	var envVars string
	err = db.QueryRow("select * from env").Scan(&envVars)
	if err != nil {
		log.Print(err)
	}

	return envVars
}

func postgresDBConnector() string {
	db, err := sql.Open(postgresHost, postgresConnectionStr)
	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	createTable := "CREATE TABLE IF NOT EXISTS env(environment text)"
	createResult, err := db.Exec(createTable)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println(createResult)
	}

	query := "INSERT INTO env(environment) VALUES ($1)"
	insertResult, err := db.Exec(query, postgresConnectionStr)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println(insertResult)
	}

	var envVars string
	err = db.QueryRow("select * from env").Scan(&envVars)
	if err != nil {
		log.Print(err)
	}

	return envVars
}
