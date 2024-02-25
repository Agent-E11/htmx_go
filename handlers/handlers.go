package handlers

import (
    "errors"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "text/template"

    "github.com/agent-e11/htmx_go/dbcontrol"
    "github.com/agent-e11/htmx_go/tools"
    "github.com/julienschmidt/httprouter"
)

type TemplateData struct {
    Products []dbcontrol.Product
    SearchString string
}

func HomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    tmpl := template.Must(template.ParseFiles("index.tmpl.html"))

    searchQuery := r.URL.Query().Get("q")

    data := TemplateData{
        Products: nil,
        SearchString: searchQuery,
    }

    tmpl.Execute(w, data)
}

func ProductList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    tmpl := template.Must(template.ParseFiles("product-list.tmpl.html"))

    noProductsHtml :=
        `<div id="product-list">
            <p class="mt-5 mx-auto" style="width: max-content">
                There are no products to display
            </p>
        </div>`

    db, err := dbcontrol.ConnectDatabase()
    if err != nil {
        log.Printf("Error connecting to db: %v", err)
        return
    }
    defer db.Close()
    log.Print("Successfully connected to db")

    // Get search value
    search := r.FormValue("search")
    if search != "" {
        w.Header().Add("HX-Replace-URL", "?q=" + url.QueryEscape(search))
    } else {
        w.Header().Add("HX-Replace-URL", "/")
    }

    // Trim whitespace and collapse multiple spaces
    search = strings.Join(strings.Fields(search), " ")

    // Escape characters and use `*` instead of `%`
    search = strings.Replace(search, "%", "\\%", -1)
    search = strings.Replace(search, "_", "\\_", -1)
    search = strings.Replace(search, "*", "%", -1)

    search, searchParams := parseSearchQuery(search)

    // Initialize the query, whereClause, and conditions
    query := fmt.Sprintf("SELECT id, name, price, available FROM product")
    whereClause := ""
    conditions := []string{}

    // Add a name condition if the search is not empty
    if search != "" {
        conditions = append(
            conditions, 
            fmt.Sprintf("name ILIKE '%s'", search),
        )
    }

    // If there are other searchParams, add them as conditions
    if len(searchParams) > 0 {
        for _, col_val := range searchParams {
            c := fmt.Sprintf("%s::varchar ILIKE '%s'", col_val[0], col_val[1]);
            conditions = append(conditions, c)
            log.Printf("Adding condition: %s", c)
        }
    }

    // If there are conditions, build a where clause from them
    if len(conditions) > 0 {
        whereClause = " WHERE " + strings.Join(conditions, " AND ")
    }
    // Add the where clause to the query
    query += whereClause

    // Add an order by clause
    query += " ORDER BY name, price"

    log.Printf("Querying the database: `%s`", query)
    rows, err := db.Query(query)
    if err != nil {
        log.Printf("Error fetching products with query `%s`: %v", query, err)
        fmt.Fprint(w, noProductsHtml)
        return
    }
    defer rows.Close()

    products := []dbcontrol.Product{}
    var id uint
    var name string
    var price float64
    var available bool

    for rows.Next() {
        err := rows.Scan(&id, &name, &price, &available)
        if err != nil {
            log.Printf("Error converting row to product: %v\n", err)
        } else {
            products = append(
                products,
                dbcontrol.Product{
                    Id: id,
                    Name: name,
                    Price: price,
                    Available: available,
                },
            )
        }
    }

    if len(products) > 0 {
        data := map[string][]dbcontrol.Product{
            "Products": products,
        }

        tmpl.Execute(w, data)
    } else {
        fmt.Fprint(w, noProductsHtml)
    }
}

