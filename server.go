package main

import (
	"log"
	"net/http"
	"os"

	"scheduler/db"
	"scheduler/graph"
	"scheduler/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

/*
TODO
DB error handling
Implement some kind of authentication
Learn algorithmics
Documentation
Expand functionality
*/

/*
BUGS
when you try to create an already existing user postgre returns the id of it instead of returning an error
*/

const defaultPort = "8080"

func main() {
	err := db.DatabaseSetUp()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		db.Connection.Close()
		db.CancelFunc()
		return
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
