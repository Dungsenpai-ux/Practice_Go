#!/bin/bash
# File: d:\Folder_of_Dung\Project\Practice_Go\deploy\k8s\debug.sh

set -e

echo "🐛 Practice Go Debug Mode"

# Load environment
if [[ ! -f ".env" ]]; then
    echo "❌ .env file not found"
    exit 1
fi

export $(grep -v '^#' .env | xargs)

echo "🔗 Port forwarding $APP_NAME service to http://localhost:8080"
echo "🏥 Health: http://localhost:8080/healthz"
echo "📊 Metrics: http://localhost:8080/metrics"
echo ""
echo "Press Ctrl+C to stop"

kubectl port-forward svc/"$APP_NAME" 8080:80 -n "$NAMESPACE"