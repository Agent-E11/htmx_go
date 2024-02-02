package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/agent-e11/htmx_go/handlers"
	mySql "github.com/agent-e11/htmx_go/sql"
	"github.com/agent-e11/htmx_go/tools"
)

type Film struct {
    Title string
    Director string
}

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
    http.HandleFunc("/load-dummy-data/", loadDummyDataHandler)
    //http.HandleFunc("/delete-all-data/", deleteAllDataHandler)

    log.Fatal(http.ListenAndServe(":8000", nil))
}

func loadDummyDataHandler(w http.ResponseWriter, r *http.Request) {
    db, err := tools.ConnectDatabase()
    defer db.Close()

    if err != nil {
        return
    }
    products, err := tools.LoadDummyData(db, "dummy.json")
    tmpl := template.Must(template.ParseFiles("index.html"))
    for _, p := range products {
        tmpl.ExecuteTemplate(w,
            "film-list-element",
            p,
        )
    }
}

//func deleteAllDataHandler(w http.ResponseWriter, r *http.Request) {
//    db, err := tools.ConnectDatabase()
//    defer db.Close()
//
//    if err != nil {
//        return
//    }
//}
