#!/bin/sh

set -e
echo "Menjalankan migration database..."
migrate -path /app/migrations \
        -database "postgres://$DB_USER:$DB_PASS@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up
echo "Migration selesai."

if [ "$WEB_ONLY" == "true" ];then
    rm -f /etc/supervisor/conf.d/worker.conf
    rm -f /etc/supervisor/conf.d/scheduler.conf
fi

if [ "$SCHEDULER_ONLY" == "true" ];then
    rm -f /etc/supervisor/conf.d/server.conf
    rm -f /etc/supervisor/conf.d/worker.conf
fi

if [ "$WORKER_ONLY" == "true" ];then
    rm -f /etc/supervisor/conf.d/server.conf
    rm -f /etc/supervisor/conf.d/scheduler.conf
fi

exec $@