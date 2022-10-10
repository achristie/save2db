INSERT INTO
  assessments (
    symbol,
    bate,
    "value",
    assessed_date,
    modified_date,
    is_corrected
  )
VALUES
  ($1, $2, $3, $4, $5, $6)
  ON CONFLICT (symbol, bate, assessed_date) DO NOTHING;