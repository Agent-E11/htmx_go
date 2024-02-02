package tools

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	mySql "github.com/agent-e11/htmx_go/sql"
)

// Connect to the database
// Errors if environment variables are set incorrectly
// Uses the following environment variables:
//
// - POSTGRES_PASSWORD : Required
// - POSTGRES_PORT : Not required (defaults to 5432)
// - POSTGRES_DBNAME : Required
func ConnectDatabase() (*sql.DB, error) {
    // Get variables from environment, handle errors
    password := os.Getenv("POSTGRES_PASSWORD")
    port := os.Getenv("POSTGRES_PORT")
    dbName := os.Getenv("POSTGRES_DBNAME")
    if password == "" {
        log.Println("Error reading POSTGRES_PASSWORD")
        return nil, errors.New("Error reading POSTGRES_PASSWORD")
    }
    if port != "5432" && port != "" {
        log.Println("Different port numbers are unsupported at this time. Use port 5432")
        return nil, errors.New("Different port numbers are unsupported at this time. Use port 5432")
    } else {
        port = "5432"
    }
    if dbName == "" {
        log.Println("Error reading POSTGRES_DBNAME")
        return nil, errors.New("Error reading POSTGRES_DBNAME")
    }

    // Construct connection string
    connStr := fmt.Sprintf("postgres://postgres:%s@localhost:%s/%s?sslmode=disable", password, port, dbName)
    
    // Return the database connection and error
    return sql.Open("postgres", connStr)
}

// Load data from file (encoded as json) into database
func LoadDummyData(db *sql.DB, filename string) error {
    bytes, err := os.ReadFile(filename)
    if err != nil {
        return err
    }

    productJson := string(bytes)

    var products []mySql.Product

    err = json.Unmarshal([]byte(productJson), &products)
    if err != nil {
        return err
    }

    for _, p := range products {
        mySql.InsertProduct(db, p)
    }

    return nil
}
