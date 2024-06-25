package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/pat"
	"github.com/yoosuf/VPay/pkg/handlers"
)

func main() {
	// Database connection
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Ensure the database connection is available
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Router setup
	router := pat.New()
	router.Post("/process-payment", handlers.ProcessPaymentHandler(db))
	router.Get("/request-otp", handlers.RequestOTPHandler(db))

	// Start the server
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", router))
}
