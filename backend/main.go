package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/shusann01116/Chatting/backend/app/config"
	"github.com/shusann01116/Chatting/backend/app/handler"
	"github.com/shusann01116/Chatting/backend/app/loader"
	"github.com/shusann01116/Chatting/backend/app/repository"
	"github.com/shusann01116/Chatting/backend/app/resolver"
	"github.com/shusann01116/Chatting/backend/app/schema"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	ctx := context.TODO()
	dbConf := config.GetDBConfig()
	// url := fmt.Sprintf("postgresql://%s:%s@%s:%s", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port)
	url := fmt.Sprintf("postgresql://%s:%s@%s:5432", dbConf.User, dbConf.Password, dbConf.Host)
	db := repository.NewDataBase(ctx, url)

	root, err := resolver.NewRoot(db)
	if err != nil {
		log.Fatal(err)
	}

	s, err := schema.String()
	if err != nil {
		log.Fatal(err)
	}

	h := handler.GraphQL{
		Schema:  graphql.MustParseSchema(s, root),
		Loaders: loader.Initialize(db),
	}

	// Register handlers to routes.
	mux := http.NewServeMux()
	mux.Handle("/", handler.GraphiQL{})
	mux.Handle("/graphql/", h)
	mux.Handle("/graphql", h) // Register without a trailing slash to avoid redirect.

	// Configure the HTTP server.
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    http.DefaultMaxHeaderBytes,
	}

	// Begin listening for requests.
	log.Printf("Listening for requests on %s", srv.Addr)

	if err = srv.ListenAndServe(); err != nil {
		log.Printf("server.ListenAndServe: %v", err)
	}

	log.Println("Shutting down server...")
}
