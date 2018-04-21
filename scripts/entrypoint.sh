#!/bin/sh

trap "exit" INT TERM
trap "kill 0" EXIT

nginx -g "daemon off;" &
export NGINX_PID=$!

/root/main &
export MAIN_PID=$!

wait "$MAIN_PID"