# Platts API Examples

## Getting Started

```
go get github.com/mattn/go-sqlite3
go run cmd/cli/main.go -apikey {APIKEY} -username {USERNAME} -password {PASSWORD}
```

Then check out the data in sqlite database

```
sqlite3 database.db
SELECT * FROM market_data
LIMIT 20;
```

## Explanation

**This is not production ready code. This is only meant to show a typical use case.**

This example is shows how to grab Market Data from the Platts API and store it in a database. This example queries by MDC and Modified Date in order to get all updates since a particular point in time. The idea here is to execute this function on a (short, 15min) interval with a sliding modified date in order to keep your local database up to date.

- Get an Access Token
- Get a list of MDC (Market Data Category) the User has access to
- For each MDC retrieve pricing data since time `t`
- Page through results (if necessary)
- Store results in database
- Get corrections (deletes) since time `t` and remove from database
- (not shown) `t` should slide so that you're updating `t` every invocation based on the time of your previous invocation
