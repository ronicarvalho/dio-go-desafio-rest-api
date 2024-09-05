#!/bin/bash

echo "Do you know? That you can open http://localhost:8086/ on your browser to see the API References!"

curl -X "GET" "http://localhost:8086/api/v1/cakes"