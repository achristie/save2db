-- SQLite
select * from ref_data;


SELECT *
FROM ref_data 
WHERE 
    json_extract(json_each.value, '$.name') LIKE '%AGP%';