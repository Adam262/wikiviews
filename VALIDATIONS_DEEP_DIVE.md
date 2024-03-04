# Validations

## Article Param

### Invalid Article

There are several simple validation on the article param:

* empty article
* article with a forbidden character
* article beginning with a lower-case letter

See examples:

```bash
# Valid request
❯ curl -X GET localhost:8080/pageviews\?article\=Michael_Phelps\&date=202402
[{"article":"Michael_Phelps","timestamp":"2024020100","views":125860}]

# Invalid - starts with lower
❯ curl -X GET localhost:8080/pageviews\?article\=michael_Phelps\&date=202402
{"error":"error: article param michael_Phelps is invalid: param must not begin with a lower case character"}

# Invalid - empty param
❯ curl -X GET localhost:8080/pageviews             
{"error":"error: article param is invalid: param cannot be empty"}
```

### Suggestions for Titlized Article

Another common case is an article that returns a 404 from the Wikipedia endpoint, because their API cannot find any references to the article.

In these cases, I added validation that suggests a title case article. That is, if the user entered `MICHAEL_PHELPS`, the validation would suggest `Michael_Phelps` or `Michael_phelps`. I include both because both may be valid (and handled by redirection in the Wikipedia UI) and it is impossible to know which one the user intended.

For example:

```bash
❯ curl -X GET localhost:8080/pageviews\?article\=MICHAEL_Phelps\&date=202402
{"error":"error: query for article param: MICHAEL_Phelps did not return any results. Consider titlizing article param as Michael_phelps or Michael_Phelps."}
```

In the above case, `Michael_Phelps` is the canonical Wikipedia article, but `Michael_phelps` is valid in the API and redirects in the UI

Another example is a non-proper noun. For example, `Man_page` is the canonical Wikipedia article, while`Man_Page` redirects in the UI but is invalid in API. If a user entered it in WikiViews, they would get a suggestion to try `Man_page` or `Man_Page`. See examples:

```bash
# Valid name
❯ curl -X GET localhost:8080/pageviews\?article\=Man_page\&date=202402      
[{"article":"Man_page","timestamp":"2024020100","views":11211}]

# Invalid - make suggestion
❯ curl -X GET localhost:8080/pageviews\?article\=MAN_page\&date=202402
{"error":"error: query for article param: MAN_page did not return any results. Consider titlizing article param as Man_page or Man_Page."}
```

Another edge case is a title that contains small words such as "of" or "the". Like any other word, these are upper case when they are the first word of the title. But they are lower case when they fall in any other position. So I handled them in my title suggestion too, for example:

```bash
# Valid
curl -X GET -H 'Accept: application/json' -H 'Content-Type: application/json' localhost:8080/pageviews\?article\=Call_of_the_wild\&date=202401
[{"article":"Call_of_the_wild","timestamp":"2024010100","views":320}]

# Invalid - make suggestion
❯ curl -X GET -H 'Accept: application/json' -H 'Content-Type: application/json' localhost:8080/pageviews\?article\=Call_Of_the_wild\&date=202401
{"error":"error: query for article param: Call_Of_the_wild did not return any results. Consider titlizing article param as Call_of_the_wild or Call_of_the_Wild."}
```

Yet another edge case is a title that contains a permitted but escapable character — e.g. "?". I automatically HTML-escaped these characters to match the Wikipedia API approach.

For example:

```bash
❯ curl -X GET -H 'Accept: application/json' -H 'Content-Type: application/json' localhost:8080/pageviews\?article\=Are_You_the_One\?\&date=202201
[{"article":"Are_You_the_One?","timestamp":"2022010100","views":104545}]

# Request as sent to Wikipedia endpoint (from tailing WikiViews logs)
wikiviews  | 2024/03/04 03:30:37 sending GET request to Wikipedia endpoint: https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia.org/all-access/all-agents/Are_You_the_One%3F/monthly/20220101/20220131
```

### V2 Proposal - Fallback to Search

This validation is a best effort for V1, but a V2 would be to fall back on Wikipedia search, e.g.:

