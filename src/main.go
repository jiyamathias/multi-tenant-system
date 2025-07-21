// Package main is the entry point of the application. It sets up the whole application and make it ready to use
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"

	"codematic/controller"
	"codematic/handler"
	"codematic/pkg/environment"
	"codematic/pkg/helper"
	"codematic/pkg/middleware"
	codematicStorage "codematic/storage"
)

//	@title			Multi-Tenant API
//	@version		1.0
//	@description	This is the API for the Multi-Tenant service
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host						localhost:5002
// @BasePath					/api/v1
// @schemes					https
// @query.collection.format	multi
// @securityDefinitions.basic	BasicAuth
func main() {
	// set global application timezone
	_ = os.Setenv("TZ", "Africa/Lagos")
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	applicationLogger := logger.With().Str(helper.LogStrKeyModule, "app").Logger()
	r := gin.New()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowAllOrigins = true

	// Limit to 10 requests per minutes
	rate, _ := limiter.NewRateFromFormatted("10-M")
	store := memory.NewStore()
	rateLimiterMiddleware := ginlimiter.NewMiddleware(limiter.New(store, rate))
	r.Use(rateLimiterMiddleware)

	r.Use(cors.New(corsConfig), gin.Recovery())
	r.Use(ginzerolog.Logger("rest"))
	r.Use(GinContextToContextMiddleware())
	r.Use(requestid.New())

	env, err := environment.New()
	if err != nil {
		applicationLogger.Fatal().Err(err)
		panic(err) // panic - this service should not start up
	}

	storage := codematicStorage.New(logger, env)
	defer storage.Close()
	// run automigration
	if err := storage.AutoMigrate(); err != nil {
		applicationLogger.Fatal().Err(err)
		panic(err) // panic - this service should not start up either
	}

	// initialize the app
	r.Use(GinContextToContextMiddleware())
	// init our custom middleware
	newMiddleware := middleware.NewMiddleware(logger, *env, storage)
	application := controller.New(logger, storage, newMiddleware)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "go ninjas are alive",
			"rest":    true,
		})
	})

	logger.Info().Msg("starting - health check is up...")

	h := handler.New(logger, env, r, *application)
	h.Build()

	port := env.Get("SERVER_PORT")
	if strings.EqualFold(port, "") {
		port = "5002"
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			applicationLogger.Fatal().Msgf("listen: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	applicationLogger.Info().Msgf("Shutdown Server ... %v", <-quit)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		applicationLogger.Fatal().Msgf("Server Shutdown: %v", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		applicationLogger.Info().Msgf("timeout of 5 seconds.")
	default:
	}

	applicationLogger.Info().Msgf("Server exiting")
}

// GinContextToContextMiddleware middleware for gin context
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), helper.GinContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
