CREATE TABLE user_questions(
    id serial primary key,
    chat_id bigint not null,
    country_id int not null,
    question text not null,
    status smallint not null
);