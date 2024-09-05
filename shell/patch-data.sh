#!/bin/bash

echo "Do you know? That you can open http://localhost:8086/ on your browser to see the API References!"

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then
    echo "Usage: $0 cake_id comments upvotes"
    exit 1
fi

comments="$2"
upvotes="$3"

payload=$(cat <<EOF
{ 
    "comments": $comments,
    "upvotes": $upvotes
}
EOF
)

curl -X "PATCH" "http://localhost:8086/api/v1/cakes/$1" \
-H 'Content-Type: application/json' \
-d "$payload" \