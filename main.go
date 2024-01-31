package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type Film struct {
    Title string
    Director string
}

func main() {
    fmt.Println("Running htmx server...")

    h1 := func (w http.ResponseWriter, r *http.Request) {

        films := map[string][]Film{
            "Films": {
                { Title: "The Godfather", Director: "Francis Ford Coppola" },
                { Title: "Blade Runner", Director: "Ridley Scott" },
                { Title: "The Thing", Director: "John Carpenter" },
            },
        }

        tmpl := template.Must(template.ParseFiles("index.html"))
        tmpl.Execute(w, films)
    }

    http.HandleFunc("/", h1)

    log.Fatal(http.ListenAndServe(":8000", nil))
}
