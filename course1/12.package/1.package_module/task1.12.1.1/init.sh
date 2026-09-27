#!/bin/bash

if [ $# -eq 0 ]; then
    echo "Module name argument is missing"
    exit 1
fi

go mod init "$1"

if [ $? -ne 0 ]; then
    exit 1
fi
