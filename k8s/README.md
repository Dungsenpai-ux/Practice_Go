# Kubernetes Deployment Guide

## Mục tiêu
Deploy Practice Go microservice lên Kubernetes (Minikube) với PostgreSQL, Memcached, Tempo (Tracing), và Grafana

## Architecture
```
┌─────────────────────────────────────────────────────────────┐
│                      Minikube Cluster                       │
│                                                             │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐ │
│  │   Grafana    │───▶│    Tempo     │◀───│Practice Go   │ │
│  │  (Port 3000) │    │(Port 3200,   │    │   App        │ │
│  │              │    │     4318)    │    │(Port 8080)   │ │
│  └──────────────┘    └──────────────┘    └──────┬───────┘ │
│                                                  │         │
│                                        ┌─────────┴─────┐   │
│                                        │               │   │
│                              ┌─────────▼──┐   ┌───────▼──┐│
│                              │ PostgreSQL │   │Memcached ││
│                              │(Port 5432) │   │(Port     ││
│                              │            │   │ 11211)   ││
│                              └────────────┘   └──────────┘│
└─────────────────────────────────────────────────────────────┘
```

## Yêu cầu
- ✅ Minikube đã cài đặt và running
- ✅ kubectl đã cài đặt
- ✅ Docker đã cài đặt
- 💾 Tối thiểu 4GB RAM cho Minikube
- 💻 2 CPU cores

## Cấu trúc thư mục k8s/
```
k8s/
├── configmap.yaml       # Application configuration (non-sensitive)
├── secret.yaml          # Credentials (DB password, etc.)
├── postgres.yaml        # PostgreSQL database (with PVC)
├── memcache.yaml        # Memcached cache
├── tempo.yaml           # Tempo distributed tracing
├── deployment.yaml      # Practice Go application
├── service.yaml         # NodePort service (port 30081)
├── grafana.yaml         # Grafana UI (port 30300)
└── README.md            # This file
```

## Quick Start - Deploy từ đầu

### Bước 1: Khởi động Minikube
```powershell
# Xóa cluster cũ (nếu có)
minikube delete

# Start Minikube với đủ resources
minikube start --driver=docker --cpus=4 --memory=4096

# Verify cluster
kubectl cluster-info
kubectl get nodes
```

### Bước 2: Build và load Docker image
```powershell
# Build image
docker build -t practice-go:latest .

# Load vào Minikube
minikube image load practice-go:latest

# Verify image đã load
minikube image ls | findstr practice-go
```

### Bước 3: Deploy theo thứ tự

**QUAN TRỌNG: Deploy đúng thứ tự để tránh lỗi!**

```powershell
# 1. Secret và ConfigMap (cần deploy trước)
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/configmap.yaml

# 2. Dependencies (PostgreSQL, Memcached, Tempo)
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/memcache.yaml
kubectl apply -f k8s/tempo.yaml

# Đợi dependencies ready (quan trọng!)
kubectl wait --for=condition=ready pod -l app=postgres --timeout=120s
kubectl wait --for=condition=ready pod -l app=memcached --timeout=60s
kubectl wait --for=condition=ready pod -l app=tempo --timeout=60s

# 3. Application
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# 4. Grafana (optional - để xem traces)
kubectl apply -f k8s/grafana.yaml
```

### Bước 4: Kiểm tra trạng thái
```powershell
# Xem tất cả pods (tất cả phải Running)
kubectl get pods

# Kết quả mong muốn:
# NAME                           READY   STATUS    RESTARTS   AGE
# postgres-xxx                   1/1     Running   0          2m
# memcached-xxx                  1/1     Running   0          2m
# tempo-xxx                      1/1     Running   0          2m
# practice-go-xxx                1/1     Running   0          1m
# grafana-xxx                    1/1     Running   0          1m

# Xem services
kubectl get svc

# Xem endpoints (phải có IP)
kubectl get endpoints practice-go-service
```

### Bước 5: Xem logs real-time
```powershell
# Logs real-time với timestamps
kubectl logs -f deployment/practice-go --timestamps

# Logs của tất cả pods
kubectl logs -f -l app=practice-go --all-containers=true --prefix=true
```

## Truy cập ứng dụng

### Cách 1: Minikube Service (Khuyến nghị)
```powershell
# Lấy URL và tự động mở browser
minikube service practice-go-service

# Chỉ lấy URL
minikube service practice-go-service --url
# Output: http://127.0.0.1:xxxxx
```

