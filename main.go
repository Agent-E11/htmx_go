package main

import (
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"

    "github.com/agent-e11/htmx_go/handlers"
    "github.com/agent-e11/htmx_go/dbcontrol"
)

func main() {
    db, err := dbcontrol.ConnectDatabase()
    defer db.Close()

    if err != nil {
        log.Fatalf("Error opening db: %v\n", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatalf("Error pinging db: %v\n", err)
    }

    log.Print("Connected to database...")

    dbcontrol.CreateProductTable(db)

    // HTTP Server
    log.Print("Running htmx server...")

    router := httprouter.New()

    router.GET("/", handlers.HomePage)
    router.POST("/add-product/", handlers.AddProduct)
    router.POST("/load-dummy-data/", handlers.LoadDummyDataHandler)
    router.POST("/product-list/", handlers.ProductList)
    router.GET("/delete-by-id/:id", handlers.DeleteById) // TODO: This should be a DELETE request
    //router.GET("/delete-all-data/", deleteAllDataHandler)

    log.Fatal(http.ListenAndServe(":8000", router))
}

//func deleteAllDataHandler(w http.ResponseWriter, r *http.Request) {
//    db, err := tools.ConnectDatabase()
//    defer db.Close()
//
//    if err != nil {
//        return
//    }
//}
