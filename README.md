# WikiViews

## Overview

WikiViews is a simple Golang server with a single JSON endpoint, `/pageviews`. Its responsibility is to respond to user queries for monthly pageview data for English-language Wikipedia articles. Although users can alternatively query the Wikipedia API directly, WikiViews provides several enhancements such as a simplied interface, response caching, param validation and sensible formatting of params that would otherwise be invalid.

## Dependencies

WikiViews is Dockerized. Its only dependency for running locally as a Docker container is [Docker Desktop](https://www.docker.com/products/docker-desktop/).

## Getting Started

To run the project locally, git clone this repo and run the below command from within the repo:

```bash
docker-compose up --build -d
docker-compose logs -f
```

To stop the project

```bash
docker-compose down
```

WikiViews may also be run locally without Docker for faster iteration. Local development takes advantage of the [air Go package](https://github.com/cosmtrek/air) for live reload.

```bash
go install github.com/cosmtrek/air@latest
make run
```

## API

### /healthcheck

This method is a simple health check. It may be used for Kubernetes liveness and readiness probes.

```
curl -X GET http://localhost:8080/ping

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

##### monthstart (int8)

The date of the first day of the target month, expressed as:

```
# e.g. 20231101 is November 1, 2023
YYYYMMDD
```

##### monthend (int8)

The date of the last day of the target month, expressed as:

```bash
# e.g. 20231130 is November 30, 2023
YYYYMMDD
```

#### Sample Request and Response

```bash
curl -X GET localhost:8080/pageviews\?article\=michael_phelps\&monthstart=20240201\&monthend=20240229
```

The response is a JSON-ified list of response objects, containing data as the article name, time period and pageview count.

```bash
[{"project":"en.wikipedia","article":"Michael_Phelps","granularity":"monthly","timestamp":"2024020100","views":125860}]
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

##### The `Michael_Phelps` vs. `Man_pages` problem

At a minimum, I knew I needed to implement validation on all passed-in params. Any of these requests will return a sensible error message:

###### Invalid and uncorrectable article*

The below article does not exist, even if I try to correct it to `Michaelphelps`. The Wikipedia endpoint returns a 404, that I wrap and pass on with a sensible error message

```bash
curl -X GET localhost:8080/pageviews\?article\=michaelphelps\&monthstart=20240201\&monthend=20240229
```

###### Invalid date

Below is entered an invalid date. The Wikipedia endpoint returns a 404, that I wrap and pass on with a sensible error message

###### Correctable article

In iterating on this project, I noticed that the Wikipedia API is quite brittle. The UI provides Search, which is forgiving of variations such as `Michael_Phelps`, `michael_phelps`, `MICHAEL_PHELPS`; they all return a top hit of the correct article with key `Michael_Phelps`. But the Pageviews API is way more strict - most variations will return a 404

My first solution was to simply correct the input. But that ran into what I call the `Michael_Phelps vs Man_pages` problem. That is, most valid article titles fall into one of the below three forms:

* *Single Word* This is easy. Just conver the input to title case, e.g.: *Dog* or *Orca*
* *Many Words, Proper Noun* These titles should have all words in title case, e.g.: *Michael_Phelps* or *New_York_City*
* *Many Words, Non-Proper Noun* These titles should have only the first word in title case, e.g.: *Man_page* or *Killer_whale*

So the issue this raises is that it is impossible to know if the intended form of an article should be `One_Two` or `One_two`. I considered returning both forms. But both forms do not always exist, although in some cases they do, because of redirects - e.g. `Michael_Phelps` and `Michael_phelps` are both valid queries to the Pageviews endpoint, with different results.

So my solution is to fall back on Wikipedia search

* Validate input article with simple checks, e.g, against forbidden characters
* Pass it to the Pageviews endpoint
* If the response is a 404, additionally pass the query to the Search endpoint
* Return the top hit, thereby giving the user the option to requery with the correct title

## Troubleshooting

Requests are logged. Troubleshooting can be done by tailing docker logs, e.g.:

```bash
docker-compose logs -f
```

## Performance

## Security

All Wikipedia endpoints used are accesible without authentication. For simplicity, V1 of this project is also accessible without auth. A V2 could implement auth.

There is rate-limiting at 20 requests per second.

## Availability

V1 of this project runs as a single web server. If deployed to production, we would use a load-balancer and multiple replicas to ensure high availability. There is a /healthcheck endpoint that may be used for Kubernetes liveness and readiness probes.
