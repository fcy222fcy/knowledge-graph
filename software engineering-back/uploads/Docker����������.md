# Docker 与容器技术完全指南

## 一、容器技术概述

### 1.1 什么是容器

容器是一种轻量级的虚拟化技术，它将应用程序及其所有依赖（库、配置文件、运行时环境）打包在一起，确保应用在任何环境中都能以相同的方式运行。与传统虚拟机不同，容器共享宿主机的操作系统内核，因此更加轻量高效。

### 1.2 容器 vs 虚拟机

| 特性 | 容器 | 虚拟机 |
|------|------|--------|
| 启动速度 | 毫秒级 | 分钟级 |
| 资源占用 | MB 级别 | GB 级别 |
| 隔离级别 | 进程级（Namespace + Cgroup） | 硬件级（Hypervisor） |
| 操作系统 | 共享宿主机内核 | 独立内核 |
| 镜像大小 | 通常 10-500MB | 通常 1-10GB |
| 性能损耗 | 接近原生（<2%） | 5-20% |
| 密度 | 单机可运行数百容器 | 单机通常数十虚拟机 |

### 1.3 容器核心技术

**Linux Namespace（命名空间）**：实现资源隔离

| Namespace | 隔离内容 | 系统调用参数 |
|-----------|---------|-------------|
| PID | 进程 ID | CLONE_NEWPID |
| NET | 网络设备、IP 地址、端口 | CLONE_NEWNET |
| MNT | 文件系统挂载点 | CLONE_NEWNS |
| UTS | 主机名和域名 | CLONE_NEWUTS |
| IPC | 进程间通信 | CLONE_NEWIPC |
| USER | 用户和组 ID | CLONE_NEWUSER |

**Linux Cgroup（控制组）**：实现资源限制

- CPU 使用率限制
- 内存使用限制
- 磁盘 I/O 带宽限制
- 网络带宽限制

---

## 二、Docker 核心概念

### 2.1 镜像（Image）

Docker 镜像是一个只读模板，包含运行应用所需的一切：代码、运行时、库、环境变量、配置文件。

**镜像分层机制**：
```
┌─────────────────────────┐
│   应用代码层 (可写)       │  ← 容器层
├─────────────────────────┤
│   依赖安装层 (只读)       │  ← Dockerfile RUN
├─────────────────────────┤
│   基础镜像层 (只读)       │  ← Dockerfile FROM
└─────────────────────────┘
```

每一层都是只读的，当容器启动时，在最上面添加一个可写层。多个容器可以共享下面的只读层，节省磁盘空间。

**常用镜像操作**：
```bash
# 拉取镜像
docker pull nginx:1.24

# 查看本地镜像
docker images

# 构建镜像
docker build -t myapp:v1.0 .

# 导出/导入镜像
docker save -o myapp.tar myapp:v1.0
docker load -i myapp.tar

# 删除镜像
docker rmi myapp:v1.0

# 推送到仓库
docker push myrepo/myapp:v1.0
```

### 2.2 容器（Container）

容器是镜像的运行实例。一个镜像可以创建多个容器，每个容器都有自己的可写层。

```bash
# 运行容器
docker run -d --name web -p 80:80 nginx:1.24

# 查看运行中的容器
docker ps

# 查看所有容器（包括已停止的）
docker ps -a

# 进入容器
docker exec -it web /bin/bash

# 查看容器日志
docker logs -f web

# 停止/启动容器
docker stop web
docker start web

# 删除容器
docker rm -f web
```

### 2.3 数据卷（Volume）

数据卷用于持久化容器数据，即使容器被删除，数据仍然保留。

```bash
# 创建数据卷
docker volume create mydata

# 挂载数据卷
docker run -v mydata:/app/data myapp:v1.0

# 挂载宿主机目录
docker run -v /host/path:/container/path myapp:v1.0

# 只读挂载
docker run -v mydata:/app/data:ro myapp:v1.0
```

**Volume vs Bind Mount vs tmpfs**：

