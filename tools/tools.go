package tools

import (
    "database/sql"

    "github.com/agent-e11/htmx_go/dbcontrol"
    "github.com/brianvoe/gofakeit/v6"
)

// Load data from file (encoded as json) into database
func LoadDummyData(db *sql.DB) []dbcontrol.Product {
    var products []dbcontrol.Product

    for i := 0; i < 5; i++ {
        products = append(products, dbcontrol.Product{
            Name: gofakeit.ProductName(),
            Price: gofakeit.Price(0.0, 200.0),
            Available: gofakeit.Bool(),
        })
    }

    for _, p := range products {
        dbcontrol.InsertProduct(db, p)
    }

    return products
}

func Filter[T any](slice []T, test func(T) bool) []T {
    var ret []T
    for _, s := range slice {
        if test(s) {
            ret = append(ret, s)
        }
    }

    return ret
}