// Add a product to the database
func AddProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    tmpl := template.Must(template.ParseFiles("form.tmpl.html"))
    // Connect to database
    db, err := dbcontrol.ConnectDatabase()
    if err != nil {
        log.Println("Error connecting to database")
        return
    }
    defer db.Close()

    // Get values from form
    name := r.PostFormValue("name")
    priceStr := r.PostFormValue("price")
    availableStr := r.PostFormValue("available")

    // Trim whitespace and collapse multiple spaces
    name = strings.Join(strings.Fields(name), " ")

    var nameErr error
    if name == "" {
        nameErr = errors.New("Name field cannot be empty")
        log.Println("Error: cannot parse empty fields")
    }
    var priceErr error
    var price float64
    if priceStr == "" {
        priceErr = errors.New("Price field cannot be empty")
        log.Println("Error: cannot parse empty fields")
    } else {
        // Convert to data types
        price, priceErr = strconv.ParseFloat(priceStr, 64)
        if priceErr != nil {
            priceErr = errors.New("Price must be a valid number")
            log.Printf("Error parsing price: %v\n", err)
        }
    }
    
    // If the string is not "", then the product is available
    available := availableStr != ""

    log.Println("name:", name)
    log.Println("priceStr:", priceStr)
    log.Println("availableStr:", availableStr)
    log.Println("nameErr:", nameErr)
    log.Println("price:", price)
    log.Println("priceErr:", priceErr)
    log.Println("available:", available)

    // Create product and insert into database
    if priceErr == nil && nameErr == nil {
        product := dbcontrol.Product{
            Name: name,
            Price: price,
            Available: available,
        }
        dbcontrol.InsertProduct(db, product)

        // Tell the client to re-search the product list
        w.Header().Add("HX-Trigger", "research-products")
    }

    tmpl.Execute(w, struct{
        Name string
        NameErr error
        Price float64
        PriceErr error
        Available bool
    }{
        Name: name,
        NameErr: nameErr,
        Price: price,
        PriceErr: priceErr,
        Available: available,
    })
}

func DeleteById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    db, err := dbcontrol.ConnectDatabase()
    if err != nil {
        log.Printf("Error connecting to db: %v", err)
        fmt.Fprintf(w, "Error connecting to db: %v", err)
        return
    }
    defer db.Close()

    id, err := strconv.Atoi(ps.ByName("id"))
    if err != nil {
        log.Printf("Error parsing id")
        fmt.Fprintln(w, "Error parsing id")
        return
    }
    
    query := fmt.Sprintf("DELETE FROM product WHERE id = %d", id)
    log.Printf("Deleting product with id `%d` with query: %s", id, query)

    _, err = db.Exec(query)

    w.Header().Add("HX-Trigger", "research-products")
}

func DeleteAllData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    db, err := dbcontrol.ConnectDatabase()
    if err != nil {
        log.Printf("Error connecting to db: %v", err)
        return
    }
    defer db.Close()

    query := "DELETE FROM product"
    log.Println("Deleting all data from database")

    _, err = db.Exec(query)

    w.Header().Add("HX-Trigger", "research-products")
}

func LoadDummyDataHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    db, err := dbcontrol.ConnectDatabase()
    if err != nil {
        log.Printf("Error connecting to db: %v", err)
        return
    }
    defer db.Close()

    products := tools.LoadDummyData(db)
    log.Printf("Random products: %v", products)

    // Tell the client to re-search the product list
    w.Header().Add("HX-Trigger", "research-products")
}

func parseSearchQuery(search string) (newSearch string, params [][2]string) {
    params = make([][2]string, 0)
    split := strings.Split(search, " ")

    filtered := tools.Filter(split, func(s string) bool {
        return !strings.Contains(s, ":")
    })

    newSearch = strings.Join(filtered, " ")

    queryStrings := tools.Filter(split, func(s string) bool {
        return strings.Contains(s, ":")
    })

    for _, s := range queryStrings {
        key_val := strings.SplitN(s, ":", 2)
        params = append(
            params,
            [2]string{key_val[0], key_val[1]},
        )
    }

    log.Printf("Parsed search string: `%s`", search)
    log.Printf("Return string: `%s`", newSearch)
    log.Printf("Queries: %v", params)

    return
}
