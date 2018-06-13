#!/bin/bash

sed -i -e 's/info.app.version=v./info.app.version=v1/g' src/main/resources/application.properties && ./dev.sh build && ./dev.sh docker v1
sed -i -e 's/info.app.version=v./info.app.version=v2/g' src/main/resources/application.properties && ./dev.sh build && ./dev.sh docker v2 Dockerfile-graceful
sed -i -e 's/info.app.version=v./info.app.version=v3/g' src/main/resources/application.properties && ./dev.sh build && ./dev.sh docker v3