package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	Item struct {
		Project     string `json:"project"`
		Article     string `json:"article"`
		Granularity string `json:"granularity"`
		Timestamp   string `json:"timestamp"`
		Views       int32  `json:"views"`
	}
	ResponseData struct {
		Items []Item `json:"items"`
	}
)

const (
	baseurl = "https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article"
)

func main() {
	e := echo.New()

	e.GET("/ping", pong)
	e.GET("/pageviews", pageviews)

	log.Fatal(e.Start(":8080"))
}

func pageviews(c echo.Context) (err error) {
	url := fmt.Sprintf("%s/en.wikipedia.org/all-access/all-agents/Orca/monthly/20240201/20240229", baseurl)

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	// Send a GET request to Wikipedia API
	response, err := client.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Unmarshal the JSON response into a struct
	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Print the message from the JSON response
	return c.JSON(http.StatusOK, responseData.Items)
}

func pong(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
