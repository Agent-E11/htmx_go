package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Product struct {
    Name string
    Price float64
    Available bool
}

func main() {
    connStr := "postgres://postgres:secret@localhost:5432/gopgtest?sslmode=disable"
    
    db, err := sql.Open("postgres", connStr)
    defer db.Close()

    if err != nil {
        log.Fatal(err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal(err)
    }

    log.Print("Connected to database...")

    createProductTable(db)

    product := Product{
        Name: "Book",
        Price: 15.55,
        Available: true,
    }
    pk := insertProduct(db, product)

    var name string
    var price string
    var available bool

    query := `SELECT name, price, available FROM product WHERE id = $1`
    err = db.QueryRow(query, pk).Scan(&name, &price, &available)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Fatalf("No rows were found with id %d", pk)
        }
        log.Fatal(err)
    }

    fmt.Printf("Name: %s\n", name)
    fmt.Printf("Name: %v\n", price)
    fmt.Printf("Name: %t\n", available)
}

func createProductTable(db *sql.DB) {
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

func insertProduct(db *sql.DB, product Product) int {
    query := `INSERT INTO product (name, price, available)
        VALUES ($1, $2, $3) RETURNING id`

    var pk int

    err := db.QueryRow(query, product.Name, product.Price, product.Available).Scan(&pk)
    if err != nil {
        log.Fatal(err)
    }

    return pk
}
