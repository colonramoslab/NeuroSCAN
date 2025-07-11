package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"neuroscan/internal/cache"
	"neuroscan/internal/database"
	"neuroscan/internal/handler"
	"neuroscan/internal/repository"
	"neuroscan/internal/router"
	"neuroscan/internal/service"
	"neuroscan/pkg/logging"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type WebCmd struct {
	Port string `optional:"" help:"Port to run the web server on." short:"p"`
}

func (cmd *WebCmd) Run(ctx *context.Context) error {
	logger := logging.NewLoggerFromEnv()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("🤯 failed to load environment variables")
		return err
	}

	cntx := logging.WithLogger(*ctx, logger)

	port := "8080"

	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	}

	if cmd.Port != "" {
		port = cmd.Port
	}

	// if no port is set, default to 8080
	if port == "" {
		port = "8080"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	db, err := database.NewFromEnv(cntx)
	if err != nil {
		logger.Fatal().Err(err).Msg("🤯 failed to connect to database")
		return err
	}

	cache, err := cache.NewCache(cntx)
	if err != nil {
		logger.Fatal().Err(err).Msg("🤯 failed to connect to cache")
		return fmt.Errorf("failed to connect to cache: %w", err)
	}

	defer db.Close(cntx)

	e := echo.New()
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))

	if appEnv == "production" {
		e.HideBanner = true
		e.HidePort = true
		e.Pre(middleware.HTTPSRedirect())

		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
		}))

		// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(100))))

		e.Use(middleware.Recover())
		e.Use(middleware.Secure())
	}

	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// config := middleware.RateLimiterConfig{
	// 	Skipper: middleware.DefaultSkipper,
	// 	Store: middleware.NewRateLimiterMemoryStoreWithConfig(
	// 		middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(100), Burst: 200, ExpiresIn: 3 * time.Minute},
	// 	),
	// 	IdentifierExtractor: func(ctx echo.Context) (string, error) {
	// 		id := ctx.RealIP()
	// 		return id, nil
	// 	},
	// 	ErrorHandler: func(context echo.Context, err error) error {
	// 		return context.JSON(http.StatusForbidden, nil)
	// 	},
	// 	DenyHandler: func(context echo.Context, identifier string, err error) error {
	// 		return context.JSON(http.StatusTooManyRequests, nil)
	// 	},
	// }

	// e.Use(middleware.RateLimiterWithConfig(config))

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "the request has timed out",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			_ = c.String(http.StatusRequestTimeout, "the request has timed out")
		},
		Timeout: 30 * time.Second,
	}))

	e.Static("/files", os.Getenv("APP_GLTF_DIR"))
	e.Static("/", os.Getenv("APP_FRONTEND_DIR"))

	neuronRepo := repository.NewPostgresNeuronRepository(db.Pool, cache)
	neuronService := service.NewNeuronService(neuronRepo)
	neuronHandler := handler.NewNeuronHandler(neuronService)

	contactRepo := repository.NewPostgresContactRepository(db.Pool, cache)
	contactService := service.NewContactService(contactRepo)
	contactHandler := handler.NewContactHandler(contactService)

	synapseRepo := repository.NewPostgresSynapseRepository(db.Pool, cache)
	synapseService := service.NewSynapseService(synapseRepo)
	synapseHandler := handler.NewSynapseHandler(synapseService)

	cphateRepo := repository.NewPostgresCphateRepository(db.Pool, cache)
	cphateService := service.NewCphateService(cphateRepo)
	cphateHandler := handler.NewCphateHandler(cphateService)

	nerveringRepo := repository.NewPostgresNerveRingRepository(db.Pool, cache)
	nerveringService := service.NewNerveRingService(nerveringRepo)
	nerveringHandler := handler.NewNerveRingHandler(nerveringService)

	scaleRepo := repository.NewPostgresScaleRepository(db.Pool, cache)
	scaleService := service.NewScaleService(scaleRepo)
	scaleHandler := handler.NewScaleHandler(scaleService)

	promoterRepo := repository.NewPostgresPromoterRepository(db.Pool, cache)
	promoterService := service.NewPromoterService(promoterRepo)
	promoterHandler := handler.NewPromoterHandler(promoterService)

	devStageRepo := repository.NewPostgresDevelopmentalStageRepository(db.Pool, cache)
	devStageService := service.NewDevelopmentalStageService(devStageRepo)
	devStageHandler := handler.NewDevelopmentalStageHandler(devStageService)

	videoRepo := repository.NewPostgresVideoRepository(db.Pool, cache)
	videoService := service.NewVideoRepository(videoRepo)
	videoHandler := handler.NewVideoHandler(videoService)

	e = router.NewRouter(e, neuronHandler, contactHandler, synapseHandler, cphateHandler, nerveringHandler, scaleHandler, promoterHandler, devStageHandler, videoHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))

	return nil
}
