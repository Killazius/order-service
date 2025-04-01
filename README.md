# OrderService gRPC

Этот проект представляет собой gRPC-сервис для управления заказами. Сервис поддерживает как gRPC, так и HTTP-запросы через grpc-gateway.
Написан по ходу специализации от Яндекс Лицея на курсе [Веб-разработка на Go](https://lyceum.yandex.ru/web-go)

## Конфигурация

### Настройка окружения

1. Создайте файл `.env` на основе примера `.env.example`:
```env
CONFIG_PATH="config/<name>.yaml"

POSTGRES_HOST: localhost
POSTGRES_PORT: 5432
POSTGRES_USER: postgres
POSTGRES_PASSWORD: postgres
POSTGRES_DB: postgres
POSTGRES_MAX_CONN: 10
POSTGRES_MIN_CONN: 5
```
   Где `<name>` — название вашего конфигурационного файла. (по умолчанию config/config.yaml)

2. Создайте конфигурационный файл по пути, указанному в `.env`. Пример содержимого конфигурационного файла:
```yaml
GRPC_PORT: 50051  # Порт для gRPC-сервера (по умолчанию 50051)
HTTP_PORT: 8080  # Порт для HTTP-сервера (по умолчанию 8080)
```

## Использование Makefile
В проекте предоставлен `Makefile` для упрощения сборки и запуска проекта. Доступные команды:
- `make build` — сборка бинарного файла.
- `make exec` — запуск бинарного файла.
- `make run` — сборка и запуск бинарного файла.
- `make proto` — генерация кода для gRPC и grpc-gateway из proto-файла.
- `make docker` — запуск контейнеров Docker.

## Запуск проекта

1. Убедитесь, что у вас установлены все зависимости и инструменты для работы с gRPC и protobuf. 
2. Выполните команду `make docker` для запуска контейнеров Docker.