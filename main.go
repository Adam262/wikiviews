package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/ping", pong)

	log.Fatal(e.Start(":8080"))
}

func pong(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
