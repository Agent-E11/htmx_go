package tools

import (
    "database/sql"
    "errors"
    "fmt"
    "log"
    "os"

    "github.com/agent-e11/htmx_go/dbcontrol"
    "github.com/brianvoe/gofakeit/v6"
)

// Connect to the database
// Errors if environment variables are set incorrectly
// Uses the following environment variables:
//
// - POSTGRES_PASSWORD : Required
// - POSTGRES_HOSTNAME : Not required (defaults to localhost)
// - POSTGRES_PORT : Not required (defaults to 5432)
// - POSTGRES_DBNAME : Required
func ConnectDatabase() (*sql.DB, error) {
    // Get variables from environment, handle errors
    password := os.Getenv("POSTGRES_PASSWORD")
    hostname := os.Getenv("POSTGRES_HOSTNAME")
    port := os.Getenv("POSTGRES_PORT")
    dbName := os.Getenv("POSTGRES_DBNAME")
    if password == "" {
        log.Println("Error reading POSTGRES_PASSWORD")
        return nil, errors.New("Error reading POSTGRES_PASSWORD")
    }
    if hostname == "" {
        hostname = "localhost"
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
    connStr := fmt.Sprintf("postgres://postgres:%s@%s:%s/%s?sslmode=disable", password, hostname, port, dbName)
    
    // Return the database connection and error
    return sql.Open("postgres", connStr)
}

// Load data from file (encoded as json) into database
func LoadDummyData(db *sql.DB) []dbcontrol.Product {
    var products []dbcontrol.Product
    for i := 0; i < 5; i++ {
        products = append(products, dbcontrol.Product{
            Name: gofakeit.ProductName(),
            Price: gofakeit.Price(0.0, 200.0),
            Available: gofakeit.Bool(),
        })
    }

    for _, p := range products {
        dbcontrol.InsertProduct(db, p)
    }

    return products
}

func Filter[T any](slice []T, test func(T) bool) []T {
    var ret []T
    for _, s := range slice {
        if test(s) {
            ret = append(ret, s)
        }
    }

    return ret
}
