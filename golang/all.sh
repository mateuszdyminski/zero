#!/bin/bash

./dev.sh build simple v1
./dev.sh docker simple v1
./dev.sh build graceful v2
./dev.sh docker graceful v2
./dev.sh build simple v3
./dev.sh docker simple v3