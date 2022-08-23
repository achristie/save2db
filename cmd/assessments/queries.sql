-- Get the symbols that have a MDC Description that contains "Electricity"
SELECT
  *
FROM
  symbols,
  json_each(mdc)
WHERE
  json_extract(json_each.value, '$.description') LIKE '%Electricity%';

-- Get the list of MDCs and their descriptions
SELECT
  json_extract(json_each.value, "$.name") name,
  json_extract(json_each.value, "$.description") description,
  count(*) symbol_count
FROM
  symbols,
  json_each(mdc)
GROUP BY
  name;

-- Get the list of Commodoties and Regions
SELECT
  commodity,
  delivery_region,
  count(*) as symbol_count
FROM
  symbols
GROUP BY
  commodity,
  delivery_region
ORDER BY
  commodity;

-- Get the largest close price per month in $/BBL
SELECT
  s.symbol,
  max(a.value) max_value,
  strftime("%Y-%m", a.assessed_date) as year_month,
  a.assessed_date
FROM
  symbols s
  JOIN assessments a on s.symbol = a.symbol
WHERE
  a.bate = "c"
  AND s.uom = "BBL"
  AND s.currency = "USD"
GROUP BY
  year_month
ORDER BY
  year_month ASC;

-- Get last 12 months of spot assessments for all symbols tagged with "Crude oil" commodity
SELECT
  s.symbol,
  s.description,
  s.delivery_region,
  s.delivery_region_basis,
  s.currency,
  s.uom,
  a.bate,
  a.assessed_date
FROM
  symbols s
  JOIN assessments a ON s.symbol = a.symbol
WHERE
  s.commodity = "Crude oil"
  AND a.assessed_date >= datetime('now', "-1 year")
  AND s.contract_type = "Spot";

-- find the latest modified date. A reasonable choice for `t` when running assessments.exe to fetch updated date
SELECT
  max(modified_date)
FROM
  assessments;

-- Get all assessments for a list of symbols and include description and currency
SELECT
  s.symbol,
  s.description,
  s.currency,
  a.bate,
  a.value,
  a.assessed_date
FROM
  symbols s
  JOIN assessments a ON s.symbol = a.symbol
WHERE
  s.symbol IN ("AAGOQ00", "AAGOP00")
  AND a.bate = "c"
ORDER BY
  a.assessed_date ASC;