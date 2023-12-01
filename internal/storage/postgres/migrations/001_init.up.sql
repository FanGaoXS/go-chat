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

CREATE TABLE IF NOT EXISTS "user_friend"
(
    user_subject   varchar(256) NOT NULL,
    friend_subject varchar(256) NOT NULL,
    added_at       timestamp NULL DEFAULT now(),
    CONSTRAINT user_friend_uq UNIQUE (user_subject, friend_subject),
    CONSTRAINT user_friend_user_fk FOREIGN KEY (user_subject) REFERENCES "user" (subject),
    CONSTRAINT user_friend_friend_fk FOREIGN KEY (friend_subject) REFERENCES "user" (subject)
);

CREATE TABLE IF NOT EXISTS "group"
(
    id         serial       NOT NULL primary key,
    name       varchar(256) NOT NULL,
    "type"     varchar(256) NOT NULL,
    is_public  bool         NOT NULL DEFAULT false,
    created_by varchar(256) NOT NULL,
    created_at timestamp    NULL DEFAULT now(),
    CONSTRAINT group_uq UNIQUE (name, created_by),
    CONSTRAINT group_user_fk FOREIGN KEY (created_by) REFERENCES "user" (subject)
);

CREATE TABLE IF NOT EXISTS "group_member"
(
    user_subject varchar(256) NOT NULL,
    group_id     bigint       NOT NULL,
    join_at      timestamp NULL DEFAULT now(),
    CONSTRAINT group_member_uq UNIQUE (user_subject, group_id),
    CONSTRAINT group_member_user_fk FOREIGN KEY (user_subject) REFERENCES "user" (subject),
    CONSTRAINT group_member_group_fk FOREIGN KEY (group_id) REFERENCES "group" (id)
);

CREATE TABLE IF NOT EXISTS "broadcast_record"
(
    id         serial       NOT NULL primary key,
    content    varchar(256) NOT NULL,
    created_by varchar(256) NOT NULL,
    created_at timestamp    NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "group_record"
(
    id         serial       NOT NULL primary key,
    content    varchar(256) NOT NULL,
    group_id   bigint       NOT NULL,
    created_by varchar(256) NOT NULL,
    created_at timestamp    NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "private_record"
(
    id         varchar(256) NOT NULL primary key,
    hash_id    varchar(256) NOT NULL,
    content    varchar(256) NOT NULL,
    sender     varchar(256) NOT NULL,
    receiver   varchar(256) NOT NULL,
    created_at timestamp    NULL DEFAULT now()
);