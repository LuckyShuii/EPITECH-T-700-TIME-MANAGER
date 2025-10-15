#!/bin/bash

PWD=$(pwd)

# Create logs directory if it doesn't exist
mkdir -p "$PWD/../logs"

# Create empty log files if they don't exist
touch "$PWD/../logs/frontend.log"
touch "$PWD/../logs/backend.log"
touch "$PWD/../logs/database.log"
touch "$PWD/../logs/apigateway.log"