CREATE TABLE users(
    id serial primary key,
    username varchar not null unique,
    password varchar not null,
    country_id int references countries
);

INSERT INTO users(username, password, country_id) VALUES ('TestAustralia', '$2a$04$YS1k0R.QaKJbF7U/UJdG/eY0tEm193vneUtVj1oOsg6ljUK5hiNS6', 1);
INSERT INTO users(username, password, country_id) VALUES ('TestAustria', '$2a$04$YS1k0R.QaKJbF7U/UJdG/eY0tEm193vneUtVj1oOsg6ljUK5hiNS6', 2);