CREATE TABLE IF NOT EXISTS matches (
    match_id TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS ratings (
    user_id BIGINT PRIMARY KEY,
    rating DOUBLE PRECISION,
    rating_deviation DOUBLE PRECISION,
    rating_volatility DOUBLE PRECISION,
    last_updated TIMESTAMPTZ
);