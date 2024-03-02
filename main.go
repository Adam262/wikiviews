package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"wikiviews/internal/httpclient"
	"wikiviews/internal/paramformatter"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// PrettyResponse struct {
	// 	OriginalArticle   string
	// 	FormattedfArticle string
	// 	Views             int32
	// 	Month             string
	// 	Year              int8
	// }
)

const (
	baseUrl   = "https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia.org/all-access/all-agents"
	userAgent = "WikiViews/1.0"
)

func main() {
	e := echo.New()
	// Rate limit to 20 rps
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	e.GET("/healthcheck", healthcheck)
	e.GET("/pageviews", pageviews)

	log.Fatal(e.Start(":8080"))
}

func pageviews(c echo.Context) (err error) {
	article := c.QueryParam("article")
	tf := paramformatter.NewTitleFormatter()
	titlizedArticle := tf.Run(article)
	// if converted {
	// 	log.Printf("Converted article param %s to %s", article, titlizedArticle)
	// }

	monthstart := c.QueryParam("monthstart")
	monthend := c.QueryParam("monthend")
	url := fmt.Sprintf("%s/%s/monthly/%s/%s", baseUrl, titlizedArticle, monthstart, monthend)

	client := httpclient.NewHttpClient()

	// Create a new HTTP GET request with our User-Agent header
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)

	// Send the request to the Wikipedia API
	response, err := client.Do(req)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return
	}

	// Unmarshal the JSON response into a struct
	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return
	}

	// Print the message from the JSON response
	return c.JSON(http.StatusOK, responseData.Items)
}

// Healthcheck probe
// May be used for Kubernetes liveness and readiness probes
func healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
