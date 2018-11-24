#!/bin/bash

pushd ../comply/

go build \
-a \
-o ../../bin/does-this-comply .

popd