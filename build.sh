#!/bin/sh
set -e

rm -rf build/*
go build -o build/coldline-miami . && ./build/coldline-miami