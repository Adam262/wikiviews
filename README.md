## WikiView

## Dependencies
* Docker Desktop
* asdf?

## Getting Started
* docker compose up -d or docker run --name wikiviews -p 8080:8080 -d wikiviews 
* Note why I had to map to port 8080 (to listen on all network interfaces rather than locahost)

## Troubleshooting
* docker logs


## To do
* Test coverage
* Docs
* Makefile?
* Error handling
* Performance
    * caching
* HA
    * replica count
* Security
    * Auth?
    * handle api keys 
