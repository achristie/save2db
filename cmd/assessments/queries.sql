-- SQLite
select
  *
from
  ref_data;

SELECT
  *
FROM
  ref_data,
  json_each(mdc)
WHERE
  json_extract(json_each.value, '$.description') LIKE '%Oil%';

select
  *
from
  assessments;

select
  *,
  count(*)
from
  assessments
group by
  symbol,
  bate,
  assessed_date
having
  count(*) > 1;

select
  count(*)
from
  assessments;

select
  *
from
  assessments
where
  assessed_date between '2022-05-01' and '2022-06-01';

-- find the largest close price per symbol per month
select
  symbol,
  max(price),
  strftime("%Y-%m", assessed_date) as year_month,
  assessed_date
from
  assessments
where
  bate = 'c'
group by
  symbol,
  year_month;

-- all the latest prices for symbols in 'Refined'
-- find the latest modified date. A reasonable choice for `t` when running assessments.exe
select
  max(modified_date)
from
  assessments;