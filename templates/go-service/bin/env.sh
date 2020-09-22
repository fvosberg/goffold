#!/bin/bash
PROJECT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." >/dev/null 2>&1 && pwd )"

setup () {
  docker-compose \
      -f "$PROJECT_ROOT/test_infrastructure/docker-compose.yml" \
      up -d
}

logs () {
  docker-compose \
      -f "$PROJECT_ROOT/test_infrastructure/docker-compose.yml" \
      logs -f
}

teardown() {
  docker-compose \
      -f "$PROJECT_ROOT/test_infrastructure/docker-compose.yml" \
      down
}
