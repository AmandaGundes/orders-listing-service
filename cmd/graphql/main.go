package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/agundes/Projects/go/orders-listing-service/graph"
	"github.com/agundes/Projects/go/orders-listing-service/graph/generated"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	http.Handle("/query", srv)
	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))

	log.Println("GraphQL service listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
