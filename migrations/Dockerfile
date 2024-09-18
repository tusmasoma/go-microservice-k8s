FROM --platform=linux/amd64 mysql:5.7

COPY ./init.d/* /docker-entrypoint-initdb.d/
COPY ./my.cnf /etc/mysql/conf.d/my.cnf