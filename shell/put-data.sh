#!/bin/bash

echo "Do you know? That you can open http://localhost:8086/ on your browser to see the API References!"

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ] || [ -z "$4" ] || [ -z "$5" ] || [ -z "$6" ]; then
    echo "Usage: $0 cake_id name picture category comments upvotes"
    exit 1
fi

name="$2"
picture="$3"
category="$4"
comments="$5"
upvotes="$6"

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

curl -X "PUT" "http://localhost:8086/api/v1/cakes/$1" \
-H 'Content-Type: application/json' \
-d "$payload" \