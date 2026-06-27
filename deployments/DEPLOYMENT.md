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
sudo kubectl create namespace dev
sudo kubectl create namespace prod
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
  --namespace=dev
```

```bash
sudo kubectl create secret docker-registry gitlab-registry \
  --docker-server=registry.gitlab.com \
  --docker-username=<username из токена> \
  --docker-password=<token из токена> \
  --namespace=prod
```

---

## 3. Деплой приложения на кластер

### Применить манифесты

Манифесты уже лежат в репозитории: `deployments/k8s/dev.yaml` и `deployments/k8s/prod.yaml`. Каждый содержит `Deployment` + `Service` + `Ingress` (Traefik) с health-проба­ми на `/health` и лимитами ресурсов.

Образ указан как `registry.gitlab.com/cu-faculty/backend:latest` (совпадает с `$CI_REGISTRY_IMAGE`). Для другого проекта замени путь на свой `registry.gitlab.com/<namespace>/<project>`.

```bash
sudo kubectl apply -f deployments/k8s/dev.yaml
sudo kubectl apply -f deployments/k8s/prod.yaml
```

> Манифесты описывают только Go backend. Фронтенд (nginx со статикой SvelteKit, проксирует `/api` на backend) деплоится отдельно и в этом гайде не рассматривается.

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

После пуша в ветку (или открытия MR) пайплайн проходит стадии:

1. **test** — `cd backend && go test` с coverage-отчётом (cobertura)
2. **build** — собирает Docker образ backend и пушит в GitLab Registry (теги `$CI_COMMIT_SHORT_SHA` и `latest`)
3. **deploy:dev** — обновляет deployment в namespace `dev`, **запускается вручную** через GitLab UI
4. **deploy:prod** — обновляет deployment в namespace `prod`, **запускается вручную** через GitLab UI

Проверить статус подов:

```bash
sudo kubectl get pods -n dev
sudo kubectl get pods -n prod
```

Приложение доступно по `http://<IP_СЕРВЕРА>` (dev через Traefik на порту 80).

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
