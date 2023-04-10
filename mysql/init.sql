CREATE USER 'seantywork'@'%' IDENTIFIED BY 'youdonthavetoknow';

GRANT ALL PRIVILEGES ON *.* TO 'seantywork'@'%';

CREATE DATABASE chfrank;

USE chfrank;


CREATE TABLE chfrank_user (uuid VARCHAR(150), user_id VARCHAR(50), user_pw VARCHAR(150), ACTIVE TINYINT);



INSERT INTO chfrank_user(uuid, user_id, user_pw, ACTIVE) values('init_test','init_test','init_test', 1);


COMMIT;