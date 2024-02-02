package handlers

import (
    "net/http"
    "strconv"
    "time"
    "log"
    "text/template"

    mySql "github.com/agent-e11/htmx_go/sql"
)

func HomePage(w http.ResponseWriter, r *http.Request) {

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

func AddProduct(w http.ResponseWriter, r *http.Request) {
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
