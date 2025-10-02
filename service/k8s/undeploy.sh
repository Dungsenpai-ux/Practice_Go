#!/bin/bash
# File: d:\Folder_of_Dung\Project\Practice_Go\deploy\k8s\undeploy.sh

set -e

echo "🗑️ Practice Go Cleanup"

# Load environment
if [[ ! -f ".env" ]]; then
    echo "❌ .env file not found"
    exit 1
fi

export $(grep -v '^#' .env | xargs)

echo "📋 Cleaning namespace: $NAMESPACE"

# Delete resources
kubectl delete -k . 2>/dev/null || echo "⚠️ Some resources may not exist"

echo "🧹 Cleaning generated files..."
rm -rf base

echo "✅ Cleanup completed!"