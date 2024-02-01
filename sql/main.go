package sql

import (
	"database/sql"
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

func main() {
    password := os.Getenv("POSTGRES_PASSWORD")
    port := os.Getenv("POSTGRES_PORT")
    dbName := os.Getenv("POSTGRES_DBNAME")
    if password == "" { log.Fatal("Error reading POSTGRES_PASSWORD") }
    if port != "5432" && port != "" {
        log.Fatal("Different port numbers are unsupported at this time. Use port 5432")
    } else {
        port = "5432"
    }
    if dbName == "" { log.Fatal("Error reading POSTGRES_DBNAME") }

    connStr := fmt.Sprintf("postgres://postgres:%s@localhost:%s/%s?sslmode=disable", password, port, dbName)
    
    db, err := sql.Open("postgres", connStr)
    defer db.Close()

    if err != nil {
        log.Fatalf("Error opening db: %v\n", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatalf("Error pinging db: %v\n", err)
    }

    log.Print("Connected to database...")

    CreateProductTable(db)

    data := []Product{}
    rows, err := db.Query("SELECT name, price, available FROM product")
    if err != nil {
        log.Fatalf("Error querying db: %v\n", err)
    }
    defer rows.Close()

    var name string
    var price float64
    var available bool

    for rows.Next() {
        err := rows.Scan(&name, &price, &available)
        if err != nil {
            log.Fatal(err)
        }
        
        data = append(data, Product{ name, price, available })
    }

    fmt.Println(data)

    p := Product{ "Book", 15.55, true }

    pk := InsertProduct(db, p)

    fmt.Println("ID =", pk)
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
        log.Fatal(err)
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
