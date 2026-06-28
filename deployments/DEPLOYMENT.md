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
  --docker-username=username_из_токена \
  --docker-password=token_из_токена \
  --namespace=dev
```

## 3. Деплой приложения на кластер

### Применить манифесты

```bash
sudo kubectl apply -f cu-niti/dev.yaml
```

## 4. Настройка GitLab CI/CD

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

> Убедись что значение — одна строка без переносов.

## 5. Проверка

Проверить статус подов:

```bash
sudo kubectl get pods -n dev
```

Приложение доступно по `http://<IP_СЕРВЕРА>`
