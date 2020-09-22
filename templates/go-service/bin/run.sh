#!/bin/bash
PROJECT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." >/dev/null 2>&1 && pwd )"
source "$PROJECT_ROOT/bin/env.sh"

DOCKER_DEV_IMAGE=mocvt_dev

case "$1" in
  "--local"):
    go run . "$@"
    exit 0
    ;;

  "--docker-dev"):
    set -x
    set -e
    setup
    docker-compose -f "$PROJECT_ROOT"/test_infrastructure/docker-compose.yml up -d --force-recreate
    exit 0
    ;;

  "--docker-prod"):
    echo "Not implemented, yet"
    exit 1
    ;;

  *)
    echo -e "\n\nMissing command. Run $0 (command)\n\n\t--local [--web-api-host]\tlocal execution of the go binary for debuggers e.g.\n\t--docker-dev\t\texecution in a docker container, but with hot reloading\n\t--docker-prod [--build]\trun the production container\n\n"
    exit 1
esac
