# Deployment Guide: Go Backend on Kubernetes (k3s) via GitLab CI

## Требования

- VDS с Ubuntu 20.04/22.04/24.04 (минимум 2 CPU, 2 GB RAM)
- GitLab репозиторий с Go проектом
- Доступ по SSH к серверу

---

## 1. Настройка сервера

### Установка k3s

```bash
curl -sfL https://get.k3s.io | sh -
```

Проверить что кластер работает:

```bash
sudo kubectl get nodes
```

### Открыть порты (UFW)

```bash
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 6443/tcp  # Kubernetes API
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
```

> Если у провайдера есть отдельный файрвол (security group) — открой те же порты и там.

### Создать namespace'ы

```bash
sudo kubectl create namespace staging
sudo kubectl create namespace production
```

---

## 2. Настройка GitLab Registry

### Создать Deploy Token

GitLab → Settings → Repository → **Deploy tokens**

- Name: `k8s-registry`
- Scope: `read_registry`
- Нажать **Create deploy token**, сохранить `username` и `token`

### Добавить секрет в кластер

```bash
sudo kubectl create secret docker-registry gitlab-registry \
  --docker-server=registry.gitlab.com \
  --docker-username=<username из токена> \
  --docker-password=<token из токена> \
  --namespace=staging
```

```bash
sudo kubectl create secret docker-registry gitlab-registry \
  --docker-server=registry.gitlab.com \
  --docker-username=<username из токена> \
  --docker-password=<token из токена> \
  --namespace=production
```

---

## 3. Деплой приложения на кластер

### Применить манифесты

Создать файл `k8s/staging.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cu-faculty-backend
  namespace: staging
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cu-faculty-backend
  template:
    metadata:
      labels:
        app: cu-faculty-backend
    spec:
      imagePullSecrets:
        - name: gitlab-registry
      containers:
        - name: cu-faculty-backend
          image: registry.gitlab.com/<your-project>/cu-faculty-backend:latest
          ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: cu-faculty-backend
  namespace: staging
spec:
  selector:
    app: cu-faculty-backend
  ports:
    - port: 80
      targetPort: 8080
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: cu-faculty-backend
  namespace: staging
spec:
  ingressClassName: traefik
  rules:
    - http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: cu-faculty-backend
                port:
                  number: 80
```

```bash
sudo kubectl apply -f k8s/staging.yaml
```

> Для production — создай аналогичный `k8s/production.yaml` с `namespace: production`.

---

## 4. Настройка GitLab CI/CD

### Получить KUBE_CONFIG

На сервере выполнить (заменить IP на реальный):

```bash
sudo cat /etc/rancher/k3s/k3s.yaml \
  | sed 's/127.0.0.1/<IP_СЕРВЕРА>/g' \
  | base64 -w0
```

### Добавить переменную в GitLab

GitLab → Settings → CI/CD → Variables → **Add variable**:

| Key            | Value                | Type     |
|----------------|----------------------|----------|
| `KUBE_CONFIG`  | вывод команды выше   | Variable |

> Убедись что значение — одна строка без переносов.

---

## 5. Проверка

После пуша в `main`:

1. **test** — запускает `go test ./...`
2. **build** — собирает Docker образ и пушит в GitLab Registry
3. **deploy:staging** — автоматически обновляет deployment в namespace `staging`
4. **deploy:production** — запускается вручную через GitLab UI

Проверить статус подов:

```bash
sudo kubectl get pods -n staging
sudo kubectl get pods -n production
```

Приложение доступно по `http://<IP_СЕРВЕРА>` (staging через Traefik на порту 80).

---

## Возможные проблемы

### `no users found` при старте пода

Приложение запущено с несуществующим пользователем. В `Dockerfile` нужно создать пользователя перед использованием:

```dockerfile
FROM alpine:3.23
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/main main
EXPOSE 8080
USER appuser:appgroup
ENTRYPOINT ["./main"]
```

### `base64: invalid input` в CI

Переменная `KUBE_CONFIG` задана неверно. Убедись что:

- Значение — вывод `base64 -w0`, а не сам kubeconfig
- Тип переменной — **Variable**, не File
- Нет лишних пробелов

### Приложение недоступно снаружи (Bad Gateway / Connection refused)

Go сервер должен слушать на всех интерфейсах, не только на localhost:

```go
// Неправильно
Addr: "localhost:8080"

// Правильно
Addr: ":8080"
```
