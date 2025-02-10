BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS system_config (
    key VARCHAR(50) PRIMARY KEY,
    value VARCHAR(255)
);

INSERT INTO system_config (key, value)
VALUES ('csv_imported', 'false')
ON CONFLICT (key) DO NOTHING;

CREATE TABLE IF NOT EXISTS assets (
  id      UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
  symbol  VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS portfolios (
  id           UUID PRIMARY KEY,
  initial_cash NUMERIC NOT NULL,
  cash         NUMERIC NOT NULL
);

CREATE TABLE IF NOT EXISTS positions (
  portfolio_id UUID NOT NULL,
  asset_id     UUID NOT NULL,
  quantity     INTEGER NOT NULL,
  PRIMARY KEY (portfolio_id, asset_id),
  FOREIGN KEY (portfolio_id) REFERENCES portfolios(id) ON DELETE CASCADE,
  FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS prices (
  asset_id  UUID NOT NULL,
  value     NUMERIC NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (asset_id, timestamp),
  FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE
);

CREATE TYPE sides as enum ('buy','sell');

CREATE TABLE IF NOT EXISTS trades (
  id        UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
  asset_id  UUID NOT NULL,
  side      sides NOT NULL,
  price     NUMERIC NOT NULL,
  quantity  INTEGER NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE
);

COMMIT;