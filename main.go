package main

import (
	"fmt"
	"log"
	"net/http"

    "github.com/agent-e11/htmx_go/handlers"
)

type Film struct {
    Title string
    Director string
}

func main() {
    fmt.Println("Running htmx server...")

    http.HandleFunc("/", handlers.HomePage)
    http.HandleFunc("/add-product/", handlers.AddProduct)

    log.Fatal(http.ListenAndServe(":8000", nil))
}
