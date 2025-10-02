#!/bin/bash
# File: d:\Folder_of_Dung\Project\Practice_Go\deploy\k8s\setup.sh

set -e

echo "🚀 Setting up Kubernetes deployment structure..."

# Create directories
mkdir -p template base

# Create template files
echo "📄 Creating template files..."

# namespace.yaml
cat > template/namespace.yaml << 'EOF'
apiVersion: v1
kind: Namespace
metadata:
  name: ${NAMESPACE}
  labels:
    name: ${NAMESPACE}
    environment: ${DEPLOY_ENV}
    app: ${APP_NAME}
EOF

# deployment.yaml
cat > template/deployment.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${APP_NAME}
  namespace: ${NAMESPACE}
  labels:
    app: ${APP_NAME}
    version: v1
    environment: ${DEPLOY_ENV}
spec:
  replicas: ${REPLICAS}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: ${APP_NAME}
  template:
    metadata:
      labels:
        app: ${APP_NAME}
        version: v1
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "${PORT}"
        prometheus.io/path: "/metrics"
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1001
        fsGroup: 1001
      containers:
      - name: ${APP_NAME}
        image: ${IMAGE}
        imagePullPolicy: ${IMAGE_PULL_POLICY}
        ports:
        - name: http
          containerPort: ${PORT}
          protocol: TCP
        env:
        - name: LOG_LEVEL
          value: "${LOG_LEVEL}"
        - name: GIN_MODE
          value: "${GIN_MODE}"
        - name: PORT
          value: "${PORT}"
        resources:
          requests:
            memory: "${MEMORY_REQUEST}"
            cpu: "${CPU_REQUEST}"
          limits:
            memory: "${MEMORY_LIMIT}"
            cpu: "${CPU_LIMIT}"
        livenessProbe:
          httpGet:
            path: /healthz
            port: ${PORT}
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /healthz
            port: ${PORT}
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        startupProbe:
          httpGet:
            path: /healthz
            port: ${PORT}
          initialDelaySeconds: 15
          periodSeconds: 10
          failureThreshold: 30
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: false
          capabilities:
            drop:
            - ALL
EOF

# service.yaml
cat > template/service.yaml << 'EOF'
apiVersion: v1
kind: Service
metadata:
  name: ${APP_NAME}
  namespace: ${NAMESPACE}
  labels:
    app: ${APP_NAME}
    service: ${APP_NAME}
spec:
  type: ClusterIP
  ports:
  - name: http
    port: 80
    targetPort: ${PORT}
    protocol: TCP
  selector:
    app: ${APP_NAME}
EOF

# hpa.yaml
cat > template/hpa.yaml << 'EOF'
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: ${APP_NAME}-hpa
  namespace: ${NAMESPACE}
  labels:
    app: ${APP_NAME}
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: ${APP_NAME}
  minReplicas: ${HPA_MIN_REPLICAS}
  maxReplicas: ${HPA_MAX_REPLICAS}
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: ${HPA_CPU_TARGET}
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: ${HPA_MEMORY_TARGET}
EOF

# kustomization.yaml
cat > kustomization.yaml << 'EOF'
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: practice-go-local

resources:
- base/namespace.yaml
- base/deployment.yaml
- base/service.yaml
- base/hpa.yaml

commonLabels:
  project: practice-go
  managed-by: kustomize
  team: backend

images:
- name: practice-go
  newTag: latest
EOF

# Make scripts executable
chmod +x *.sh

echo "✅ Kubernetes deployment structure created successfully!"
echo ""
echo "📁 Structure:"
echo "├── .env"
echo "├── setup.sh"
echo "├── deploy.sh"
echo "├── debug.sh"
echo "├── undeploy.sh"
echo "├── apply_envsubst.sh"
echo "├── kustomization.yaml"
echo "├── template/"
echo "│   ├── namespace.yaml"
echo "│   ├── deployment.yaml"
echo "│   ├── service.yaml"
echo "│   └── hpa.yaml"
echo "└── base/ (generated)"
echo ""
echo "🚀 Ready to deploy:"
echo "   ./deploy.sh"