#!/bin/sh
set -e

# Hanya jalankan service yang diinginkan, hapus folder service yang tidak aktif
[ "$WEB" != "true" ] && rm -rf /etc/services.d/server
[ "$WORKER" != "true" ] && rm -rf /etc/services.d/worker
[ "$SCHEDULER" != "true" ] && rm -rf /etc/services.d/scheduler

# Jalankan s6-overlay
exec /init