### Cách 2: Port Forward (Dùng localhost)
```powershell
# Forward sang localhost:8080
kubectl port-forward service/practice-go-service 8080:80
```

### Cách 3: NodePort trực tiếp
```powershell
# Lấy Minikube IP
minikube ip

# Truy cập: http://<minikube-ip>:30081
```

## Test API

### Endpoints có sẵn:

**Health Check:**
```bash
GET http://<url>/healthz
```

**Search Movies:**
```bash
GET http://<url>/movies/search?q=inception
```

**Get Movie by ID:**
```bash
GET http://<url>/movies/1
```

**Create Movie:**
```bash
POST http://<url>/movies
Content-Type: application/json

{
  "title": "The Dark Knight",
  "genres": "Action,Crime,Drama",
  "year": 2008
}
```

**Metrics:**
```bash
GET http://<url>/metrics
GET http://<url>/debug/vars
```


## Grafana - Xem Distributed Traces

### Truy cập Grafana

**Cách 1: Port-forward (Khuyến nghị)**
```powershell
kubectl port-forward service/grafana 3000:3000
```
Truy cập: **http://localhost:3000**

**Cách 2: Minikube Service**
```powershell
minikube service grafana --url
```

**Login:**
- Username: `admin`
- Password: `admin`

### Xem Traces trong Grafana

1. **Vào Explore** (biểu tượng la bàn ở sidebar trái)
2. **Chọn Data Source**: Tempo
3. **Query Type**: Search
4. **Service Name**: `practice-go-api`
5. Click **Run query**

Bạn sẽ thấy danh sách traces với:
- ⏱️ Duration
- 🏷️ HTTP method, status code
- 📊 Span timeline

### Tạo traffic để có traces

Test nhiều API calls trong Postman để tạo traces:
```bash
GET /healthz
GET /movies/search?q=inception
GET /movies/1
POST /movies
```

Sau đó xem traces trong Grafana Explore.

## Health Probes

### Readiness Probe
- ✅ Kiểm tra pod đã sẵn sàng nhận traffic
- Endpoint: `GET /healthz`
- Initial delay: 5 seconds
- Period: 5 seconds
- Timeout: 3 seconds

### Liveness Probe
- ❤️ Kiểm tra pod còn sống không
- Endpoint: `GET /healthz`
- Initial delay: 15 seconds
- Period: 10 seconds
- Timeout: 3 seconds

Nếu health check fail, Kubernetes sẽ tự động restart pod.

## Monitoring & Debugging

### Xem logs real-time
```powershell
# Logs với timestamps
kubectl logs -f deployment/practice-go --timestamps

# Logs tất cả pods với prefix
kubectl logs -f -l app=practice-go --all-containers=true --prefix=true

# Logs một pod cụ thể
kubectl logs -f <pod-name>

# Xem logs trước đó (nếu pod restart)
kubectl logs --previous <pod-name>
```

### Ví dụ logs khi call API:
```
2025-10-08T14:20:10.123456Z [otel] tracing enabled protocol=http endpoint=http://tempo:4318 service=practice-go-api sampler=parentbased_always_on
2025-10-08T14:20:15.234567Z [GIN] 2025/10/08 - 14:20:15 | 200 |     1.234ms | 10.244.0.1 | GET      "/healthz"
2025-10-08T14:20:20.345678Z [GIN] 2025/10/08 - 14:20:20 | 200 |    45.678ms | 10.244.0.1 | GET      "/movies/search?q=inception"
2025-10-08T14:20:25.456789Z [GIN] 2025/10/08 - 14:20:25 | 200 |    12.345ms | 10.244.0.1 | GET      "/movies/1"
2025-10-08T14:20:30.567890Z [GIN] 2025/10/08 - 14:20:30 | 201 |    78.901ms | 10.244.0.1 | POST     "/movies"
```

### Kiểm tra chi tiết pods
```powershell
# Describe pod để xem events và config
kubectl describe pod <pod-name>

# Exec vào pod
kubectl exec -it deployment/practice-go -- sh

# Trong pod, test kết nối:
nc -zv postgres 5432
nc -zv memcached 11211
nc -zv tempo 4318
exit
```

### Xem resource usage
```powershell
# CPU và Memory usage của pods
kubectl top pods

# Nodes resource usage
kubectl top nodes
```

