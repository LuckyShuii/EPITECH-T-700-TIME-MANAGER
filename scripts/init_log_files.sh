#!/bin/bash

PWD=$(pwd)

# Create logs directory if it doesn't exist
mkdir -p "$PWD/../logs"
mkdir -p "$PWD/../logs/archives"

mkdir -p "$PWD/../logs/archive/frontend"
mkdir -p "$PWD/../logs/archive/backend"
mkdir -p "$PWD/../logs/archive/database"
mkdir -p "$PWD/../logs/archive/apigateway"

# Create empty log files if they don't exist
touch "$PWD/../logs/frontend.log"
touch "$PWD/../logs/backend.log"
touch "$PWD/../logs/database.log"
touch "$PWD/../logs/apigateway.log"