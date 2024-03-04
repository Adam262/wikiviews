# WikiViews

## Overview

WikiViews is a simple Golang server with a single JSON endpoint, `/pageviews`. Its responsibility is to respond to user queries for monthly pageview data for English-language Wikipedia articles. Although users can alternatively query the Wikipedia API directly, WikiViews provides several enhancements such as a simplied interface and param validation.

## Dependencies

WikiViews is Dockerized. Its only dependency for running locally as a Docker container is [Docker Desktop](https://www.docker.com/products/docker-desktop/).

## Getting Started

### Via Docker Compose

To run the project locally, git clone this repo and run the below command from within the repo:

```bash
docker-compose up -d 
docker-compose logs -f
```

To stop the project

```bash
docker-compose down
```

### Via local server

WikiViews may also be run locally without Docker for faster iteration. Local development takes advantage of the [air Go package](https://github.com/cosmtrek/air) for live reload.

```bash
go install github.com/cosmtrek/air@latest
make run
```

### Running tests

Run unit tests via:

```bash
make test
```

### Vetting and package management

Other Make targets are exposed for local development

```bash
# Run Golang compile checks
make vet

# Reconcile modules and vendorize dependencies
make mod
```

## API

### /healthcheck

This method is a simple health check. It may be used for Kubernetes liveness and readiness probes.

```
curl -X GET http://localhost:8080/healthcheck

$ pong
```

### /pageviews

This endpoint accepts JSON queries to the [Wikipedia Pageviews REST API](https://wikimedia.org/api/rest_v1/#/Pageviews%20data). It returns a JSON-ified list of response objects, containing data as the article name, time period and pageview count.

#### Params

##### article (string)

The title of the Wikipedia article. It must follow the [naming conventions](https://en.wikipedia.org/wiki/Wikipedia:Naming_conventions_(technical_restrictions)) defined by Wikipedia, which can be summarized as:

* The title must begin with a capital lettter. Any additional words in the title may be either capital or lower-case
* A space between words must be entered as a single underscore
* These characters are forbidden anywhere in the article title: `# < > [ ] { } |`

##### date (int8)

The target year and month to query, expressed as:

```bash
# e.g. 202311 is November, 2023
YYYYMMDD
```

This format is an optimization of the underlying Wikipedia enpoint, which expects separate params for the start and end of the monthly period. That is, to express `November 2023`, the user would need two path parameters in their query:

* start -- `20241101`
* end -- `20241130`

#### Sample Request and Response

```bash
curl -X GET localhost:8080/pageviews\?article\=michael_phelps\&date=20240201
```

The response is a JSON-ified list of response objects, containing data as the article name, time period and pageview count.

```bash
[{"article":"Michael_Phelps","timestamp":"2024020100","views":125860}]
```

#### Endpoint Design Decisions

I made several design decisions in V1 of the endpoint for the sake of simplifying the interface, validating params and trying to mitigate the brittleness of the underlying endpoint

##### English-language only

I made the decision to only query on English-language articles. This decision had two benefits:

* Removed an additional param -- *project* -- that users would otherwise need to pass in
* Simplied the regex for validating article titles by removing the need to deal with non-English characters

##### Hard-code other params

I hard-coded three other params (from the Wikipedia endpoint) to simplify the user interface:

* *access*. This param filters by page access method, e.g.: *desktop*, *mobile-app* or *mobile-web*. I hard-coded to *all-access*.
* *agent*. This param filters by page agent, e.g.: *user*, *automated* or *spider*. I hard-coded to *all-agents*.
* *granularity*. This param sets the time unit for the response data, e.g.: *daily* or *monthly*. I hard-coded to *monthly*.

##### Validations

There are several simple validation on the article param:

* empty article
* article with a forbidden character
* article beginning with a lower-case letter

Another common case is an article that returns a 404 from the Wikipedia endpoint, because their API cannot find any references to the article.

In these cases, I added validation that suggests a title case article. That is, if the user entered `MICHAEL_PHELPS`, the validation would suggest `Michael_Phelps` or `Michael_phelps`. I include both because both may be valid (and handled by redirection in the Wikipedia UI) and it is impossible to know which one the user intended. For example:

* `Michael_Phelps` is the canonical Wikipedia article, but `Michael_phelps` is valid in the API and redirects in the UI
* `Man_page` is the canonical Wikipedia article. `Man_Page` redirects in the UI but is invalid in API. If a user entered it in WikiViews, they would get a suggestion to try `Man_page` or `Man_Page`.

Another edge case is a title that contains small words such as "of" or "the". Like any other word, these are upper case when they are the first word of the title. But they are lower case when they fall in any other position. So I handled them in my title suggestion too, e.g.

```bash
❯ curl -X GET -H 'Accept: application/json' -H 'Content-Type: application/json' localhost:8080/pageviews\?article\=Call_Of_the_wild\&date=202401

{"error":"error: query for article param: Call_Of_the_wild did not return any results. Consider titlizing article param as Call_of_the_wild or Call_of_the_Wild."}
```

This validation is a best effort for V1, but a V2 would be to fall back on Wikipedia search, e.g.:

* Continue validate input article with simple checks, e.g, against forbidden characters
* Pass it to the Pageviews endpoint
* If the response is a 404, additionally pass the query to the Search endpoint
* Return the top hit, thereby giving the user the option to requery with the correct title

###### Invalid date

The Wikipedia endpoint is also brittle for date inputs. For monthly granularity, it turns out that the user must enter below exactly:

* for `start`, the first day of the month as YYYYMMDD
* for `end`, the last day of the same month as YYYYMMDD

Any variations from these rules will return a 404 in my testing. For example:

Valid month period -- Jan 1, 2024 - Jan 31, 2024

```bash
adambarcan in Code on main ❯ curl -s -X 'GET' \
  'https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia.org/all-access/all-agents/Man_page/monthly/20240101/20240131' \
  -H 'accept: application/json' | jq '.items[0].views'
11482
```

Invalid month period -- Jan 1, 2024 - Jan 30, 2024

```bash
adambarcan in Code on main ❯ curl -s -X 'GET' \
  'https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia.org/all-access/all-agents/Man_page/monthly/20240101/20240130' \
  -H 'accept: application/json' | jq '.items[0].views'
null
```

Invalid month period -- Jan 2, 2024 - Jan 31, 2024

```bash
adambarcan in Code on main ❯ curl -s -X 'GET' \
  'https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia.org/all-access/all-agents/Man_page/monthly/20240102/20240131' \
  -H 'accept: application/json' | jq '.items[0].views'
null
```

Invalid month period -- Jan 2, 2024 - Feb 2, 2024

```bash
adambarcan in Code on main ❯ curl -s -X 'GET' \
  'https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia.org/all-access/all-agents/Man_page/monthly/20240102/20240202' \
  -H 'accept: application/json' | jq '.items[0].views'
null
```

My solution is to only ask for a single date input, in the form

```bash
# e.g., 202402
YYYYMM
```

I then map this to the approriate `start` and `end` paramns when calling to the Wikipedia endpoint. There is validation for empty or mal-formed params, and I account for leap year.

## Troubleshooting

Requests to Wikipedia, user requests and errors are logged. Troubleshooting can be done by tailing docker logs, e.g.:

```bash
docker-compose logs -f
```

## Performance

A V2 optimization would be to cache requests.

## Security

Given all Wikipedia endpoints used are accesible without authentication, V1 of this project is also accessible without auth.

There is rate-limiting at 20 requests per second.

## Availability

V1 of this project runs as a single web server. If deployed to production, we would use a load-balancer and multiple replicas to ensure high availability. There is a `/healthcheck` endpoint that may be used for Kubernetes liveness and readiness probes.
