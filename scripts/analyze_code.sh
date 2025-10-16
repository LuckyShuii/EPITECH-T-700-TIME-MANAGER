#!/bin/sh

# Charge les variables d'environnement du .env
export $(grep -v '^#' ../.env | xargs)

NETWORK_NAME="t-dev-700-project-par_19_default"
SONAR_HOST_URL="http://sonarqube:9000/sonarqube"

echo "🚀 Running SonarQube analysis..."
echo "🔗 Using network: $NETWORK_NAME"
echo "🌐 Server: $SONAR_HOST_URL"
echo "🔑 Token loaded from .env"

docker run --rm \
  --network "$NETWORK_NAME" \
  -e SONAR_HOST_URL="$SONAR_HOST_URL" \
  -e SONAR_TOKEN="$SONAR_TOKEN" \
  -v "$(pwd)/..:/usr/src" \
  sonarsource/sonar-scanner-cli:latest
