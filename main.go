package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {

	dbEnv := getEnv("DBENV")

	db, err := sql.Open("mysql", "api:api@tcp(mariadb:3305)/infrastructure")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// var version string
	// err = db.QueryRow("select version()").Scan(&version)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(version)

	// Route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, dbEnv)
	})

	port := ":8080"
	fmt.Println("Server is running on port" + port)

	// Start server on port specified above
	log.Fatal(http.ListenAndServe(port, nil))
}

func getEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}
