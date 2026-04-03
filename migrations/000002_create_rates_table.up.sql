CREATE TYPE side_type AS ENUM ('ask', 'bid');

CREATE TABLE rates (
    id SERIAL PRIMARY KEY,
    price FLOAT NOT NULL,
    side side_type NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

SELECT attach_updated_at_trigger('rates');