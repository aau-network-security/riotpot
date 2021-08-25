#!/bin/sh


# wait for the database to be up
if [ $DB_HOST ]
then
  echo "Waiting for mongodb..."

  while ! nc -z $DB_HOST $DB_PORT; do
    sleep 0.1
  done

  echo "MongoDB started"
fi

# run riotpot
./riotpot

exec "$@"