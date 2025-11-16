# avitoTest


- Доп. эндпоинт статистики: `/stats/assignments` (текущее количество назначений по пользователям и по PR)


## для запуска 

### В корне проекта
cp .env.example .env 
docker compose up --build

### локальный запуск 
1) поднять бд:
docker compose up -d db
2) поднять миграции:
make migrate-up
3) запуск:
make run


Архитектура(слоев)

internal/domain — сущности домена и их типы

internal/ports — интерфейсы (контракты) сервисов

internal/repository — доступ к данным (Postgres, pgx)

internal/service — бизнес-логика (транзакции через TxManager)

internal/http — хендлеры/роутер/DTO/мидлвары (Gin)

internal/app — wiring (инициализация, HTTP-сервер)

migrations — SQL-миграции