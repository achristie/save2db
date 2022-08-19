SELECT
  *
FROM
  trades,
  json_each(markets)
WHERE
  json_extract(json_each.value, '$.name') LIKE '%EU BFOE%';

select
  markets
from
  trades
limit
  5;