-- No migrations yet because because I'm still spiking it out
-- and I don't want to figure out a migration system yet

BEGIN;

CREATE TABLE IF NOT EXISTS users (
  uri         TEXT PRIMARY KEY,
  username    TEXT NOT NULL,
  public_key  TEXT NOT NULL,
  private_key TEXT NOT NULL
);

COMMIT;
