package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

type Film struct {
    Title string
    Director string
}

func main() {
    fmt.Println("Running htmx server...")

    h1 := func(w http.ResponseWriter, r *http.Request) {

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

    h2 := func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(time.Second*1)
        title := r.PostFormValue("title")
        director := r.PostFormValue("director")

        tmpl := template.Must(template.ParseFiles("index.html"))
        tmpl.ExecuteTemplate(w,
            "film-list-element",
            Film{
                Title: title,
                Director: director,
            },
        )
    }

    http.HandleFunc("/", h1)
    http.HandleFunc("/add-film/", h2)

    log.Fatal(http.ListenAndServe(":8000", nil))
}
