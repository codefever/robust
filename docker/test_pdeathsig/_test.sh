#!/bin/sh
set -e

wrapper_bin="/robust/subprocess/wrapper/wrapper"
looper_bin="/robust/docker/looper.sh"

echo "Before run..." >>/dev/stderr
$wrapper_bin --command "$looper_bin >1.log 2>&1" &  # not piped
echo "After run..." >>/dev/stderr
sleep 1

wrapper_pid=$(pgrep $(basename $wrapper_bin))
looper_pid=$(pgrep $(basename $looper_bin))
echo "PID: $wrapper_pid -> $looper_pid"
sleep 1

echo "Try to kill the parent..." >>/dev/stderr
kill $wrapper_pid
sleep 3

echo "Detect whether the child is alive..." >>/dev/stderr
msg=$(kill -0 $looper_pid 2>&1 || true)
echo $msg | grep -i "no such process" >/dev/null

echo "\n@@PASS@@"
