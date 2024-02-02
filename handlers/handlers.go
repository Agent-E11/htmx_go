package handlers

import (
    "net/http"
    "strconv"
    "log"
    "text/template"

    mySql "github.com/agent-e11/htmx_go/sql"
    "github.com/agent-e11/htmx_go/tools"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
    db, err := tools.ConnectDatabase()
    defer db.Close()
    if err != nil {
        log.Println("Error connecting to database")
        return
    }

    products := []mySql.Product{}
    rows, err := db.Query("SELECT name, price, available FROM product")
    if err != nil {
        log.Println("Error fetching products")
        return
    }
    defer rows.Close()

    var name string
    var price float64
    var available bool

    for rows.Next() {
        err := rows.Scan(&name, &price, &available)
        if err != nil {
            log.Printf("Error converting row to product: %v\n", err)
        }
        
        products = append(products, mySql.Product{ Name: name, Price: price, Available: available })
    }

    data := map[string][]mySql.Product{
        "Products": products,
    }

    tmpl := template.Must(template.ParseFiles("index.html"))
    tmpl.Execute(w, data)
}

// Add a product to the database
func AddProduct(w http.ResponseWriter, r *http.Request) {
    // Connect to database
    db, err := tools.ConnectDatabase()
    defer db.Close()
    if err != nil {
        log.Println("Error connecting to database")
        return
    }

    // Get values from form
    name := r.PostFormValue("name")
    priceStr := r.PostFormValue("price")
    availableStr := r.PostFormValue("available")

    // Convert to data types
    price, err := strconv.ParseFloat(priceStr, 64)
    if err != nil {
        log.Fatalf("Error parsing price: %v\n", err)
    }
    available := availableStr != ""

    // Create product and insert into database
    product := mySql.Product{
        Name: name,
        Price: price,
        Available: available,
    }
    mySql.InsertProduct(db, product)

    // Return product as html fragment
    tmpl := template.Must(template.ParseFiles("index.html"))
    tmpl.ExecuteTemplate(w,
        "film-list-element",
        product,
    )
}
