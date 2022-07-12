-- Prepare --
CREATE SCHEMA IF NOT EXISTS geolocation;

-- TODO Use normalized schema with separate tables for country, city, ip_address, etc.
-- It's impossible to do it with example data, that totally randomized.
CREATE TABLE geolocation.ip_location
(
    ip_address inet PRIMARY KEY,
    country_code varchar(2),
    country_name text,
    city text,
    latitude float,
    longitude float,
    mystery_value bigint
);

-- ---- create above / drop below ----

DROP TABLE geolocation.ip_location;

DROP SCHEMA geolocation;
