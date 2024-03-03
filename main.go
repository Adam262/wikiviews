package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"wikiviews/internal/httpclient"
	"wikiviews/internal/paramformatter"
	"wikiviews/internal/paramvalidator"

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
	// Enable logger
	e.Use(middleware.Logger())
	// Rate limit to 20 rps
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	e.GET("/healthcheck", healthcheck)
	e.GET("/pageviews", pageviews)

	log.Fatal(e.Start(":8080"))
}

func pageviews(c echo.Context) (err error) {
	article := c.QueryParam("article")
	tv := paramvalidator.NewTitleValidator()
	// titlizedArticle := tf.Run(article)
	// if converted {
	// 	log.Printf("Converted article param %s to %s", article, titlizedArticle)
	// }
	ok, err := tv.Run(article)
	if !ok {
		log.Println("Error:", err)
		return c.JSON(http.StatusBadRequest, errorMessage((err)))
	}

	monthstart := c.QueryParam("monthstart")
	monthend := c.QueryParam("monthend")
	url := fmt.Sprintf("%s/%s/monthly/%s/%s", baseUrl, article, monthstart, monthend)

	client := httpclient.NewHttpClient()

	// Create a new HTTP GET request with our User-Agent header
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	req.Header.Set("User-Agent", userAgent)

	// Send the request to the Wikipedia API
	response, err := client.Do(req)
	if err != nil {
		log.Println("Response Error:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return
	}

	// Handle 404 error response code
	if response.StatusCode == http.StatusNotFound {
		pf := paramformatter.NewTitleFormatter()
		baseMessage := "error: query for article param: %s did not return any results. Consider titlizing article param as %s"

		if pf.IsSingleWord(article) {
			fullMessage := fmt.Sprintf(baseMessage, article, pf.Run(article, true))
			err = fmt.Errorf(fullMessage)
		} else if pf.IsMultiWord(article) {
			titleOptions := fmt.Sprintf("%s or %s", pf.Run(article, true), pf.Run(article, false))
			fullMessage := fmt.Sprintf(baseMessage, article, titleOptions)
			err = fmt.Errorf(fullMessage)
		} else {
			err = fmt.Errorf("No results found")
		}

		log.Println("Error:", err)
		return c.JSON(http.StatusNotFound, errorMessage((err)))
	}

	// Unmarshal the JSON response into a struct
	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return
	}

	// Set response headers
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	// Print the message from the JSON response
	return json.NewEncoder(c.Response()).Encode(responseData.Items)
}

// Healthcheck probe
// May be used for Kubernetes liveness and readiness probes
func healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func errorMessage(err error) map[string]string {
	return map[string]string{
		"error": err.Error(),
	}
}
