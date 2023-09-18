CREATE TABLE orders (
  order_uid UUID PRIMARY KEY,
  track_number VARCHAR(64),
  entry VARCHAR(64),
  locale VARCHAR(8),
  internal_signature VARCHAR(128),
  customer_id VARCHAR(64),
  delivery_service VARCHAR(64),
  shardkey INTEGER,
  sm_id INTEGER,
  date_created TIMESTAMP,
  oof_shard INTEGER
);

CREATE TABLE delivery (
  order_uid UUID PRIMARY KEY REFERENCES orders(order_uid),
  name VARCHAR(128),
  phone VARCHAR(64),
  zip VARCHAR(64),
  city VARCHAR(128),
  address VARCHAR(256),
  region VARCHAR(128),
  email VARCHAR(128)
);

CREATE TABLE payment (
  order_uid UUID PRIMARY KEY REFERENCES orders(order_uid),
  transaction VARCHAR(128),
  request_id VARCHAR(128),
  currency VARCHAR(4),
  provider VARCHAR(64),
  amount INTEGER,
  payment_dt INTEGER,
  bank VARCHAR(64),
  delivery_cost INTEGER,
  goods_total INTEGER,
  custom_fee INTEGER
);

CREATE TABLE items (
  chrt_id INTEGER PRIMARY KEY,
  order_uid UUID REFERENCES orders(order_uid),
  track_number VARCHAR(64),
  price INTEGER,
  rid VARCHAR(128),
  name VARCHAR(256),
  sale INTEGER,
  size VARCHAR(64),
  total_price INTEGER,
  nm_id INTEGER,
  brand VARCHAR(256),
  status INTEGER
);
