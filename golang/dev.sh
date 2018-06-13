#!/bin/bash

usage() {
	cat <<EOF
Usage: $(basename $0) <command> <server-type> <version>

Wrappers around core binaries:
    build                   Builds the zero app.
    docker                  Builds docker image and pushes it to DockerHub.
EOF
	exit 1
}

if [ "$#" -lt 3 ]; then
  echo "Pass parameters: <command> <server-type> <version>"
  usage
  exit 1
fi

CMD="$1"
TYPE="$2"
VERSION="$3"

GIT_VERSION=$(git describe --always)
LAST_COMMIT_USER="$(tr -d '[:space:]' <<<"$(git log -1 --format=%cn)<$(git log -1 --format=%ce)>")"
LAST_COMMIT_HASH=$(git log -1 --format=%H)
LAST_COMMIT_TIME=$(git log -1 --format=%cd --date=format:'%Y-%m-%d_%H:%M:%S')

LDFLAGS="-X main.appVersion=$VERSION -X main.gitVersion=$GIT_VERSION -X main.lastCommitTime=$LAST_COMMIT_TIME -X main.lastCommitHash=$LAST_COMMIT_HASH -X main.lastCommitUser=$LAST_COMMIT_USER -X main.buildTime=$(date -u +%Y-%m-%d_%H:%M:%S)"
		
shift
case "$CMD" in
	build)
		cd $TYPE; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o main -a -tags netgo .
	;;
	docker)
		docker build --build-arg TYPE=$TYPE -t "mateuszdyminski/zero-golang:$VERSION" . && docker push "mateuszdyminski/zero-golang:$VERSION"
	;;
	*)
		usage
	;;
esac