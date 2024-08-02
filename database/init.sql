-- database/init.sql

-- Ensure the 'postgres' role exists
DO
$do$
BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'postgres') THEN
CREATE ROLE postgres;
END IF;
END
$do$;

-- Create the hft database if it does not exist
DO
$do$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'hft') THEN
        CREATE DATABASE hft;
END IF;
END
$do$;

-- Create a new user with a password if it does not exist
DO
$do$
BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'postgres') THEN
            CREATE USER postgres WITH ENCRYPTED PASSWORD 'postgres';
END IF;
END
$do$;

-- Grant all privileges on the hft database to the postgres user
GRANT ALL PRIVILEGES ON DATABASE hft TO postgres;

-- Connect to the hft database
\c hft

-- Create the ticker_data table
CREATE TABLE IF NOT EXISTS ticker_data (
           id SERIAL PRIMARY KEY,
           event_type TEXT,
           event_time BIGINT,
           symbol TEXT,
           price_change TEXT,
           price_change_percent TEXT,
           weighted_avg_price TEXT,
           first_trade_price TEXT,
           last_price TEXT,
           last_quantity TEXT,
           best_bid_price TEXT,
           best_bid_quantity TEXT,
           best_ask_price TEXT,
           best_ask_quantity TEXT,
           open_price TEXT,
           high_price TEXT,
           low_price TEXT,
           volume TEXT,
           quote_volume TEXT,
           statistics_open_time BIGINT,
           statistics_close_time BIGINT,
           first_trade_id BIGINT,
           last_trade_id BIGINT,
           total_trades BIGINT
);