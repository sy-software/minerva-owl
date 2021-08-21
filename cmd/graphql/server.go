package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	minervaLog "github.com/sy-software/minerva-go-utils/log"
	"github.com/sy-software/minerva-owl/cmd/graphql/graph"
	"github.com/sy-software/minerva-owl/cmd/graphql/graph/generated"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/service"
	"github.com/sy-software/minerva-owl/internal/handlers"
	"github.com/sy-software/minerva-owl/internal/repositories"
	"github.com/sy-software/minerva-owl/internal/repositories/mongodb"
)

// Defining the Graphql handler
func graphqlHandler(config *domain.Config, resolver *graph.Resolver) gin.HandlerFunc {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		// TODO: Implement bug tracker
		log.Error().Msgf("Unexpected error: %v", err)
		return errors.New("internal server error")
	})

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	minervaLog.ConfigureLogger(minervaLog.LogLevel(os.Getenv("LOG_LEVEL")), os.Getenv("CONSOLE_OUTPUT") != "")
	log.Info().Msg("Starting server")

	configRepo := repositories.ConfigRepo{}
	config := configRepo.Get()

	mdbInstance, err := mongodb.GetMongoDB(config.MongoDBConfig)

	if err != nil {
		log.Error().Stack().Err(err).Msg("Can't initialize Mongo DB")
		os.Exit(1)
	}

	repo, err := mongodb.NewMongoRepo(mdbInstance, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't start Mongo DB Repo")
		os.Exit(1)
	}

	defer mdbInstance.Close()

	orgService := service.NewOrgService(repo, config)
	usrService := service.NewUserService(repo, config)
	orgHandler := handlers.NewOrgGraphqlHandler(*orgService)
	usrHandler := handlers.NewUserGraphqlHandler(*usrService)

	r := gin.Default()

	r.POST("/query", graphqlHandler(&config, &graph.Resolver{
		OrgHandler: *orgHandler,
		UsrHandler: *usrHandler,
	}))
	r.GET("/", playgroundHandler())

	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	srv := &http.Server{
		Addr:    address,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Info().Msgf("connect to http://%s/ for GraphQL playground", address)
		err := srv.ListenAndServe()

		if err != http.ErrServerClosed {
			log.Panic().Err(err).Msg("Server crashed")
		} else {
			log.Info().Msg("Server closed")
		}
	}()

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Info().Msg("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Panic().Err(err).Msg("Server forced to shutdown")
	}
}
