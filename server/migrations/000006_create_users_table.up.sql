CREATE TABLE users(
    id serial primary key,
    username varchar not null unique,
    password varchar not null,
    country_id int references countries
)