* Continue to validate the article param with simple checks, e.g, against forbidden characters
* Pass the article param to the Wikipedia `/pageviews` endpoint
* If the response is a 404, additionally pass the param to the Wikipedia `/search/title` endpoint
* Return the top hit, thereby giving the user the option to requery with the correct title

## Date Param

### Brittleness in Wikipedia endpoint

The underlying Wikipedia endpoint is also brittle for date inputs. For monthly granularity, it turns out that the user must enter below exactly:

* for `start`, the first day of the month as YYYYMMDD
* for `end`, the last day of the same month as YYYYMMDD

Any variations from these rules will return a 404 in my testing. For example:

```bash

# Valid month period — Jan 1, 2024 - Jan 31, 2024
❯ curl -s -X 'GET' \
  'https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia.org/all-access/all-agents/Man_page/monthly/20240101/20240131' \
  -H 'accept: application/json' | jq '.items[0].views'
11482

# Invalid month period — Jan 1, 2024 - Jan 30, 2024
❯ curl -s -X 'GET' \
  'https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia.org/all-access/all-agents/Man_page/monthly/20240101/20240130' \
  -H 'accept: application/json' | jq '.items[0].views'
null

# Invalid month period — Jan 2, 2024 - Jan 31, 2024
❯ curl -s -X 'GET' \
  'https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia.org/all-access/all-agents/Man_page/monthly/20240102/20240131' \
  -H 'accept: application/json' | jq '.items[0].views'
null

# Invalid month period — Jan 2, 2024 - Feb 2, 2024
❯ curl -s -X 'GET' \
  'https://wikimedia.org/api/rest_v1/metrics/pageviews/per-article/en.wikipedia.org/all-access/all-agents/Man_page/monthly/20240102/20240202' \
  -H 'accept: application/json' | jq '.items[0].views'
null
```

### WikiViews solution

The WikiViews `/pageviews` endpoint has a simple date param interface, in the form

```bash
# e.g., 202402
YYYYMM
```

This is mapped to the appropriate start and end params when calling to the Wikipedia endpoint. There is validation for empty or mal-formed params, and accounting for leap year.

For example, below are all valid queries:

```bash
❯ curl -X GET localhost:8080/pageviews\?article\=Michael_Phelps\&date=202402
[{"article":"Michael_Phelps","timestamp":"2024020100","views":125860}]

❯ curl -X GET -H 'Accept: application/json' -H 'Content-Type: application/json' localhost:8080/pageviews\?article\=Orca\&date=202201
[{"article":"Orca","timestamp":"2022010100","views":8420}]

❯ curl -X GET -H 'Accept: application/json' -H 'Content-Type: application/json' localhost:8080/pageviews\?article\=Breath_of_the_Wild\&date=202201
[{"article":"Breath_of_the_Wild","timestamp":"2022010100","views":980}]

 ❯ curl -X GET -H 'Accept: application/json' -H 'Content-Type: application/json' localhost:8080/pageviews\?article\=Are_You_the_One\?\&date=202201
[{"article":"Are_You_the_One?","timestamp":"2022010100","views":104545}]
```

And below are caught by validation:

```bash
# Missing date param
❯ curl -X GET -H 'Accept: application/json' -H 'Content-Type: application/json' localhost:8080/pageviews\?article\=Michael_Phelps             
{"error":"error: date param is invalid: param cannot be empty. Please enter in form YYYYMM"}

# Malformed date param
❯ curl -X GET -H 'Accept: application/json' -H 'Content-Type: application/json' localhost:8080/pageviews\?article\=Michael_Phelps\&date=2024  
{"error":"error: date param is invalid: please enter a valid year and month in form YYYYMM"}

# Date param with invalid year
❯ curl -X GET -H 'Accept: application/json' -H 'Content-Type: application/json' localhost:8080/pageviews\?article\=Michael_Phelps\&date=302401
{"error":"error: date param is invalid: please enter a valid year and month in form YYYYMM"}

# Date param with invalid month
❯ curl -X GET -H 'Accept: application/json' -H 'Content-Type: application/json' localhost:8080/pageviews\?article\=Michael_Phelps\&date=202313
{"error":"error: date param is invalid: please enter a valid year and month in form YYYYMM"}
```
