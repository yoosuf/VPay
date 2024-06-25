#!/bin/sh

# wait-for-it.sh -- wait for a TCP connection to become available

HOST="$1"
PORT="$2"
shift 2
CMD="$@"

until nc -z -v -w30 $HOST $PORT
do
  echo "Waiting for database connection at $HOST:$PORT..."
  sleep 1
done

echo "Database connection available at $HOST:$PORT"
exec $CMD
