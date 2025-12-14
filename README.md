# Anti-Bruteforce Service

Сервис предназначен для борьбы с подбором паролей при авторизации в какой-либо системе.
Сервис вызывается перед авторизацией пользователя и может либо разрешить, либо заблокировать попытку.
Предполагается, что сервис используется только для server-server, т.е. скрыт от конечного пользователя.

Ограничивает количество попыток аутентификации по:
- логину
- паролю
- IP-адресу

Поддерживает whitelist и blacklist IP-адресов.

---

## Стек

- Go
- PostgreSQL
- Docker
- Docker Compose

---

## Сборка и запуск

Проект управляется через `Makefile`.

---
### Сборка бинарника

```bash
make build
```

В результате будет собран бинарный файл:

```bash
bin/antibruteforce
```

---
### Запуск сервиса

```bash
make run
```

Команда запускает сервис и базу данных через `docker compose up`.

Для остановки:

```bash
make down
```

---

## Тестирование

```bash
make test
```

Запускаются все unit и integration тесты.

---

## Конфигурация

Сервис настраивается через переменные окружения.

Пример `.env`:

```env
PORT=8080
LIMIT_LOGIN=10
LIMIT_PASSWORD=100
LIMIT_IP=1000

# Для контейнера
POSTGRES_HOST=antibruteforce_db
POSTGRES_PORT=5432
POSTGRES_USER=antibruteforce
POSTGRES_PASSWORD=secret
POSTGRES_DB=antibruteforce

# Для локального теста
LOCAL_POSTGRES_HOST=127.0.0.1
LOCAL_POSTGRES_PORT=5432
```
---

## CLI

### whitelist / blacklist
  `antibruteforce-cli whitelist add <cidr>`
  `antibruteforce-cli whitelist remove <cidr>`
  `antibruteforce-cli blacklist add <cidr>`
  `antibruteforce-cli blacklist remove <cidr>`

### backet
  `antibruteforce-cli bucket reset`

---

## API

### whitelist / blacklist

#### Добавление IP в whitelist / blacklist
`/whitelist/add"`
`/blacklist/add`

#### Удаление IP из whitelist / blacklist
`/whitelist/remove`
`/blacklist/remove`


### backet
`/check`

#### Сброс bucket
`/bucket/reset`
---

## Логика ограничений

* лимит по логину
* лимит по паролю
* лимит по IP-адресу

Whitelist всегда разрешает запросы, blacklist — всегда блокирует.
