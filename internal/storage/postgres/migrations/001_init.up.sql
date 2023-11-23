CREATE TABLE IF NOT EXISTS "user"
(
    subject    varchar(256) NOT NULL primary key,
    nickname   varchar(256) NOT NULL,
    username   varchar(256) NOT NULL,
    password   varchar(256) NOT NULL,
    phone      varchar(256) NOT NULL,
    created_at timestamp    NULL DEFAULT now(),
    CONSTRAINT nickname_uq UNIQUE (nickname),
    CONSTRAINT username_uq UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS "group"
(
    id         serial       NOT NULL primary key,
    name     varchar(256) NOT NULL,
    "type"     varchar(256) NOT NULL,
    is_public  bool         NOT NULL DEFAULT false,
    created_by varchar(256) NOT NULL,
    created_at timestamp    NULL DEFAULT now(),
    CONSTRAINT group_uq UNIQUE (name, created_by),
    CONSTRAINT group_user_fk FOREIGN KEY (created_by) REFERENCES "user" (subject)
);