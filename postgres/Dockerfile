FROM postgres:14-alpine

EXPOSE 5432

RUN mkdir -p /docker-entrypoint-initdb.d
COPY ./init.sql /docker-entrypoint-initdb.d
