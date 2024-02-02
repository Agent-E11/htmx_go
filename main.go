package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/agent-e11/htmx_go/handlers"
	mySql "github.com/agent-e11/htmx_go/sql"
	"github.com/agent-e11/htmx_go/tools"
)

type Film struct {
    Title string
    Director string
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

    mySql.CreateProductTable(db)
    err = tools.LoadDummyData(db, "dummy.json")
    if err != nil {
        log.Fatalf("Error loading dummy data: %v\n", err)
    }

    // HTTP Server
    fmt.Println("Running htmx server...")

    http.HandleFunc("/", handlers.HomePage)
    http.HandleFunc("/add-product/", handlers.AddProduct)

    log.Fatal(http.ListenAndServe(":8000", nil))
}
