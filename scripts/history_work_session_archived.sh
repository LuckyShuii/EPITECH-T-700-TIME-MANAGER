#!/bin/bash

PWD=$(pwd)
source .env

SQL_QUERY=$(cat "$PWD/sql/history_work_session_archived.sql")

docker exec -i time-manager-database psql -U "$DB_USER" -d "$DB_NAME" -c "$SQL_QUERY"
