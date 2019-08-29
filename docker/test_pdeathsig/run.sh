#!/bin/sh
set -e
cd $(dirname $0)/../..
docker build -t test_pdeathsig -f docker/test_pdeathsig/Dockerfile .
docker run --rm test_pdeathsig /_test.sh
