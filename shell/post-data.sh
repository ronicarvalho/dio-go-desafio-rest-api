#!/bin/bash

echo "Do you know? That you can open http://localhost:8086/ on your browser to see the API References!"

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ] || [ -z "$4" ] || [ -z "$5" ]; then
    echo "Usage: $0 name picture category comments upvotes"
    exit 1
fi

name="$1"
picture="$2"
category="$3"
comments="$4"
upvotes="$5"

payload=$(cat <<EOF
{ 
    "name": "$name",
    "picture": "$picture",
    "category": "$category",
    "comments": $comments,
    "upvotes": $upvotes
}
EOF
)

curl -X "POST" "http://localhost:8086/api/v1/cakes" \
-H 'Content-Type: application/json' \
-d "$payload" \