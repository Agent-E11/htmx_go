package main

import (
	"log"
	"net/http"

	"github.com/agent-e11/htmx_go/handlers"
	mySql "github.com/agent-e11/htmx_go/sql"
	"github.com/agent-e11/htmx_go/tools"
)

func main() {
    db, err := tools.ConnectDatabase()
    defer db.Close()

    if err != nil {
        log.Fatalf("Error opening db: %v\n", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatalf("Error pinging db: %v\n", err)
    }

    log.Print("Connected to database...")

    mySql.CreateProductTable(db)

    // Load dummy data
    //err = tools.LoadDummyData(db, "dummy.json")
    //if err != nil {
    //    log.Fatalf("Error loading dummy data: %v\n", err)
    //}

    // HTTP Server
    log.Print("Running htmx server...")

    http.HandleFunc("/", handlers.HomePage)
    http.HandleFunc("/add-product/", handlers.AddProduct)
    http.HandleFunc("/load-dummy-data/", handlers.LoadDummyDataHandler)
    //http.HandleFunc("/delete-all-data/", deleteAllDataHandler)

    log.Fatal(http.ListenAndServe(":8000", nil))
}

//func deleteAllDataHandler(w http.ResponseWriter, r *http.Request) {
//    db, err := tools.ConnectDatabase()
//    defer db.Close()
//
//    if err != nil {
//        return
//    }
//}
