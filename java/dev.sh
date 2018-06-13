#!/bin/bash

usage() {
	cat <<EOF
Usage: $(basename $0) <command> <version>

Wrappers around core binaries:
    build                   Builds the zero app.
    docker                  Builds docker image and pushes it to DockerHub.
EOF
	exit 1
}

if [ "$#" -lt 1 ]; then
  echo "Pass parameters: <command> <version> <dockerfile>"
  usage
  exit 1
fi

CMD="$1"
VERSION="$2"
DOCKERFILE="${3:-Dockerfile}"

shift
case "$CMD" in
	build)
		mvn package
	;;
	docker)
		docker build -t "mateuszdyminski/zero-java:$VERSION" -f $DOCKERFILE . && docker push "mateuszdyminski/zero-java:$VERSION"
	;;
	*)
		usage
	;;
esac