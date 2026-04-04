CREATE TYPE app.side_type AS ENUM ('ask', 'bid');

CREATE TABLE app.rates (
    id SERIAL PRIMARY KEY,
    price FLOAT NOT NULL,
    side app.side_type NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

SELECT app.attach_updated_at_trigger('rates');