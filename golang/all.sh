#!/bin/bash

./dev.sh build simple v1
./dev.sh docker simple v1
./dev.sh build simple v2
./dev.sh docker simple v2