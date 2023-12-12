ALTER TABLE user_friend
    RENAME CONSTRAINT user_friend_uq TO friendship_uq;
ALTER TABLE user_friend
    RENAME CONSTRAINT user_friend_user_fk TO friendship_user_fk;
ALTER TABLE user_friend
    RENAME CONSTRAINT user_friend_friend_fk TO friendship_friend_fk;
ALTER TABLE user_friend
    RENAME TO friendship;


CREATE TABLE IF NOT EXISTS "friend_request_log"
(
    id        serial       NOT NULL PRIMARY KEY,
    sender    varchar(256) NOT NULL,
    receiver  varchar(256) NOT NULL,
    status    varchar(256) NOT NULL,
    created_at timestamp   NULL DEFAULT now(),
    CONSTRAINT friend_request_log_sender_fk FOREIGN KEY (sender) REFERENCES "user" (subject),
    CONSTRAINT friend_request_log_receiver_fk FOREIGN KEY (receiver) REFERENCES "user" (subject)
);

CREATE TABLE IF NOT EXISTS "group_request_log"
(
    id serial NOT NULL PRIMARY KEY,
    group_id bigint NOT NULL,
    sender varchar(256) NOT NULL,
    status varchar(256) NOT NULL,
    created_at timestamp   NULL DEFAULT now(),
    CONSTRAINT group_request_log_group_fk FOREIGN KEY (group_id) REFERENCES "group" (id),
    CONSTRAINT group_request_log_sender_fk FOREIGN KEY (sender) REFERENCES "user" (subject)
);

CREATE TABLE IF NOT EXISTS "group_invitation_log"
(
    id serial NOT NULL PRIMARY KEY,
    group_id bigint NOT NULL,
    sender varchar(256) NOT NULL,
    receiver varchar(256) NOT NULL,
    status varchar(256) NOT NULL,
    created_at timestamp   NULL DEFAULT now(),
    CONSTRAINT group_invitation_log_group_fk FOREIGN KEY (group_id) REFERENCES "group" (id),
    CONSTRAINT group_invitation_log_sender_fk FOREIGN KEY (sender) REFERENCES "user" (subject),
    CONSTRAINT group_invitation_log_receiver_fk FOREIGN KEY (receiver) REFERENCES "user" (subject)
);