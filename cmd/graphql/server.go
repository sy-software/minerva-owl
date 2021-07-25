package main

import (
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sy-software/minerva-owl/cmd/graphql/graph"
	"github.com/sy-software/minerva-owl/cmd/graphql/graph/generated"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/service"
	"github.com/sy-software/minerva-owl/internal/handlers"
	"github.com/sy-software/minerva-owl/internal/repositories/mongodb"
)

const defaultConfigFile = "./config.json"

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting server")

	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = defaultConfigFile
	}

	config := domain.LoadConfiguration(configFile)
	mdbInstance, err := mongodb.GetMongoDB(config.MongoDBConfig)

	if err != nil {
		log.Error().Stack().Err(err).Msg("Can't initialize Cassandra DB:")
		os.Exit(1)
	}

	repo, err := mongodb.NewMongoRepo(mdbInstance, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't start Cassandra DB Repo:")
		os.Exit(1)
	}

	defer mdbInstance.Close()

	orgService := service.NewOrgService(repo, config)
	usrServce := service.NewUserService(repo, config)
	orgHandler := handlers.NewOrgGraphqlHandler(*orgService)
	usrHanlder := handlers.NewUserGraphqlHandler(*usrServce)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		OrgHandler: *orgHandler,
		UsrHandler: *usrHanlder,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Info().Msgf("connect to http://localhost:%s/ for GraphQL playground", config.Port)
	log.Fatal().
		Err(http.ListenAndServe(config.Host+":"+config.Port, nil)).
		Msg("Can't start server")
}
