#!/bin/bash

go build

for i in 1 2 3 4 5; do
    echo "./run -id $i"
    x-terminal-emulator --new-tab -x "./run -id $i"
done
