CREATE TABLE user_messages(
    id serial primary key,
    chat_id bigint not null,
    message text not null,
    from_operator boolean not null,
    question_id int not null
)