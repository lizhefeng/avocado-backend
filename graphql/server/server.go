package main

import (
	"context"
	"github.com/lizhefeng/avocado-backend/keystore"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/lizhefeng/avocado-backend/graphql"
	"github.com/lizhefeng/avocado-backend/sql"
)

const (
	defaultPort = "8080"
	dbPath      = "avocado.db"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	store := sql.NewSQLite3(dbPath)

	ks := keystore.NewKeyStore(store)
	if err := ks.RegisterDefaultProtocols(); err != nil {
		log.Fatal("Failed to register default protocols")
	}

	ctx := context.Background()
	if err := ks.Start(ctx); err != nil {
		log.Fatal("Failed to start the keystore: ", err)
	}

	defer func() {
		if err := ks.Stop(ctx); err != nil {
			log.Fatal("Failed to stop the keystore: ", err)
		}
	}()

	http.Handle("/", graphqlHandler(handler.Playground("GraphQL playground", "/query")))
	http.Handle("/query", graphqlHandler(handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{KeyStore: ks}}))))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func graphqlHandler(playgroundHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		playgroundHandler.ServeHTTP(w, r)
	})
}
