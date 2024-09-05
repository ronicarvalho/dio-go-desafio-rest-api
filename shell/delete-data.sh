#!/bin/bash

echo "Do you know? That you can open http://localhost:8086/ on your browser to see the API References!"

if [ -z "$1" ]; then
    echo "Usage: $0 cake_id"
    exit 1
fi

curl -X "DELETE" "http://localhost:8086/api/v1/cakes/$1"