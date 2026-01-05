#!/bin/bash

PWD=$(pwd)

# Create logs directory if it doesn't exist
mkdir -p "$PWD/../logs"
mkdir -p "$PWD/../logs/archives"

mkdir -p "$PWD/../logs/archives/frontend"
mkdir -p "$PWD/../logs/archives/backend"
mkdir -p "$PWD/../logs/archives/database"
mkdir -p "$PWD/../logs/archives/apigateway"

# Create empty log files if they don't exist
touch "$PWD/../logs/frontend.log"
touch "$PWD/../logs/backend.log"
touch "$PWD/../logs/database.log"
touch "$PWD/../logs/apigateway.log"

# Create empty folder for KPI exports csv files
mkdir "$PWD/../backend/data/kpi