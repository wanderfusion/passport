create table users
(
    id  uuid    default gen_random_uuid()   not null    primary key,
    created_at  timestamp   default CURRENT_TIMESTAMP,
    updated_at  timestamp   default CURRENT_TIMESTAMP,
    email   varchar(255)    not null    unique,
    hashed_password varchar(255)    not null,
    username    varchar(100)    unique,
    profile_picture text
);

---- create above / drop below ----

drop table users;
