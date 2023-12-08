ALTER TABLE user_friend
    RENAME CONSTRAINT user_friend_uq TO friendship_uq;
ALTER TABLE user_friend
    RENAME CONSTRAINT user_friend_user_fk TO friendship_user_fk;
ALTER TABLE user_friend
    RENAME CONSTRAINT user_friend_friend_fk TO friendship_friend_fk;
ALTER TABLE user_friend
    RENAME TO friendship;


CREATE TABLE IF NOT EXISTS "friend_request"
(
    sender    varchar(256) NOT NULL,
    receiver  varchar(256) NOT NULL,
    status    varchar(256) NOT NULL,
    created_at timestamp   NULL DEFAULT now(),
    CONSTRAINT friend_request_sender_fk FOREIGN KEY (sender) REFERENCES "user" (subject),
    CONSTRAINT friend_request_receiver_fk FOREIGN KEY (receiver) REFERENCES "user" (subject)
);