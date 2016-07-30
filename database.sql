CREATE KEYSPACE "orders" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

USE "orders";

CREATE TABLE "order"(id uuid primary key, number varchar, reference varchar, status varchar, created_at TIMESTAMP,
  updated_at TIMESTAMP, notes text, price INT);

CREATE TABLE "order_item"( sku uuid primary key, order_id uuid, unit_price INT, quantity INT);

CREATE TABLE "transaction"(id uuid primary key, order_id uuid, external_id varchar, amount INT, type varchar, authorization_code varchar,
  card_brand varchar, card_bin varchar, card_last varchar);