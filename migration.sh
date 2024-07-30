#!/bin/bash
source .env

export MIGRATION_DSN="postgresql://$DB_USER:$DB_PASSWORD@pg:$DB_PORT/$DB_DATABASE_NAME?sslmode=disable"

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v