# go-market
![Go](https://img.shields.io/badge/-Go-00ADD8?logo=go)
![REST](https://img.shields.io/badge/-REST-FF6C37?logo=rest&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-85EA2D?logo=swagger&logoColor=black)
![PostgreSQL](https://img.shields.io/badge/-PostgreSQL-4169E1?logo=postgresql)

## Как запустить контейнер
Сборка бинарных файлов сервера:
```shell
task build-server
```

Запустите локально Docker:
```shell
docker-compose up -d
```

## Конфигурация
### Конфигурация Сервера
Сервер поддерживает настройку через переменные окружения, аргументы командной строки:

| Переменная | Флаг | По умолчанию | Описание                                             |
|------------|------|--------------|------------------------------------------------------|
| `RUN_ADDRESS` | `-a` | `localhost:8080` | Адрес HTTP сервера                                   |
| `ACCRUAL_SYSTEM_ADDRESS` | `-r` | `localhost:8081` | Адрес HTTP сервера системы расчета баллов лояльности |
| `LOG_LEVEL` | `-l` | `INFO` | Уровень логирования (DEBUG, INFO, WARN, ERROR)       |
| `DATABASE_URI` | `-d` | - | URI для подключения к PostgreSQL                     |
| `AUTH_KEY` | -    | hex-ключ | Ключ для аутентификации JWT                          |
| `TERMINATION_TIMEOUT` | -    | `30s` | Таймаут graceful shutdown                            |
| `WORKER_COUNT` | -    | `10` | Количество воркеров                                  |