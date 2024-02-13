package handlers

import (
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

    search_query := r.URL.Query().Get("q")

    data := TemplateData{
        Products: nil,
        SearchString: search_query,
    }

    tmpl.Execute(w, data)
}

func ProductList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    tmpl := template.Must(template.ParseFiles("product-list.tmpl.html"))

    no_products_html :=
        `<div id="product-list">
            <p class="mt-5 mx-auto" style="width: max-content">
                There are no products to display
            </p>
        </div>`

    db, err := tools.ConnectDatabase()
    if err != nil {
        log.Printf("Error connecting to db: %v", err)
        return
    }
    log.Print("Successfully connected to db")

    // Get search value
    search := r.FormValue("search")
    if search != "" {
        w.Header().Add("HX-Replace-URL", "?q=" + url.QueryEscape(search))
    } else {
        w.Header().Add("HX-Replace-URL", "/")
    }

    // Escape characters and use `*` instead of `%`
    search = strings.Replace(search, "%", "\\%", -1)
    search = strings.Replace(search, "_", "\\_", -1)
    search = strings.Replace(search, "*", "%", -1)

    search, queries := parseSearchQuery(search) // FIXME: `queries` is a confusing name

    // Initialize the query, where_clause, and conditions
    query := fmt.Sprintf("SELECT name, price, available FROM product")
    where_clause := ""
    conditions := []string{}

    // Add a name condition if the search is not empty
    if search != "" {
        conditions = append(
            conditions, 
            fmt.Sprintf("name ILIKE '%s'", search),
        )
    }

    // If there are other queries, add them as conditions
    if len(queries) > 0 {
        for _, col_val := range queries {
            c := fmt.Sprintf("%s::varchar ILIKE '%s'", col_val[0], col_val[1]);
            conditions = append(conditions, c)
            log.Printf("Adding condition: %s", c)
        }
    }

    // If there are conditions, build a where clause from them
    if len(conditions) > 0 {
        where_clause = " WHERE " + strings.Join(conditions, " AND ")
    }
    // Add the where clause to the query
    query += where_clause

    log.Printf("Querying the database: `%s`", query)
    rows, err := db.Query(query)
    if err != nil {
        log.Printf("Error fetching products with query `%s`: %v", query, err)
        fmt.Fprint(w, no_products_html)
        return
    }
    defer rows.Close()

    products := []dbcontrol.Product{}
    var name string
    var price float64
    var available bool

    for rows.Next() {
        err := rows.Scan(&name, &price, &available)
        if err != nil {
            log.Printf("Error converting row to product: %v\n", err)
        } else {
            products = append(products, dbcontrol.Product{ Name: name, Price: price, Available: available })
        }
    }

    if len(products) > 0 {
        data := map[string][]dbcontrol.Product{
            "Products": products,
        }

        tmpl.Execute(w, data)
    } else {
        fmt.Fprint(w, no_products_html)
    }
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
    product := dbcontrol.Product{
        Name: name,
        Price: price,
        Available: available,
    }
    dbcontrol.InsertProduct(db, product)

    // Return product as html fragment
    tmpl := template.Must(template.ParseFiles("product-list.tmpl.html"))
    tmpl.ExecuteTemplate(w,
        "product-list-element",
        product,
    )
}

func LoadDummyDataHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    db, err := tools.ConnectDatabase()
    defer db.Close()

    if err != nil {
        log.Printf("Error connecting to db: %v", err)
        return
    }
    products := tools.LoadDummyData(db)
    log.Printf("Random products: %v", products)

    tmpl := template.Must(template.ParseFiles("product-list.tmpl.html"))
    for i, p := range products {
        log.Printf("Product %d: %v", i, p)
        tmpl.ExecuteTemplate(w,
            "product-list-element",
            p,
        )
    }
}

func parseSearchQuery(search string) (return_string string, queries [][2]string) {
    queries = make([][2]string, 0)
    split := strings.Split(search, " ")

    filtered := tools.Filter(split, func(s string) bool {
        return !strings.Contains(s, ":")
    })

    return_string = strings.Join(filtered, " ")

    query_strings := tools.Filter(split, func(s string) bool {
        return strings.Contains(s, ":")
    })

    for _, s := range query_strings {
        key_val := strings.SplitN(s, ":", 2)
        queries = append(
            queries,
            [2]string{key_val[0], key_val[1]},
        )
    }

    log.Printf("Parsed search string: `%s`", search)
    log.Printf("Return string: `%s`", return_string)
    log.Printf("Queries: %v", queries)

    return
}
