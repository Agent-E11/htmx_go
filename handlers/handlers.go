package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	mySql "github.com/agent-e11/htmx_go/sql"
	"github.com/agent-e11/htmx_go/tools"
	"github.com/julienschmidt/httprouter"
)

func HomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

    tmpl := template.Must(template.ParseFiles("index.gotmpl"))
    tmpl.Execute(w, data)
}

// Add a product to the database
func AddProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

    if name == "" || priceStr == "" {
        log.Println("Error: cannot parse empty fields")
        w.Header().Add("HX-Retarget", "#form-error")
        w.Header().Add("HX-Reswap", "innerHTML")
        fmt.Fprint(w, "No form fields can be empty")
        return
    }

    // Convert to data types
    price, err := strconv.ParseFloat(priceStr, 64)
    if err != nil {
        log.Printf("Error parsing price: %v\n", err)
        w.Header().Add("HX-Retarget", "#form-error")
        w.Header().Add("HX-Reswap", "innerHTML")
        fmt.Fprintf(w, "`%s` is not a valid number", priceStr)
        return
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
    tmpl := template.Must(template.ParseFiles("index.gotmpl"))
    tmpl.ExecuteTemplate(w,
        "film-list-element",
        product,
    )
}

func LoadDummyDataHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    db, err := tools.ConnectDatabase()
    defer db.Close()

    if err != nil {
        return
    }
    products, err := tools.LoadDummyData(db, "dummy.json")
    tmpl := template.Must(template.ParseFiles("index.gotmpl"))
    for _, p := range products {
        tmpl.ExecuteTemplate(w,
            "film-list-element",
            p,
        )
    }
}

