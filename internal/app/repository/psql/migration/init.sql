CREATE DATABASE IF NOT EXISTS urls_shorten;
CREATE TABLE IF NOT EXISTS shorten_urls
(
    short_url text      PRIMARY KEY,
    origin_url text     NOT NULL,
    visits BIGINT       DEFAULT 0,
    date_created        timestamp WITH TIME ZONE DEFAULT NOW(),
    date_updated        timestamp WITH TIME ZONE DEFAULT NOW()
);