| 类型 | 说明 | 适用场景 |
|------|------|---------|
| Volume | Docker 管理的命名存储 | 生产环境数据持久化 |
| Bind Mount | 挂载宿主机路径 | 开发环境热重载 |
| tmpfs | 内存临时存储 | 敏感信息临时存储 |

### 2.4 网络（Network）

Docker 提供多种网络模式：

| 网络模式 | 说明 | 使用场景 |
|---------|------|---------|
| bridge | 默认模式，容器通过虚拟网桥通信 | 单机多容器通信 |
| host | 容器直接使用宿主机网络 | 需要高性能网络 |
| none | 无网络 | 安全隔离场景 |
| overlay | 跨主机容器通信 | Docker Swarm 集群 |
| macvlan | 容器拥有独立 MAC 地址 | 需要独立网络身份 |

```bash
# 创建自定义网络
docker network create --driver bridge mynet

# 运行容器并加入网络
docker run --network mynet --name app1 myapp:v1.0
docker run --network mynet --name app2 myapp:v1.0

# app1 可以通过容器名直接访问 app2
# ping app2 即可
```

---

## 三、Dockerfile 详解

### 3.1 Dockerfile 指令

```dockerfile
# 基础镜像
FROM golang:1.21-alpine

# 元数据标签
LABEL maintainer="developer@example.com"
LABEL version="1.0"
LABEL description="示例 Go 应用"

# 设置工作目录
WORKDIR /app

# 设置环境变量
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0

# 复制依赖文件（利用缓存机制）
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN go build -o main .

# 暴露端口（文档作用）
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/health || exit 1

# 启动命令
CMD ["./main"]
```

### 3.2 多阶段构建

多阶段构建可以显著减小最终镜像大小：

```dockerfile
# 第一阶段：构建
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o main .

# 第二阶段：运行（仅包含二进制文件）
FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### 3.3 Dockerfile 最佳实践

1. **利用构建缓存**：将不常变化的指令放在前面（如安装依赖），常变化的放在后面（如复制源代码）
2. **使用 .dockerignore**：排除不需要的文件（如 .git、node_modules、*.md）
3. **减少层数**：合并相关的 RUN 指令
4. **使用小基础镜像**：alpine 替代 ubuntu/debian
5. **非 root 用户运行**：添加 USER 指令提高安全性
6. **单个容器单个进程**：每个容器只运行一个主要进程

```dockerignore
# .dockerignore 示例
.git
.gitignore
*.md
docker-compose*.yml
.env
node_modules
__pycache__
*.pyc
.DS_Store
```

---

## 四、Docker Compose

### 4.1 基本配置

Docker Compose 用于定义和运行多容器应用：

```yaml
version: '3.8'

services:
  # Web 应用
  web:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_PORT=3306
      - REDIS_HOST=redis
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_started
    networks:
      - app-net
    restart: unless-stopped

  # MySQL 数据库
  db:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: rootpass
      MYSQL_DATABASE: myapp
      MYSQL_USER: appuser
      MYSQL_PASSWORD: apppass
    volumes:
      - db-data:/var/lib/mysql
    ports:
      - "3306:3306"
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app-net

  # Redis 缓存
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    networks:
      - app-net

volumes:
  db-data:
  redis-data:

networks:
  app-net:
    driver: bridge
```

### 4.2 常用命令

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f web

# 停止所有服务
docker-compose down

# 重建并启动
docker-compose up -d --build

# 进入某个服务容器
docker-compose exec web sh

# 扩展服务实例数
docker-compose up -d --scale web=3
```

---

## 五、Kubernetes 核心概念

### 5.1 架构概览

