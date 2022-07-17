-- Prepare --
CREATE SCHEMA IF NOT EXISTS geolocation;

-- TODO Use normalized schema with separate tables for country, city, ip_address, etc.
-- It's impossible to do it with example data, that totally randomized.
CREATE TABLE geolocation.ip_location
(
    -- As we have many records with the same ip_address and don't have reliable deduplication logic, use id as primary key.
    id serial PRIMARY KEY,
    ip_address inet,
    country_code varchar(2),
    country_name text,
    city text,
    latitude float,
    longitude float,
    mystery_value bigint
);

CREATE INDEX ip_location_ip_idx ON geolocation.ip_location USING hash(ip_address);

-- ---- create above / drop below ----

DROP INDEX ip_location_ip_idx;

DROP TABLE geolocation.ip_location;

DROP SCHEMA geolocation;