### Xem events
```powershell
# Events của cluster (sắp xếp theo thời gian)
kubectl get events --sort-by=.metadata.creationTimestamp

# Events của một pod cụ thể
kubectl describe pod <pod-name> | Select-String -Pattern "Events:" -Context 0,20
```

## Quản lý Deployment

### Scale replicas
```powershell
# Scale lên 3 replicas
kubectl scale deployment/practice-go --replicas=3

# Verify
kubectl get pods -l app=practice-go
```

### Update image
```powershell
# Rebuild image
docker build -t practice-go:latest .

# Load vào Minikube
minikube image load practice-go:latest

# Delete pods để dùng image mới
kubectl delete pods -l app=practice-go

# Hoặc dùng rollout restart
kubectl rollout restart deployment/practice-go
```

### Rollback deployment
```powershell
# Xem history
kubectl rollout history deployment/practice-go

# Rollback về version trước
kubectl rollout undo deployment/practice-go

# Rollback về revision cụ thể
kubectl rollout undo deployment/practice-go --to-revision=2
```

### Check rollout status
```powershell
kubectl rollout status deployment/practice-go
```

## Xóa resources

### Xóa tất cả
```powershell
kubectl delete -f k8s/
```

### Xóa từng loại
```powershell
# Xóa app
kubectl delete -f k8s/deployment.yaml
kubectl delete -f k8s/service.yaml

# Xóa dependencies
kubectl delete -f k8s/postgres.yaml
kubectl delete -f k8s/memcache.yaml
kubectl delete -f k8s/tempo.yaml
kubectl delete -f k8s/grafana.yaml

# Xóa config
kubectl delete -f k8s/configmap.yaml
kubectl delete -f k8s/secret.yaml
```

## Troubleshooting

### 1. Pods bị CrashLoopBackOff
```powershell
# Xem logs để tìm lỗi
kubectl logs deployment/practice-go

# Xem logs của container trước đó
kubectl logs <pod-name> --previous

# Describe pod để xem events
kubectl describe pod <pod-name>
```

**Nguyên nhân thường gặp:**
- ❌ Database chưa ready
- ❌ Migration lỗi
- ❌ File CSV không tồn tại
- ❌ Biến môi trường sai

### 2. Không kết nối được database
```powershell
# Kiểm tra Postgres pod
kubectl get pods -l app=postgres
kubectl logs deployment/postgres

# Test kết nối từ app pod
kubectl exec -it deployment/practice-go -- sh
nc -zv postgres 5432
exit
```

**Fix:**
- Đợi Postgres ready trước khi deploy app
- Kiểm tra Secret có key `db-password`

### 3. Image pull error / ImagePullBackOff
```powershell
# Verify image trong Minikube
minikube image ls | findstr practice-go

# Nếu không có, load lại
minikube image rm practice-go:latest
docker build -t practice-go:latest .
minikube image load practice-go:latest

# Delete pods để recreate
kubectl delete pods -l app=practice-go
```

**Lưu ý:** Trong deployment.yaml phải có `imagePullPolicy: Never` để dùng image local.

### 4. Service không accessible
```powershell
# Kiểm tra service
kubectl get svc practice-go-service
kubectl describe svc practice-go-service

# Kiểm tra endpoints (phải có IP)
kubectl get endpoints practice-go-service

# Nếu endpoints trống, kiểm tra pods
kubectl get pods -l app=practice-go

# Test từ bên trong cluster
kubectl run curl --image=curlimages/curl -it --rm -- sh
curl http://practice-go-service/healthz
exit
```

### 5. Tempo connection refused
```powershell
# Kiểm tra Tempo pod
kubectl get pods -l app=tempo
kubectl logs deployment/tempo

# Restart Tempo
kubectl delete pod -l app=tempo

# Verify service
kubectl get svc tempo
kubectl get endpoints tempo
```

### 6. Minikube cluster issues
```powershell
# Check cluster status
minikube status

# Restart Minikube
minikube stop
minikube start

# Nếu vẫn lỗi, delete và start lại
minikube delete
minikube start --driver=docker --cpus=4 --memory=4096
```

### 7. Secret key not found
```powershell
# Kiểm tra secret
kubectl get secret practice-go-secret -o yaml

# Verify keys
kubectl describe secret practice-go-secret
```

Phải có các keys:
- `DB_USER`
- `DB_PASSWORD`
- `db-password` (cho Postgres deployment)
- `DB_URL`

## Best Practices

