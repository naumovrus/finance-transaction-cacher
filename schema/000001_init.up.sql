CREATE TABLE users
(
    id            serial       not null unique,
    first_name          varchar(255) not null,
    last_name       varchar(255) not null,
    username      varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE money_users
(
    id            serial       not null unique,
    amount  DECIMAL DEFAULT 0.0,
    user_id int references users (id) on delete cascade not null
);

CREATE TABLE transactionsSend
(
    id            serial       not null unique,
    user_id_from int references users (id) on delete cascade not null,
    user_id_to int references users (id) on delete cascade not null,
    date_time TIMESTAMP not null 
);

CREATE TABLE transactionsBalance
(
    id            serial       not null unique,
    user_id int references users (id) on delete cascade not null,
    date_time TIMESTAMP not null 
);