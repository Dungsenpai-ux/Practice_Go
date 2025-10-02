# Kubernetes Deployment Guide

## Mục tiêu
Deploy Practice Go microservice lên Kubernetes (Minikube)

## Yêu cầu
- Minikube đã cài đặt
- kubectl đã cài đặt  
- Docker đã cài đặt

## Các bước deploy

### 1. Build Docker image
```powershell
docker build -t practice-go:latest .
```

### 2. Load image vào Minikube
```powershell
minikube image load practice-go:latest
```

### 3. Deploy dependencies (PostgreSQL + Memcached)
```powershell
kubectl apply -f k8s/dependencies.yaml
```

### 4. Deploy ConfigMap và Secret
```powershell
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
```

### 5. Deploy ứng dụng
```powershell
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

### 6. Kiểm tra trạng thái
```powershell
# Xem pods
kubectl get pods

# Xem services
kubectl get svc

# Xem logs
kubectl logs deployment/practice-go
```

### 7. Truy cập ứng dụng

**Cách 1: NodePort**
```powershell
minikube service practice-go-service --url
```

**Cách 2: Port forward**
```powershell
kubectl port-forward service/practice-go-service 8080:80
```

Sau đó truy cập:
- Health check: http://localhost:8080/healthz
- Metrics: http://localhost:8080/metrics

## Health Probes

### Readiness Probe
- Kiểm tra xem pod đã sẵn sàng nhận traffic chưa
- Endpoint: `/healthz`
- Initial delay: 5s
- Period: 5s

### Liveness Probe  
- Kiểm tra xem pod còn sống không
- Endpoint: `/healthz`
- Initial delay: 15s
- Period: 10s

## Xóa tất cả resources
```powershell
kubectl delete -f k8s/
```

## Troubleshooting

### Pods bị CrashLoopBackOff
```powershell
# Xem logs
kubectl logs deployment/practice-go

# Describe pod để xem chi tiết
kubectl describe pod <pod-name>
```

### Không kết nối được database
```powershell
# Kiểm tra postgres đã chạy chưa
kubectl get pods -l app=postgres

# Xem logs postgres
kubectl logs deployment/postgres
```

### Image pull error
```powershell
# Load lại image vào minikube
minikube image load practice-go:latest

# Hoặc list images trong minikube
minikube image ls | grep practice-go
```
