#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 [folder]"
  exit 1
fi

if [ "$1" = "clear" ]; then
  rm parser
  rm to_analyse.txt
  rm analyser
  rm result.json
  exit 0
fi

go build -o parser ./src_parser
go build -o analyser ./src_analyser

./parser -path=$1
./analyser