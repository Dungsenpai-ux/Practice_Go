#!/bin/bash
# File: d:\Folder_of_Dung\Project\Practice_Go\deploy\k8s\deploy.sh

set -e

echo "🚀 ====================================="
echo "   Practice Go Kubernetes Deployment"
echo "====================================="

# Check prerequisites
echo "🔍 Checking prerequisites..."

if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl not found"
    echo "💡 Install: https://kubernetes.io/docs/tasks/tools/install-kubectl/"
    exit 1
fi

if ! command -v docker &> /dev/null; then
    echo "❌ Docker not found"
    echo "💡 Install Docker"
    exit 1
fi

if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Cannot connect to Kubernetes cluster"
    echo "💡 Start cluster: minikube start"
    exit 1
fi

echo "✅ Prerequisites check passed"

# Load environment
if [[ ! -f ".env" ]]; then
    echo "❌ .env file not found"
    exit 1
fi

export $(grep -v '^#' .env | xargs)

echo "📋 Environment:"
echo "  • Image: $IMAGE"
echo "  • Namespace: $NAMESPACE"
echo "  • Deploy Env: $DEPLOY_ENV"

# Build Docker image
echo ""
echo "🐳 Building Docker image..."
cd ../..
if ! docker build -t "$IMAGE" .; then
    echo "❌ Docker build failed"
    exit 1
fi
echo "✅ Docker build successful"

# Load image to minikube if needed
if minikube status &> /dev/null; then
    echo "📦 Loading image to minikube..."
    minikube image load "$IMAGE"
fi

cd service/k8s

# Process templates
echo ""
echo "🔄 Processing templates..."
if ! ./apply_envsubst.sh; then
    echo "❌ Template processing failed"
    exit 1
fi

# Validate manifests
echo "🔍 Validating manifests..."
if ! kubectl apply --dry-run=client -k . &> /dev/null; then
    echo "❌ Invalid manifests"
    exit 1
fi

# Apply manifests
echo ""
echo "☸️ Applying Kubernetes manifests..."
if ! kubectl apply -k .; then
    echo "❌ Deployment failed"
    exit 1
fi

# Wait for deployment
echo ""
echo "⏳ Waiting for deployment to be ready..."
if ! kubectl wait --for=condition=available --timeout=300s deployment/"$APP_NAME" -n "$NAMESPACE"; then
    echo "❌ Deployment timeout"
    echo "🔍 Debug info:"
    kubectl get pods -n "$NAMESPACE"
    kubectl describe deployment "$APP_NAME" -n "$NAMESPACE"
    exit 1
fi

echo ""
echo "✅ ====================================="
echo "   Deployment Completed Successfully!"
echo "====================================="

echo ""
echo "📊 Status:"
kubectl get all -n "$NAMESPACE"

echo ""
echo "🔗 Access your application:"
echo "  • Port forward: ./debug.sh"
echo "  • Health check: curl http://localhost:8080/healthz"
echo ""
echo "🐛 Debug commands:"
echo "  • Logs: kubectl logs -f deployment/$APP_NAME -n $NAMESPACE"
echo "  • Pods: kubectl get pods -n $NAMESPACE"
echo "  • Events: kubectl get events -n $NAMESPACE --sort-by='.lastTimestamp'"