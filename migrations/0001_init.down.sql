-- Drop initial schema (reverse of 0001_init.up.sql)

DROP TABLE IF EXISTS settlements;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS vouchers;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;

-- Note: extension is left in place by default
