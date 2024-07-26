FROM ghcr.io/kukymbr/goose-docker:3.21.1

WORKDIR /root

ADD migrations/*.sql /migrations/

ARG DB_USER
ARG DB_PASSWORD
ARG DB_PORT
ARG DB_DATABASE_NAME
ARG DB_HOST
ARG MIGRATION_DIR
ARG SSL_MODE
ENV DBSTRING="host=$DB_HOST user=$DB_USER password=$DB_PASSWORD dbname=$DB_DATABASE_NAME sslmode=$SSL_MODE"

RUN goose -dir ${MIGRATION_DIR} postgres "$DBSTRING" up

