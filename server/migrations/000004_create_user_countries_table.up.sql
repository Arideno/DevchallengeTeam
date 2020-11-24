CREATE TABLE user_countries(
    chatId bigint not null unique,
    countryId int not null,
    topicId int
);