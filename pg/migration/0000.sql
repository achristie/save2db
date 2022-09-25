CREATE TABLE
  assessments (
    symbol CHAR(7) NOT NULL,
    bate CHAR(1) NOT NULL,
    "value" DECIMAL NOT NULL,
    assessed_date TIMESTAMPTZ NOT NULL,
    modified_date TIMESTAMPTZ NOT NULL,
    is_corrected CHAR(1) NOT NULL,
    PRIMARY KEY (symbol, bate, assessed_date)
  );

CREATE INDEX assessments_ad ON assessments (assessed_date);

CREATE INDEX assessments_bate ON assessments (bate);