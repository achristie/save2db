-- Get last 100 trades in EU BFOE
SELECT
  *
FROM
  trades,
  json_each(markets)
WHERE
  json_extract(json_each.value, '$.name') LIKE '%EU BFOE%'
ORDER BY
  order_date
LIMIT
  100;

-- Get trades in the last 1 month in the Markets Asia Japan Rack, Asia MD (PVO) for the Product Platts GO 10ppm
SELECT
  *
FROM
  trades,
  json_each(markets)
WHERE
  json_extract(json_each.value, '$.name') IN ("Asia Japan Rack", "ASIA MD (PVO)")
  AND product = "Platts GO 10ppm"
  AND order_date >= datetime('now', '-1 month');

-- Get the latest update time. A useful initial value for `t`
SELECT
  max(update_time)
FROM
  trades;

-- Get the Buyers and Sellers from the MOC on the last business day of last month
SELECT
  buyer,
  seller,
  count(*) AS trade_count
FROM
  trades
WHERE
  order_date = (
    SELECT
      max(order_date)
    FROM
      trades
    WHERE
      order_date <= datetime('now', 'start of month', '-1 day')
  )
GROUP BY
  buyer,
  seller
ORDER BY
  buyer,
  trade_count DESC;

-- Get the Max/Min Price each product was traded at per seller on the first business day of this month
SELECT
  seller,
  product,
  min(price),
  max(price)
FROM
  trades
WHERE
  order_date = (
    SELECT
      min(order_date)
    FROM
      trades
    WHERE
      order_date >= datetime('now', 'start of month')
  )
GROUP BY
  seller,
  product
ORDER BY
  product;

-- Get consumated Trades with BP Products North America is the Market Maker or Counterparty
SELECT
  *
FROM
  trades
WHERE
  order_state = "consummated"
  AND (
    market_maker = "BP Products North America Inc."
    OR counterparty = "BP Products North America Inc."
  )
ORDER BY
  order_date DESC;