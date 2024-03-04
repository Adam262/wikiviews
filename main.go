package main

import (
	"log"
	"net/http"
	"wikiviews/internal/pageviews"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	// Enable logger
	e.Use(middleware.Logger())
	// Rate limit to 20 rps
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	pageviewsHandler := pageviews.NewPageviewsHandler()
	e.GET("/healthcheck", healthcheck)
	e.GET("/pageviews", pageviewsHandler.List)

	log.Fatal(e.Start(":8080"))
}

// Healthcheck probe
// May be used for Kubernetes liveness and readiness probes
func healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
