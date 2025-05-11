CREATE TABLE listings (
  id SERIAL PRIMARY KEY,
  item_name TEXT,
  is_closed BOOLEAN DEFAULT FALSE
);

CREATE TABLE bids (
  id SERIAL PRIMARY KEY,
  bidder TEXT,
  amount INT,
  listing_id INT REFERENCES listings(id)
);
