package tools

import (
	"database/sql"
	"encoding/json"
	"os"

	mySql "github.com/agent-e11/htmx_go/sql"
)

func LoadDummyData(db *sql.DB, filename string) error {
    bytes, err := os.ReadFile(filename)
    if err != nil {
        return err
    }

    productJson := string(bytes)

    var products []mySql.Product

    err = json.Unmarshal([]byte(productJson), &products)
    if err != nil {
        return err
    }

    for _, p := range products {
        mySql.InsertProduct(db, p)
    }

    return nil
}
