package dbcontrol

import (
    "database/sql"
    "errors"
    "fmt"
    "log"
    "os"

    _ "github.com/joho/godotenv/autoload"
    _ "github.com/lib/pq"
)

type Product struct {
    Name string
    Price float64
    Available bool
}

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

func CreateProductTable(db *sql.DB) {
    query := `CREATE TABLE IF NOT EXISTS product (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        price NUMERIC(6,2) NOT NULL,
        available BOOLEAN,
        created timestamp DEFAULT NOW()
    )`

    _, err := db.Exec(query)
    if err != nil {
        log.Fatalf("Error creating database table: %v\n", err)
    }
}

func InsertProduct(db *sql.DB, product Product) int {
    query := `INSERT INTO product (name, price, available)
        VALUES ($1, $2, $3) RETURNING id`

    var pk int

    err := db.QueryRow(query, product.Name, product.Price, product.Available).Scan(&pk)
    if err != nil {
        log.Fatal(err)
    }

    return pk
}
