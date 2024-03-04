package pageviews

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"wikiviews/internal/httpclient"
	"wikiviews/internal/paramformatter"
	"wikiviews/internal/paramvalidator"

	"github.com/labstack/echo/v4"
)

type (
	Item struct {
		Article   string `json:"article"`
		Timestamp string `json:"timestamp"`
		Views     int32  `json:"views"`
	}
	ResponseData struct {
		Items []Item `json:"items"`
	}

	PageviewsHandler struct{}
)

const (
	baseUrl   = "https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia.org/all-access/all-agents"
	userAgent = "WikiViews/1.0"
)

func (ph *PageviewsHandler) List(c echo.Context) (err error) {
	// Query escape all incoming article params
	article := url.QueryEscape(c.QueryParam("article"))

	// Validate title input
	tv := paramvalidator.NewTitleValidator()
	tvok, err := tv.Run(article)
	if !tvok {
		log.Println("error:", err)
		return c.JSON(http.StatusBadRequest, errorMessage((err)))
	}

	// Validate date input
	date := c.QueryParam("date")
	dv := paramvalidator.NewDateValidator()
	dvok, err := dv.Run(date)
	if !dvok {
		log.Println("error:", err)
		return c.JSON(http.StatusBadRequest, errorMessage((err)))
	}

	// Return start and end date params that Wikipedia API needs from single data input
	df := paramformatter.NewDateFormatter()
	monthStart, monthEnd, err := df.Run(date)
	if err != nil {
		log.Println("error:", err)
		return c.JSON(http.StatusBadRequest, errorMessage((err)))
	}

	url := fmt.Sprintf("%s/%s/monthly/%s/%s", baseUrl, article, monthStart, monthEnd)
	client := httpclient.NewHttpClient()

	// Create a new HTTP GET request with our User-Agent header
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("error:", err)
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	log.Printf("sending GET request to Wikipedia endpoint: %s\n", req.URL)

	// Send the request to the Wikipedia API
	response, err := client.Do(req)
	if err != nil {
		log.Println("response error:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("error reading response:", err)
		return
	}

	// Handle 404 error response code
	if response.StatusCode == http.StatusNotFound {
		pf := paramformatter.NewTitleFormatter()
		baseMessage := "error: query for article param: %s did not return any results. Consider titlizing article param as %s."

		if pf.IsSingleWord(article) {
			fullMessage := fmt.Sprintf(baseMessage, article, pf.Run(article, true))
			err = fmt.Errorf(fullMessage)
		} else if pf.IsMultiWord(article) {
			titleOptions := fmt.Sprintf("%s or %s", pf.Run(article, true), pf.Run(article, false))
			fullMessage := fmt.Sprintf(baseMessage, article, titleOptions)
			err = fmt.Errorf(fullMessage)
		} else {
			err = fmt.Errorf("no results found")
		}

		log.Println("error:", err)
		return c.JSON(http.StatusNotFound, errorMessage((err)))
	}

	// Unmarshal the JSON response into a struct
	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Println("error unmarshalling JSON:", err)
		return
	}

	// Set response headers
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)

	// Print the message from the JSON response
	return json.NewEncoder(c.Response()).Encode(responseData.Items)
}

func errorMessage(err error) map[string]string {
	return map[string]string{
		"error": err.Error(),
	}
}

func NewPageviewsHandler() *PageviewsHandler {
	return &PageviewsHandler{}
}
