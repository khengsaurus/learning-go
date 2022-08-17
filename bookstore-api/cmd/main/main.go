package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	// "github.com/jinzhu/gorm"

	// "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/khengsaurus/go-tutorials/bookstore-api/pkg/routes"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)
	http.Handle("/", r)
	fmt.Printf("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
