#!/bin/bash

go build

for i in 4 3 2 1; do
    if [$i -eq 1]
    then
        echo "./run -id $i -init true"
        x-terminal-emulator --new-tab -x "./run -id $i -init true"
    else
        echo "./run -id $i"
        x-terminal-emulator --new-tab -x "./run -id $i"
    fi
done
