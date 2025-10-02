# // filepath: d:\Folder_of_Dung\Project\Practice_Go\service\k8s\apply_envsubst.sh
#!/bin/bash
set -euo pipefail

echo "🔄 Processing Kubernetes templates..."

[ -f ".env" ] || { [ -f "../../.env" ] && ln -sf ../../.env .env || { echo "❌ .env file not found"; exit 1; }; }

set -a
. ./.env
set +a

echo "✅ Env: $APP_NAME -> $NAMESPACE"

mkdir -p base
[ -d template ] || { echo "❌ template/ missing"; exit 1; }

# List biến dùng trong template
VARS='${APP_NAME} ${NAMESPACE} ${IMAGE} ${IMAGE_PULL_POLICY} ${REPLICAS} ${CPU_REQUEST} ${CPU_LIMIT} ${MEMORY_REQUEST} ${MEMORY_LIMIT} ${DEPLOY_ENV} ${PORT} ${HPA_MIN_REPLICAS} ${HPA_MAX_REPLICAS} ${HPA_CPU_TARGET} ${HPA_MEMORY_TARGET} ${DB_HOST} ${DB_PORT} ${DB_NAME} ${DB_USER} ${DB_PASSWORD}'

process() {
  local f=$1
  if [ -f "template/$f.yaml" ]; then
    echo "  • $f"
  envsubst "$VARS" < "template/$f.yaml" > "base/$f.yaml"
  fi
}

for t in namespace postgres deployment service hpa ingress; do
  process "$t"
done

echo "✅ Done"