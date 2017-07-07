#!/bin/bash

HOST=localhost
#HOST=107.23.249.68
PORT=7676

curl -v -H IMPORTANTSTUFF:YEAH $HOST:$PORT/example-webservice/v1/hello
