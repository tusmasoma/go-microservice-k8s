FROM mysql:8

COPY ./init.d/* /docker-entrypoint-initdb.d/
COPY ./my.cnf /etc/mysql/conf.d/my.cnf