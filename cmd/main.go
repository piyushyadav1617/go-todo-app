package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/piyushyadav1617/go-todo-app/internal/database"
	"github.com/piyushyadav1617/go-todo-app/internal/route"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if err := database.InitDB(connStr); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	if err := database.Migrate(); err != nil {
		log.Fatalf("%v", err)
	}

	route.TodoRoutes()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}
