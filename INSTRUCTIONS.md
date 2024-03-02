## Grow Therapy Platform Engineering Take-Home

This programming project is meant to gauge back-end web app development skills, ability to
use third-party APIs, and ability to use containers.

### Project Requirements

This exercise involves creating a containerized back-end web application:

- Back-end web application that acts as a wrapper around the Wikipedia API. The
application should expose an endpoint that returns the view count for a given article for a
given month.
- The application does not need a front-end.
- We recommend using a mainstream language and web framework, e.g. TypeScript/Node/Express, Python/Flask, Java/Spring.
- The application should use communicate with the Wikipedia API directly and not use a third-party package, e.g. python-mviews
- Application is containerized using Docker
- Feel free to use additional wrappers around `docker`, e.g. a Makefile

## Tips for a Great Take Home

- Add a README file
- Document your API endpoint
- Write tests for your application
- Make it easy for our team to build and run your application locally
- Utilize best practices for containerization and Docker
- Think through how you would want to handle corner cases and invalid requests
- Think through performance, security, and reliability implications including what you would want to do in Production

## Next Steps

Reach out to Rachel (<rachel@growtherapy.com>) with any questions.

After you’re done,

1. Upload your project to a private Github repository
1. Give Rachel (<rachel@growtherapy.com>) access
1. Email Rachel a link

Note that we will not be looking at commit history or commit messages.

## To do

- Test coverage - ~
  - Should solve Michael_Phelps v Michael_phelps. Return both?
    - yes, throw out null, otherwise return both as pretty item
  - Orca - Done
  - Man_page - Done
- Docs
- Makefile? - DONE
- CLI with flags
  - this would make it easier to validate date
  - but start with just year query param - validate it is YYYY
  - month query param - accept Feb, February,
- Error handling - this is pretty good but should wrap a 404 if title is say: ddffdfd
  - v2 would add search for that
- Performance
  - No idea
  - Cache in redis?
    - Michael_Phelps202402
      - Handle leap year
- HA
  - replica count ?
  - add note saying v2 is on k8s with load balancer
  - Hard in Docker to do this. How about a health check
- Security
  - Auth? Not needed, there is no auth to API
  - handle api keys  - NONE
