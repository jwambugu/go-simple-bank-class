#!/bin/sh

set -e

echo "running db migrations.."
/app/migrate -path /app/migrations -database "$DB_SOURCE" --verbose up

echo "starting the app.."
exec "$@"