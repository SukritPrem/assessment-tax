#!bin/bash

export	DB_HOST=localhost
export  DATABASE_URL=postgres
export	DB_PORT=5432
export	DB_USER=postgres
export	DB_PASSWORD=postgres
export	DB_NAME=ktaxes
export	DB_SSL_MODE=disable
export  ADMIN_USERNAME=adminTax
export  ADMIN_PASSWORD=admin!
export  PORT=8080

exec "$@"