CREATE DATABASE url_shorten;
CREATE TABLE IF NOT EXISTS shorten_urls (
                            short_url text primary key unique,
                            origin_url text not null,
                            visits int default 0,
                            date_created timestamp,
                            date_updated timestamp
);