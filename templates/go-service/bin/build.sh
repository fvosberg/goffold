#!/bin/bash
PROJECT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." >/dev/null 2>&1 && pwd )"
#
# Current git commit is tagged and no uncommitted changes
if [ $(git describe --exact-match --tags 2> /dev/null) ] && [ "$(git diff-index --quiet HEAD --)" == "" ] ; then
  BUILD_TAG=$(git describe --exact-match --tags 2> /dev/null || git rev-parse --short HEAD)
elif [ "$CI_PIPELINE_ID" ]; then
  BUILD_TAG="ci-${CI_COMMIT_BRANCH:-$(git symbolic-ref --short HEAD)}-${CI_PIPELINE_ID}"
else
  BUILD_TAG="dev-$(git symbolic-ref --short HEAD)-$(date -u +"%Y%m%dT%H%M%SZ")"
fi

CI_COMMIT_SHA=${CI_COMMIT_SHA:-$(git rev-parse HEAD)}
DOCKER_REPOSITORY="{{ .DockerRepository }}"
DOCKER_IMAGE="${DOCKER_REPOSITORY}:${BUILD_TAG}"


case "$1" in
  "--bin")

  export CGO_ENABLED=0
  export GO111MODULE=on
  export GOOS=linux
  BIN_NAME="${2:-{{ .CommandName }}}"

  set -x
  go build -ldflags "-X main.build=${BUILD_TAG}" -a -o "$BIN_NAME" -installsuffix cgo $PROJECT_ROOT/cmds/{{ .CommandName }}
  RETURN=$?
  set +x
  echo -e "\n\nBuilt ./${BIN_NAME}\n\n"
  exit $RETURN
  ;;

  "--container")
  set -x

  if [ "$2" != "" ] && [ "$2" != "--push" ]; then
    DOCKER_IMAGE="$2"
  fi

  docker build \
    -t "$DOCKER_IMAGE" \
    --build-arg BUILD_TAG="$BUILD_TAG" \
    --build-arg BUILD_DATE="$BUILD_DATE" \
    "$PROJECT_ROOT"

  if [ "$2" == "--push" ] || [ "$3" == "--push" ]; then
    docker push "$DOCKER_IMAGE"
    docker tag "$DOCKER_IMAGE" "$DOCKER_REPOSITORY:latest"
    docker push "$DOCKER_REPOSITORY:latest"
  fi

  exit $?
  ;;

  "--get-image-tag")
  echo "$DOCKER_IMAGE"
  exit 0
  ;;
  *)
    echo -e "\n\nMissing command. Run \n\n\t$0 bin [binary-name]\n\t$0 container [docker-image-tag] [--push]\n\t$0 --get-image-tag\n\n"
    exit 1
    ;;
esac
