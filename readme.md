# Platts API Examples

## Table of Contents

**THIS IS NOT PRODUCTION READY CODE**

Each folder in CMD shows a different, but complimentary, use case for the Platts Market Data API

- [ Replacing a Datafeed with the Market Data API ](#replacing-a-datafeed-with-the-market-data-api)
- [ Augmenting with Corrections ](#incorporating-corrections)
- [ Getting Symbol Reference Data](#getting-reference-data)
- [ Listing the MDCs I have access to](#replacing-a-datafeed-with-api)

## Replacing a Datafeed with the Market Data API

This example shows how to use the History API to get all assessments since `t` (`modified_date`) and store the results in a local database. The idea is to invoke this program every `n` minutes in order to keep your database up to date.

| Parameter | Description                                                                                          |
| :-------- | :--------------------------------------------------------------------------------------------------- |
| apikey    | Your MarketData API Key. **Required**                                                                |
| username  | Your Platts Username. **Required**                                                                   |
| password  | Your Platts Password. **Required**                                                                   |
| t         | Modified Date. Get assessments since `t`. Defaults to `Now - 3Days`. Format is `2006-01-02T15:04:05` |
| p         | Page size. Defaults to 1000                                                                          |
| mdc       | Market Date Category. A grouping of symbols. _Optional_                                              |

### Getting Started

```go
go get github.com/mattn/go-sqlite3
go run cmd/assessments/assessments.go -t 2022-06-10T00:00:00 -apikey {APIKEY} -username {USERNAME} -password {PASSWORD} -mdc {MDC}
```

You will see logs in the console as API calls are made. Market data will be added to the `market_data` table in `database.db`.

```sql
sqlite3 database.db
SELECT * FROM market_data
LIMIT 20;
```

## Incorporating Corrections

This example compliments the above example by showing how to use the Corrections API to get Deletes/Backfills since `t` (`modified_date`) and update the local database accordingly. The idea is to invoke this program every `n` minutes in order to keep your database up to date.

| Parameter | Description                                                                                          |
| :-------- | :--------------------------------------------------------------------------------------------------- |
| apikey    | Your MarketData API Key. **Required**                                                                |
| username  | Your Platts Username. **Required**                                                                   |
| password  | Your Platts Password. **Required**                                                                   |
| t         | Modified Date. Get corrections since `t`. Defaults to `Now - 3Days`. Format is `2006-01-02T15:04:05` |
| p         | Page size. Defaults to 1,000                                                                         |

### Getting Started

```go
go get github.com/mattn/go-sqlite3
go run cmd/corrections/corrections.go -t 2022-06-10T00:00:00 -apikey {APIKEY} -username {USERNAME} -password {PASSWORD}
```

You will see logs in the console as API calls are made. Corrections will be reflected in the `market_data` table in `database.db`.

## Getting Reference Data

This example compliments the above by getting the reference data associated with a Symbol.

| Parameter | Description                           |
| :-------- | :------------------------------------ |
| apikey    | Your MarketData API Key. **Required** |
| username  | Your Platts Username. **Required**    |
| password  | Your Platts Password. **Required**    |
| p         | Page size. Defaults to 1,000          |

### Getting Started

```go
go get github.com/mattn/go-sqlite3
go run cmd/refdata/refdata.go -apikey {APIKEY} -username {USERNAME} -password {PASSWORD}
```

You will see logs in the console as API calls are made. Data will be added to `database.db`. Three tables are created:
| Table Name | Description |
| :-------- | :------------------------------------ |
| ref_data | One row per `symbol`
| sym_bate | One row per `symbol`, `bate`
| sym_mdc | One row per `symbol`, `mdc` |

Basic usage

```sql
sqlite3 database.db
SELECT * FROM ref_data
LIMIT 20;
```

Join with `market_data` on `symbol`

```sql
sqlite3 database.db
SELECT * FROM market_data md
INNER JOIN ref_data rd ON md.symbol = rd.symbol
LIMIT 10;
```
