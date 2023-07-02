FROM mysql:8

COPY ./create-tables.sql /tmp

CMD [ "mysqld", "--init-file=/tmp/create-tables.sql" ]