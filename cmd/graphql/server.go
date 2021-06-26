package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sy-software/minerva-owl/cmd/graphql/graph"
	"github.com/sy-software/minerva-owl/cmd/graphql/graph/generated"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/service"
	"github.com/sy-software/minerva-owl/internal/handlers"
	"github.com/sy-software/minerva-owl/internal/repositories"
)

const defaultConfigFile = "./config.json"

func main() {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = defaultConfigFile
	}

	config := domain.LoadConfiguration(configFile)
	cassandra := repositories.GetCassandra(config.CassandraDB)
	repo := repositories.NewOrgRepo(cassandra)

	defer cassandra.Close()

	service := service.NewOrgService(repo)
	handlerInstance := handlers.NewOrgGraphqlHandler(*service)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Handler: *handlerInstance,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", config.Port)
	log.Fatal(http.ListenAndServe(config.Host+":"+config.Port, nil))
}
