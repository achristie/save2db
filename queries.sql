-- SQLite
select * from ref_data;


SELECT *
FROM ref_data, json_each(mdc)
WHERE 
    json_extract(json_each.value, '$.description') LIKE '%Refined%';


select *
from assessments
limit 100;


select *, count(*) from assessments
group by symbol, bate, assessed_date
having count(*) > 1;


select * from assessments
where symbol = "AABRZ05" and bate = "c" and assessed_date = datetime('2022-06-20T 00:00:00');