### 1. Deploy Order
Luôn deploy theo thứ tự:
1. Secret & ConfigMap
2. Databases & Dependencies
3. Application
4. Monitoring tools

### 2. Health Checks
- ✅ Luôn có readinessProbe và livenessProbe
- ⏱️ Set initialDelaySeconds phù hợp với startup time
- 🔄 periodSeconds không quá thấp (tránh spam)

### 3. Resource Limits
- 💾 Set requests và limits cho tất cả containers
- 📊 Monitor resource usage: `kubectl top pods`
- 🔧 Adjust dựa trên actual usage

### 4. Logging
- 📝 Luôn xem logs khi deploy: `kubectl logs -f`
- 🐛 Debug với timestamps: `--timestamps`
- 📊 Export logs để analyze: `kubectl logs deployment/practice-go > app.log`

### 5. Image Management
- 🏷️ Tag images properly (không chỉ dùng `latest`)
- 🔄 Verify image loaded: `minikube image ls`
- 🗑️ Clean old images: `minikube image rm`

## Useful Commands Cheat Sheet

```powershell
# --- Pod Management ---
kubectl get pods                              # List all pods
kubectl get pods -o wide                      # List with more details
kubectl get pods -w                           # Watch pod changes
kubectl describe pod <pod-name>               # Detailed info
kubectl logs -f <pod-name>                    # Stream logs
kubectl exec -it <pod-name> -- sh             # Shell into pod
kubectl delete pod <pod-name>                 # Delete pod

# --- Service Management ---
kubectl get svc                               # List services
kubectl describe svc <service-name>           # Service details
kubectl get endpoints                         # Service endpoints
minikube service <service-name> --url         # Get service URL

# --- Deployment Management ---
kubectl get deployments                       # List deployments
kubectl scale deployment/<name> --replicas=3  # Scale
kubectl rollout restart deployment/<name>     # Restart
kubectl rollout status deployment/<name>      # Check status
kubectl rollout history deployment/<name>     # View history
kubectl rollout undo deployment/<name>        # Rollback

# --- Resource Viewing ---
kubectl get all                               # All resources
kubectl get all -l app=practice-go            # Filter by label
kubectl top pods                              # Resource usage
kubectl top nodes                             # Node usage

# --- Debugging ---
kubectl describe <resource> <name>            # Detailed info
kubectl logs -f deployment/<name>             # Stream logs
kubectl logs <pod> --previous                 # Previous logs
kubectl get events --sort-by=.metadata.creationTimestamp  # Events

# --- Minikube ---
minikube status                               # Cluster status
minikube ip                                   # Cluster IP
minikube service list                         # All services
minikube dashboard                            # Web UI
minikube image ls                             # List images
minikube image load <image>                   # Load image
minikube image rm <image>                     # Remove image
```

## Architecture Details

### Components:

**Application Tier:**
- `practice-go`: Go microservice (2 replicas)
- Health checks: `/healthz`
- Metrics: `/metrics`, `/debug/vars`

**Data Tier:**
- `PostgreSQL 16`: Primary database
- PersistentVolumeClaim: 1Gi
- Auto-creates database via `POSTGRES_DB` env

**Cache Tier:**
- `Memcached 1.6`: In-memory cache
- Port: 11211

**Observability Tier:**
- `Tempo`: Distributed tracing backend
- OpenTelemetry OTLP receiver (HTTP: 4318, gRPC: 4317)
- `Grafana`: Traces visualization UI

### Network:

**Internal (ClusterIP):**
- postgres:5432
- memcached:11211
- tempo:3200, 4317, 4318
- grafana:3000 (also exposed via NodePort)

**External (NodePort):**
- practice-go-service:30081
- tempo:30320
- grafana:30300

## Screenshots để submit cho mentor

1. ✅ `kubectl get pods` - Tất cả Running
2. ✅ `kubectl get svc` - Services
3. ✅ `kubectl logs -f deployment/practice-go --timestamps` - Logs với [GIN] entries
4. ✅ Postman - API requests & responses
5. ✅ Grafana - Tempo datasource connected
6. ✅ Grafana Explore - Traces list
7. ✅ Trace detail - Span timeline

## References

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Minikube Documentation](https://minikube.sigs.k8s.io/docs/)
- [Grafana Tempo](https://grafana.com/oss/tempo/)
- [OpenTelemetry](https://opentelemetry.io/)

---

**Created:** October 2025  
**Author:** Practice Go Team  
**Branch:** feature/Deploy_Kurbenetes
