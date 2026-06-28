# Deployment

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
sudo ufw allow 22/tcp
sudo ufw allow 6443/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### Создать namespace

```bash
sudo kubectl create namespace dev
```

## 2. GitLab Registry

### Создать Deploy Token

GitLab → Settings → Repository → **Deploy tokens**

- Name: `k8s-registry`
- Scope: `read_registry`
- Нажать **Create deploy token**, сохранить `username` и `token`

### Добавить секрет в кластер

```bash
sudo kubectl create secret docker-registry gitlab-registry \
  --docker-server=registry.gitlab.com \
  --docker-username=username_из_токена \
  --docker-password=token_из_токена \
  --namespace=dev
```

## 3. Переменные окружения

Backend и PostgreSQL читают конфиг из общего Secret `cu-faculty-backend-env`
(единый источник правды). Postgres использует из него `POSTGRES_DB/USER/PASSWORD`,
остальные ключи игнорирует.

1. Скопировать `.env` на сервер (например через scp) и поправить значения под прод - как минимум `ALLOWED_ORIGINS` (публичный URL вместо `localhost:3000`).

2. Создать Secret из файла:

```bash
sudo kubectl create secret generic cu-faculty-backend-env \
  --from-env-file=k8s/.env -n dev
```

## 4. Деплой

### PostgreSQL

```bash
sudo kubectl apply -f k8s/postgres.yaml
```

### Garage (S3)

В `k8s/garage.yaml` поля `rpc_secret` и `admin_token` пустые - их нужно
сгенерировать и подставить при применении (в git они не хранятся):

```bash
RPC_SECRET=$(openssl rand -hex 32)
ADMIN_TOKEN=$(openssl rand -base64 32)

sed -e "s|rpc_secret = \"\"|rpc_secret = \"$RPC_SECRET\"|" \
    -e "s|admin_token = \"\"|admin_token = \"$ADMIN_TOKEN\"|" \
    k8s/garage.yaml | sudo kubectl apply -f -
```

### Backend

```bash
sudo kubectl apply -f k8s/dev.yaml
```

## 5. Инициализация Garage (выполнить один раз)

После первого запуска нужно собрать layout кластера и создать bucket с ключом
доступа. Внутри кластера Garage доступен по адресу `http://garage:3900`.

```bash
# 1. Узнать node ID
sudo kubectl exec -n dev garage-0 -- /garage status

# 2. Назначить layout (взять префикс node ID из вывода выше)
sudo kubectl exec -n dev garage-0 -- /garage layout assign -z dc1 -c 10G NODE_ID
sudo kubectl exec -n dev garage-0 -- /garage layout apply --version 1

# 3. Создать bucket
sudo kubectl exec -n dev garage-0 -- /garage bucket create cu-faculty

# 4. Создать ключ доступа (запомнить Key ID и Secret из вывода)
sudo kubectl exec -n dev garage-0 -- /garage key create cu-faculty-app

# 5. Выдать ключу права на bucket
sudo kubectl exec -n dev garage-0 -- \
  /garage bucket allow --read --write cu-faculty --key cu-faculty-app
```

Полученные креды прописать в `.env` и пересоздать Secret (см. раздел 7):

```plaintext
S3_ENDPOINT="http://garage:3900"
S3_REGION="garage"
S3_BUCKET="cu-faculty"
S3_ACCESS_KEY="<Key ID>"
S3_SECRET_KEY="<Secret>"
```

## 6. Проверка

```bash
sudo kubectl get pods -n dev
sudo kubectl get pvc -n dev
```

Все поды должны быть `Running`/`Ready`. Приложение доступно по
`http://<IP_СЕРВЕРА>`.

---

## 7. Обновление конфигурации

### Изменить env (.env)

```bash
sudo kubectl create secret generic cu-faculty-backend-env \
  --from-env-file=k8s/.env -n dev \
  --dry-run=client -o yaml | sudo kubectl apply -f -

sudo kubectl rollout restart deployment/cu-faculty-backend -n dev
sudo kubectl rollout restart statefulset/postgres -n dev
```

### Обновить образ backend

```bash
sudo kubectl rollout restart deployment/cu-faculty-backend -n dev
```

---

## 8. GitLab CI/CD

### Получить KUBE_CONFIG

На сервере выполнить (заменить IP на реальный):

```bash
sudo cat /etc/rancher/k3s/k3s.yaml \
  | sed 's/127.0.0.1/IP_СЕРВЕРА/g' \
  | base64 -w0
```

### Добавить переменную в GitLab

GitLab → Settings → CI/CD → Variables → **Add variable**:

| Key            | Value                | Type     |
|----------------|----------------------|----------|
| `KUBE_CONFIG`  | вывод команды выше   | Variable |
