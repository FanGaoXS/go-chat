CREATE TABLE "user"
(
    id         serial       NOT NULL primary key,
    nickname   varchar(256) NOT NULL,
    username   varchar(256) NOT NULL,
    password   varchar(256) NOT NULL,
    phone      varchar(256) NOT NULL,
    created_at timestamp    NULL DEFAULT now(),
    CONSTRAINT nickname_uq UNIQUE (nickname),
    CONSTRAINT username_uq UNIQUE (username)
);