```
┌─────────────────────────────────────────────────────────┐
│                    Kubernetes 集群                       │
│                                                         │
│  ┌─────────────── Master Node ───────────────┐         │
│  │  API Server  Scheduler  Controller Manager│         │
│  │              etcd (存储)                   │         │
│  └───────────────────────────────────────────┘         │
│                                                         │
│  ┌──────── Worker Node 1 ────┐  ┌──── Worker Node 2 ──┐│
│  │  kubelet  kube-proxy      │  │  kubelet  kube-proxy ││
│  │  ┌─────┐ ┌─────┐ ┌─────┐ │  │  ┌─────┐ ┌─────┐    ││
│  │  │Pod 1│ │Pod 2│ │Pod 3│ │  │  │Pod 4│ │Pod 5│    ││
│  │  └─────┘ └─────┘ └─────┘ │  │  └─────┘ └─────┘    ││
│  │  Container Runtime        │  │  Container Runtime   ││
│  └───────────────────────────┘  └──────────────────────┘│
└─────────────────────────────────────────────────────────┘
```

### 5.2 Pod

Pod 是 Kubernetes 最小的调度单元，包含一个或多个紧密关联的容器。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  labels:
    app: myapp
    env: production
spec:
  containers:
    - name: myapp
      image: myapp:v1.0
      ports:
        - containerPort: 8080
      resources:
        requests:
          memory: "128Mi"
          cpu: "250m"
        limits:
          memory: "256Mi"
          cpu: "500m"
      livenessProbe:
        httpGet:
          path: /health
          port: 8080
        initialDelaySeconds: 10
        periodSeconds: 5
      readinessProbe:
        httpGet:
          path: /ready
          port: 8080
        initialDelaySeconds: 5
        periodSeconds: 3
    - name: sidecar
      image: log-agent:v1.0
```

### 5.3 Deployment

Deployment 提供声明式更新，支持滚动更新和回滚：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp-deployment
  labels:
    app: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1        # 滚动更新时最多多出1个Pod
      maxUnavailable: 0  # 更新期间不允许减少可用Pod
  template:
    metadata:
      labels:
        app: myapp
        version: v1.0
    spec:
      containers:
        - name: myapp
          image: myapp:v1.0
          ports:
            - containerPort: 8080
          env:
            - name: DB_HOST
              valueFrom:
                secretKeyRef:
                  name: db-secret
                  key: host
      imagePullSecrets:
        - name: registry-secret
```

```bash
# 常用操作
kubectl apply -f deployment.yaml
kubectl get deployments
kubectl get pods -l app=myapp
kubectl rollout status deployment/myapp-deployment
kubectl rollout history deployment/myapp-deployment
kubectl rollout undo deployment/myapp-deployment
kubectl scale deployment/myapp-deployment --replicas=5
```

### 5.4 Service

Service 为 Pod 提供稳定的网络访问入口：

```yaml
# ClusterIP（集群内部访问）
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  type: ClusterIP
  selector:
    app: myapp
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
---
# NodePort（外部通过节点端口访问）
apiVersion: v1
kind: Service
metadata:
  name: myapp-nodeport
spec:
  type: NodePort
  selector:
    app: myapp
  ports:
    - port: 80
      targetPort: 8080
      nodePort: 30080
---
# LoadBalancer（云环境负载均衡）
apiVersion: v1
kind: Service
metadata:
  name: myapp-lb
spec:
  type: LoadBalancer
  selector:
    app: myapp
  ports:
    - port: 80
      targetPort: 8080
```

### 5.5 Ingress

Ingress 管理外部 HTTP/HTTPS 访问：

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: myapp-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - myapp.example.com
      secretName: myapp-tls
  rules:
    - host: myapp.example.com
      http:
        paths:
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: api-service
                port:
                  number: 80
          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend-service
                port:
                  number: 80
```

### 5.6 ConfigMap 和 Secret

```yaml
# ConfigMap
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  APP_ENV: "production"
  LOG_LEVEL: "info"
  config.yaml: |
    server:
      port: 8080
      readTimeout: 30s
    database:
      maxOpenConns: 100
      maxIdleConns: 10
---
# Secret
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
type: Opaque
data:
  DB_PASSWORD: cGFzc3dvcmQxMjM=   # base64 编码
  API_KEY: YWJjZGVmZzEyMzQ1Ng==
```

### 5.7 HPA 自动扩缩容

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: myapp-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: myapp-deployment
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
        - type: Pods
          value: 2
          periodSeconds: 60
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
        - type: Percent
          value: 10
          periodSeconds: 60
```

