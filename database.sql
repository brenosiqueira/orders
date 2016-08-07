CREATE KEYSPACE "orders" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

USE "orders";

CREATE TYPE order_item(sku uuid, unit_price INT, quantity INT);
CREATE TYPE transaction(id uuid, external_id varchar, amount INT, type varchar, authorization_code varchar,
  card_brand varchar, card_bin varchar, card_last varchar);

CREATE TABLE orders(id uuid primary key, number varchar, reference varchar, status varchar, created_at TIMESTAMP, updated_at TIMESTAMP, notes text, price INT,
  items list<frozen<order_item>>, transactions list<frozen<transaction>>);

