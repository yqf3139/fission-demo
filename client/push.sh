#!/bin/bash

set -e

tag=$1
if [ -z "$tag" ]
then
    tag=latest
fi

. ./build.sh

docker build -t fission-demo-client .
docker tag fission-demo-client yqf3139/fission-demo-client:$tag
docker push yqf3139/fission-demo-client:$tag
