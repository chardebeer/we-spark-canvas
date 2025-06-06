package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/chardebeer/we-spark-canvas/graph"
	"github.com/vektah/gqlparser/v2/ast"
  "github.com/chardebeer/we-spark-canvas/graph/generated"
  "github.com/chardebeer/we-spark-canvas/server/storage"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
 // Initialize Postgres connection
  db := storage.NewPostgresDB()
  resolver := &graph.Resolver{DB: db}

  // Create GraphQL server
  srv := handler.NewDefaultServer(
    generated.NewExecutableSchema(generated.Config{Resolvers: resolver}),
  )

  // Serve Playground at “/”
  http.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
  http.Handle("/graphql", srv)

  log.Printf("🚀 We Spark Canvas API running at http://localhost:%s/", port)
  log.Fatal(http.ListenAndServe(":"+port, nil))
}