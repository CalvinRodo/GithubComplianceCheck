#!/bin/bash

pushd ..

docker build . -t githubcompliance:latest

popd 