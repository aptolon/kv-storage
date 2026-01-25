# Key-Value Storage Server

Key-Value Storage Server на Go с in-memory хранилищем и сохранением состояния в PostgreSQL.

Сервер работает по TCP-протоколу, поддерживает конкурентных клиентов, корректное завершение работы и восстановление данных после перезапуска.

---

## Возможности

- TCP сервер с текстовым протоколом
- In-memory Key-Value хранилище
- Потокобезопасный доступ
- Изоляция данных
- Сохранение состояния в PostgreSQL
- Восстановление данных при старте
- Graceful shutdown
- Docker и docker-compose
- Unit, integration и stress тесты

---

### Компоненты

- **storage**
  - интерфейс `Storage`
  - реализация `MemoryStorage`
  - потокобезопасность
  - изоляция данных
  - поддержка snapshot

- **server**
  - TCP сервер
  - обработка нескольких клиентов
  - отдельная goroutine на соединение
  - остановка через `context.Context`

- **persistence**
  - snapshot всего состояния (`map[string][]byte`)
  - сохранение и загрузка из PostgreSQL

---

## Архитектура 

Client -> TCP Server -> In-memory Storage -> PostgreSQL (snapshot)

---

## TCP протокол

Сервер принимает текстовые команды, по одной на строку.

### Команды

SET key value   -> OK
GET key         -> VALUE value | NULL
DEL key         -> OK

---

## Persistence

- при запуске сервера состояние загружается из PostgreSQL
- при завершении работы текущее состояние сохраняется целиком
- данные хранятся в таблице с ключом и значением (`BYTEA`)

---

## Запуск

Сервер можно запускать локально и в Docker контейнере

### Локально 

```bash
make build 
make run
```

### В контейнере

```bash
make docker-build
make docker-up
```

---

## Тестирование

```bash
make test / make test-run
```

Покрываются:
- операции хранилища
- конкурентный доступ
- TCP команды
- graceful shutdown
- сохранение и восстановление данных
- нагрузочные сценарии

---

## Возможные улучшения
- бинарный протокол
- HTTP API