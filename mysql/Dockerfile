FROM mysql
MAINTAINER seantywork@gmail.com

ENV MYSQL_ROOT_PASSWORD youdonthavetoknow
ENV MYSQL_HOST localhost

EXPOSE 3306

ADD ./init.sql /docker-entrypoint-initdb.d