---

## 六、Helm 包管理

### 6.1 Chart 结构

```
mychart/
├── Chart.yaml          # Chart 元信息
├── values.yaml         # 默认配置值
├── templates/          # 模板文件
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   ├── configmap.yaml
│   ├── _helpers.tpl    # 辅助模板
│   └── NOTES.txt       # 安装后提示
└── charts/             # 依赖的子 Chart
```

### 6.2 常用命令

```bash
# 创建 Chart
helm create mychart

# 打包 Chart
helm package mychart

# 安装
helm install myrelease mychart -f values-prod.yaml

# 升级
helm upgrade myrelease mychart --set image.tag=v2.0

# 回滚
helm rollback myrelease 1

# 列出所有 release
helm list

# 卸载
helm uninstall myrelease
```

---

## 七、生产环境最佳实践

### 7.1 镜像安全

1. 使用官方基础镜像，避免使用 latest 标签
2. 定期扫描镜像漏洞（Trivy、Snyk）
3. 最小化镜像体积（多阶段构建、alpine 基础镜像）
4. 不在镜像中存储敏感信息（使用 Secret）
5. 使用非 root 用户运行

### 7.2 容器安全

1. 设置资源限制（CPU、内存）
2. 使用只读文件系统（readOnlyRootFilesystem）
3. 禁用特权模式
4. 使用 PodSecurityPolicy 或 PodSecurityAdmission
5. 限制容器能力（drop ALL, add 最小必要能力）

### 7.3 监控与日志

```
监控栈推荐：
├── Prometheus     → 指标采集与存储
├── Grafana        → 可视化仪表板
├── Alertmanager   → 告警管理
├── Loki           → 日志聚合
├── Tempo          → 链路追踪
└── Jaeger         → 分布式追踪
```

### 7.4 CI/CD 流水线

```
代码提交 → 构建镜像 → 安全扫描 → 推送仓库 → 部署到 Staging → 集成测试 → 部署到 Production
```

---

## 八、常见排错命令

```bash
# 查看 Pod 事件
kubectl describe pod <pod-name>

# 查看 Pod 日志
kubectl logs <pod-name> -c <container-name> --previous

# 查看集群节点状态
kubectl get nodes -o wide

# 查看所有资源
kubectl get all -A

# 进入 Pod 调试
kubectl exec -it <pod-name> -- /bin/sh

# 查看资源使用情况
kubectl top pods
kubectl top nodes

# 干运行（不实际执行）
kubectl apply -f deployment.yaml --dry-run=client

# 强制删除卡住的 Pod
kubectl delete pod <pod-name> --force --grace-period=0
```

---

## 九、Docker vs Podman 对比

| 特性 | Docker | Podman |
|------|--------|--------|
| 守护进程 | 需要 dockerd | 无守护进程 |
| 权限 | 需要 root（默认） | 支持 rootless |
| 兼容性 | 行业标准 | 兼容 Docker CLI |
| 安全性 | 较低 | 较高（无守护进程攻击面） |
| Pod 支持 | 不支持 | 原生支持 Pod |
| Systemd 集成 | 需要额外配置 | 原生支持 |
| Kubernetes YAML | 需要生成 | 支持 podman generate kube |

---

## 十、常用 Docker 命令速查

```bash
# 系统管理
docker system df          # 查看磁盘使用
docker system prune -a    # 清理所有未使用资源
docker system info        # 查看系统信息

# 镜像管理
docker images -a          # 列出所有镜像
docker image prune -a     # 清理未使用镜像
docker tag <id> myapp:v1  # 标记镜像

# 容器管理
docker container ls -a    # 列出所有容器
docker container prune    # 清理已停止容器
docker inspect <container> # 查看容器详情
docker stats              # 实时资源使用统计

# 网络管理
docker network ls         # 列出网络
docker network inspect bridge  # 查看网络详情

# 数据卷管理
docker volume ls          # 列出数据卷
docker volume inspect mydata  # 查看数据卷详情
```
