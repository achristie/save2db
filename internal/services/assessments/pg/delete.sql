DELETE FROM
  assessments
WHERE
  symbol = $1
  AND bate = $2
  AND assessed_date = $3