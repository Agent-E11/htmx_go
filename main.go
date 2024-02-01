package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
    "strconv"

    mySql "github.com/agent-e11/htmx_go/sql"
)

type Film struct {
    Title string
    Director string
}

func main() {
    fmt.Println("Running htmx server...")

    page := func(w http.ResponseWriter, r *http.Request) {

        products := map[string][]mySql.Product{
            "Products": {
                { Name: "Book", Price: 15.55, Available: true },
                { Name: "TV", Price: 199.99, Available: true },
                { Name: "Mouse", Price: 9.99, Available: true },
                { Name: "Keyboard", Price: 12.99, Available: true },
            },
        }

        tmpl := template.Must(template.ParseFiles("index.html"))
        tmpl.Execute(w, products)
    }

    addFilm := func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(time.Second*1)
        name := r.PostFormValue("name")
        priceStr := r.PostFormValue("price")
        availableStr := r.PostFormValue("available")

        price, err := strconv.ParseFloat(priceStr, 64)
        if err != nil {
            log.Fatalf("Error parsing price: %v\n", err)
        }

        available := availableStr != ""

        tmpl := template.Must(template.ParseFiles("index.html"))
        tmpl.ExecuteTemplate(w,
            "film-list-element",
            mySql.Product {
                Name: name,
                Price: price,
                Available: available,
            },
        )
    }

    http.HandleFunc("/", page)
    http.HandleFunc("/add-film/", addFilm)

    log.Fatal(http.ListenAndServe(":8000", nil))
}
