#!/bin/bash

cur=`dirname $(readlink -f $0)`
cwd=`pwd`

cd ${cur}/..

bin/checkah example --format=json > configs/example.json
bin/checkah example --format=yaml > configs/example.yaml
bin/checkah example --format=json --local > configs/localhost.json
bin/checkah example --format=yaml --local > configs/localhost.yaml
