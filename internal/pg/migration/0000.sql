CREATE TABLE
  assessments (
    symbol CHAR(7) NOT NULL,
    bate CHAR(1) NOT NULL,
    "value" DECIMAL NOT NULL,
    assessed_date TIMESTAMP NOT NULL,
    modified_date TIMESTAMP NOT NULL,
    is_corrected CHAR(1) NOT NULL,
    PRIMARY KEY (symbol, bate, assessed_date)
  );

CREATE INDEX plt_assessments_ad ON assessments (assessed_date);

CREATE INDEX plt_assessments_bate ON assessments (bate);

CREATE TABLE
  symbols (
    symbol CHAR(7) NOT NULL PRIMARY KEY,
    assessment_frequency VARCHAR(50),
    commodity VARCHAR(50),
    contract_type VARCHAR(50),
    description VARCHAR(100),
    publication_frequency_code VARCHAR(10),
    currency VARCHAR(30),
    quotation_style VARCHAR(30),
    delivery_region VARCHAR(80),
    delivery_region_basis VARCHAR(80),
    settlement_type VARCHAR(50),
    active VARCHAR(50),
    TIMESTAMP VARCHAR(50),
    uom VARCHAR(50),
    day_of_publication VARCHAR(50),
    shipping_terms VARCHAR(80),
    standard_lot_size INT,
    commodity_grade VARCHAR(80),
    standard_lot_units VARCHAR(50),
    decimal_places INT,
    mdc JSON,
    bates JSON
  );

CREATE INDEX plt_symbols_cmdty ON symbols (commodity);

CREATE INDEX plt_symbols_drb ON symbols (delivery_region_basis);

CREATE TABLE
  trades (
    markets json,
    product TEXT,
    strip TEXT,
    hub TEXT,
    update_time TEXT,
    market_maker TEXT,
    order_type TEXT,
    order_state TEXT,
    buyer TEXT,
    seller TEXT,
    price REAL,
    price_unit TEXT,
    order_quantity REAL,
    lot_size INTEGER,
    lot_unit TEXT,
    order_begin TEXT,
    order_end TEXT,
    order_date TEXT,
    order_time TEXT,
    order_id INTEGER,
    order_sequence INTEGER,
    deal_id INTEGER,
    deal_begin TEXT,
    deal_end TEXT,
    deal_quantity REAL,
    deal_quantity_max REAL,
    deal_quantity_min REAL,
    deal_terms TEXT,
    counterparty_parent TEXT,
    counterparty TEXT,
    market_maker_parent text,
    buyer_parent TEXT,
    seller_parent TEXT,
    buyer_mnemonic TEXT,
    seller_mnemonic TEXT,
    market_maker_mnemonic TEXT,
    counterparty_mnemonic TEXT,
    window_region TEXT,
    market_type TEXT,
    c1_price_basis TEXT,
    c1_percentage REAL,
    c1_price REAL,
    c1_basis_period TEXT,
    c1_basis_period_details TEXT,
    c2_price_basis TEXT,
    c2_percentage REAL,
    c2_price REAL,
    c2_basis_period TEXT,
    c2_basis_period_details TEXT,
    c3_price_basis TEXT,
    c3_percentage REAL,
    c3_price REAL,
    c3_basis_period TEXT,
    c3_basis_period_details TEXT,
    window_state TEXT,
    order_classification TEXT,
    oco_order_id TEXT,
    reference_order_id INTEGER,
    order_platts_id INTEGER,
    order_cancelled TEXT,
    order_derived TEXT,
    order_quantity_total REAL,
    order_repeat TEXT,
    leg_prices TEXT,
    parent_deal_id TEXT,
    order_spread TEXT,
    order_state_detail TEXT,
    PRIMARY KEY (
      order_id,
      order_sequence,
      order_platts_id,
      order_time
    )
  );

CREATE INDEX plt_trades_product_idx ON trades (product);