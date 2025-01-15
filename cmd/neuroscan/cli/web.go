package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"neuroscan/internal/database"
	"neuroscan/internal/handler"
	"neuroscan/internal/logging"
	"neuroscan/internal/repository"
	"neuroscan/internal/router"
	"neuroscan/internal/service"


	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type WebCmd struct {}

func (cmd *WebCmd) Run(ctx *context.Context) error {
	logger := logging.NewLoggerFromEnv()

	cntx := logging.WithLogger(*ctx, logger)

	port := os.Getenv("PORT")
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

		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(100))))

		e.Use(middleware.Recover())
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(20), Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	e.Use(middleware.RateLimiterWithConfig(config))

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "the request has timed out",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			_ = c.String(http.StatusRequestTimeout, "the request has timed out")
		},
		Timeout: 30 * time.Second,
	}))


	neuronRepo := repository.NewPostgresNeuronRepository(db.Pool)
	neuronService := service.NewNeuronService(neuronRepo)
	neuronHandler := handler.NewNeuronHandler(neuronService)

	contactRepo := repository.NewPostgresContactRepository(db.Pool)
	contactService := service.NewContactService(contactRepo)
	contactHandler := handler.NewContactHandler(contactService)

	e = router.NewRouter(e, neuronHandler, contactHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))

	return nil
}
