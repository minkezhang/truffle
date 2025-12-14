#!/usr/bin/env bash

while true; do
    go build -o .build/tui && pkill -f '_build/tui'
    inotifywait -e attrib $(find . -name '*.go') || exit
done
