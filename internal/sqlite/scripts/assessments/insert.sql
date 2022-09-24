INSERT
OR REPLACE INTO assessments (
  symbol,
  bate,
  VALUE,
  assessed_date,
  modified_date,
  is_corrected
)
VALUES
  (?, ?, ?, ?, ?, ?)