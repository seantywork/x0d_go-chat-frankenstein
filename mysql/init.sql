CREATE USER 'seantywork'@'%' IDENTIFIED BY 'youdonthavetoknow';

GRANT ALL PRIVILEGES ON *.* TO 'seantywork'@'%';

CREATE DATABASE chfrank;

USE chfrank;


CREATE TABLE chfrank_user (uuid VARCHAR(150), usid VARCHAR(150), user_id VARCHAR(50), user_pw VARCHAR(150), ACTIVE TINYINT);

CREATE TABLE chfrank_channel (cuid VARCHAR(150), channel_passphrase VARCHAR(50), ACTIVE TINYINT);

CREATE TABLE user_in_channel (uuid VARCHAR(150), cuid VARCHAR(150), cidx INT);

CREATE TABLE channel_chat_record (cuid VARCHAR(150), cidx INT, uuid VARCHAR(150), record TEXT);

INSERT INTO chfrank_user(uuid, usid, user_id, user_pw, ACTIVE) values('init_test','N','init_test','init_test', 1);


COMMIT;