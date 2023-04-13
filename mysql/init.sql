CREATE USER 'seantywork'@'%' IDENTIFIED BY 'youdonthavetoknow';

GRANT ALL PRIVILEGES ON *.* TO 'seantywork'@'%';

CREATE DATABASE chfrank;

USE chfrank;


CREATE TABLE chfrank_user (uuid VARCHAR(150), usid VARCHAR(150), user_id VARCHAR(50), user_pw VARCHAR(150), ACTIVE TINYINT);

CREATE TABLE chfrank_channel (cuid VARCHAR(150), channel_name VARCHAR(50), channel_passphrase VARCHAR(50), ACTIVE TINYINT);

CREATE TABLE user_in_channel (uuid VARCHAR(150), cuid VARCHAR(150), cidx INT);

CREATE TABLE channel_chat_record (cuid VARCHAR(150), cidx INT, uuid VARCHAR(150), record TEXT);

INSERT INTO chfrank_user(uuid, usid, user_id, user_pw, ACTIVE) values('b03b606740340d6f50128c0a81c40390','N','test1','test1', 1);

INSERT INTO chfrank_user(uuid, usid, user_id, user_pw, ACTIVE) values('141f84d83505dc2261ebadf1f6424a72','N','test2','test2', 1);

INSERT INTO chfrank_channel (cuid, channel_name, channel_passphrase, ACTIVE) values('0f398445adddc9590a1fa44ada00dadc','test_chan1','N',1)

INSERT INTO chfrank_channel (cuid, channel_name, channel_passphrase, ACTIVE) values('564e07aa11f709dd6685ff7dc9aa3f53','test_chan2','secretchannel',1)



COMMIT;