#!/bin/bash
export $(cat conf/env.dev |xargs)
./bin/example-app
