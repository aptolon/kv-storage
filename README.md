# Key-Value хранилище

## Контракт `Storage`

`Storage` — абстракция для работы с Key-Value данными.

### `Set(key string, value []byte) error`

Записывает значение по ключу.

- если ключ существует — значение перезаписывается
- возвращает `error` только в случае внутренней ошибки

---

### `Get(key string) (value []byte, err error)`

Получает значение по ключу.

- если ключ существует — возвращает `(value, nil)`
- если ключ не существует — возвращает `(nil, nil)`
- возвращает `error` только в случае внутренней ошибки

---

### `Delete(key string) error`

Удаляет значение по ключу.

- если ключ не существует — ничего не делает
- возвращает `error` только в случае внутренней ошибки

---

## `MemoryStorage`

`MemoryStorage` — быстрое in-memory KV-хранилище, реализующее контракт `Storage`.

### `MemoryStorage` struct

Структура содержит:
- `sync.RWMutex` для потокобезопасного доступа
- внутреннее хранилище на основе `map[string][]byte`

---

### `NewMemoryStorage`

Конструктор, возвращающий инициализированное in-memory хранилище.

---

### Методы `Set`, `Get`, `Delete`

Методы реализуют контракт `Storage` и являются:
- потокобезопасными
- изолирующими данные (копирование `[]byte`)
- корректными при конкурентном доступе

---

## Тестирование `Storage`

Тесты проверяют контракт `Storage` и применимы к любой его реализации.

### `TestStorageSetGet`
Проверка корректной записи и получения значения.

### `TestStorageSetOverride`
Проверка корректной перезаписи значения по существующему ключу.

### `TestStorageGetMissingKey`
Проверка возврата `(nil, nil)` при отсутствии ключа.

### `TestStorageDeleteMissingKey`
Проверка корректной работы `Delete` при отсутствии ключа.

### `TestStorageDeleteKey`
Проверка корректного удаления существующего ключа.

### `TestStorageIsolatedSet`
Проверка изоляции данных при записи.

### `TestStorageIsolatedGet`
Проверка изоляции данных при получении.

### `TestStorageConcurrentSetNoRace`
Проверка конкурентной записи без data race.

### `TestStorageConcurrentGetNoRace`
Проверка конкурентного чтения без data race.

### `TestStorageConcurrentSetGetNoRace`
Проверка параллельной записи и чтения без data race.
