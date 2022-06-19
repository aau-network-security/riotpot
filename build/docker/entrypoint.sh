#!/bin/sh

# Wait for the database to be up
if [ $DB_HOST ]
then
  echo "Waiting for database..."

  while ! nc -z $DB_HOST $DB_PORT; do
    sleep 0.1
  done

  echo "database started"
fi

# run riotpot
./riotpot

exec "$@"