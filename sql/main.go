package main

import (
	"database/sql"
	"log"

    _ "github.com/lib/pq"
)

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

    log.Print("Successfully connected to database")
}
