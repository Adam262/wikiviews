# WikiViews

## Overview

WikiViews is a simple Golang server with a single JSON endpoint, `/pageviews`. Its responsibility is to respond to user queries for monthly pageview data for English-language Wikipedia articles. Although users can alternatively query the Wikipedia API directly, WikiViews provides several enhancements such as a simplied interface and param validation.

See [project instructions](./INSTRUCTIONS.md)

## Dependencies

WikiViews is Dockerized. Its only dependency for running locally as a Docker container is [Docker Desktop](https://www.docker.com/products/docker-desktop/). To develop on your machine without Docker, [Golang](https://go.dev/doc/install) is required.

## Getting Started

### Via Docker Compose

To run the project locally, git clone this repo and run the below command from within the repo:

```bash
# Start the API server (builds a Docker image if needed). It is exposed to port :8080
❯ docker-compose up -d 

# Run a health check. Note it may take a few seconds for the server to be healthy and the port exposed
# During this time, you may see the response: url: (56) Recv failure: Connection reset by peer
❯ curl -X GET http://localhost:8080/healthcheck

# Tail logs
❯ docker-compose logs -f

# Gracefully kill the server
❯ docker-compose down

# Build a Docker image only without running it 
❯ make docker-build
```

### Run locally

WikiViews may also be run locally without Docker for faster iteration. Local development takes advantage of the [air Go package](https://github.com/cosmtrek/air) for live reload.

#### As live reload server

```bash
❯ go install github.com/cosmtrek/air@latest
❯ make run
```

WikiViews may also be run locally by compiling and running a binary

#### As binary

```bash
❯ make build
❯ ./cmd/http-server/http-server
```

### Running tests

Run unit tests via:

```bash
❯ make test
```

### Vetting and package management

Other Make targets are exposed for local development

```bash
# Run Golang compile checks
❯ make vet

# This command performs several housekeeping functions such as `go mod tidy` and `go mod vendor`
# It is important to run this command whenever you change a module invocation - e.g, when you add or remove an import
❯ make mod
```

## API

### /healthcheck

This method is a simple health check. It may be used for Kubernetes liveness and readiness probes.

```bash
❯ curl -X GET http://localhost:8080/healthcheck
ok
```

### /pageviews

This endpoint accepts JSON queries to the [Wikipedia Pageviews REST API](https://wikimedia.org/api/rest_v1/#/Pageviews%20data). It returns a JSON-ified list of response objects, containing data as the article name, time period and pageview count.

#### Params

##### article (string)

The title of the Wikipedia article. It must follow the [naming conventions](https://en.wikipedia.org/wiki/Wikipedia:Naming_conventions_(technical_restrictions)) defined by Wikipedia, which can be summarized as:

* The title must begin with a capital letter. Any additional words in the title may be either capital or lower-case
* A space between words must be entered as a single underscore
* These characters are forbidden anywhere in the article title: `# < > [ ] { } |`

##### date (int8)

The target year and month to query, expressed as:

```bash
# e.g. 202311 is November, 2023
YYYYMM
```

This format is an optimization of the underlying Wikipedia endpoint, which expects separate params for the start and end of the monthly period. That is, to express `November 2023`, the Wikipedia API caller would need two path parameters in their query:

* start — `20241101`
* end — `20241130`

See [Validations — Date Param](./VALIDATIONS_DEEP_DIVE.md#date-param) for more discussion and examples.

#### Sample Request and Response

```bash
❯ curl -X GET localhost:8080/pageviews\?article\=MichaeL_Phelps\&date=202402

# The response is a JSON-ified list of response objects, containing data as the article name, time period and pageview count.
[{"article":"Michael_Phelps","timestamp":"2024020100","views":125860}]
```

#### Endpoint Design Decisions

I made several design decisions in V1 of the endpoint for the sake of simplifying the interface, validating params and trying to mitigate the brittleness of the underlying endpoint

##### English-language only

I made the decision to only query on English-language articles. This decision had two benefits:

* Removed an additional param — *project* — that users would otherwise need to pass in
* Simplified the regex for validating article titles by removing the need to deal with non-English characters

##### Hard-code other params

I hard-coded three other params (from the Wikipedia endpoint) to simplify the user interface:

* *access*. This param filters by page access method, e.g.: *desktop*, *mobile-app* or *mobile-web*. I hard-coded to *all-access*.
* *agent*. This param filters by page agent, e.g.: *user*, *automated* or *spider*. I hard-coded to *all-agents*.
* *granularity*. This param sets the time unit for the response data, e.g.: *daily* or *monthly*. I hard-coded to *monthly*.

##### Validations

I employed several validations and formatters for both article and date. Please see [Validations Deep Dive](./VALIDATIONS_DEEP_DIVE.md)

## Troubleshooting

Requests to Wikipedia, user requests and errors are logged. Troubleshooting can be done by tailing docker logs, e.g.:

```bash
❯ docker-compose logs -f
```

## Performance

A V2 optimization would be to cache requests and pageviews in a Redis database. A schema could be:

* Key - `Article_YYYY_MM` string
* Value - `Views` int

```bash
{
  "Michael_Phelps_2024_02": 125860,
  "Michael_Phelps_2024_01": 168202,
  ...and_so on
}
```

Upon receiving an incoming query, WikiViews would first check the cache. If there is a cache hit, WikiViews would respond to the request without needing a Wikipedia API call. If there is a miss, WikiViews would query the Wikipedia API, and write to the cache before responding to the user.

## Security

Given all Wikipedia endpoints used are accessible without authentication, V1 of this project is also accessible without auth.

All user article param input is html-escaped.

There is rate-limiting at 20 requests per second.

## Availability

V1 of this project runs as a single web server. If deployed to production, we would use a load-balancer and multiple replicas to ensure high availability. There is a `/healthcheck` endpoint that may be used for Kubernetes liveness and readiness probes.
