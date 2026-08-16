# Документация Amai API

Базовый URL: `http://<host>:<port>/`

Все ответы возвращаются в формате JSON. Ответы с ошибками обычно имеют следующую структуру:

```json
{ "message": "понятное человеку описание ошибки" }

```

---

## Аутентификация

Аутентификация осуществляется по двухэтапной схеме:

1. **`POST /auth/login`** — аутентификация с использованием **HTTP Basic Auth** и одной хардкодной учетной записи администратора (данные берутся из переменных окружения `admin_login` / `admin_password`; система пользователей отсутствует). В случае успеха на стороне сервера в карту памяти сохраняется непрозрачный случайно сгенерированный **ID сессии (session ID)**, который возвращается клиенту в виде `HttpOnly`, `Secure` куки с именем `sessionId`.
2. **Последующие запросы администратора** — аутентифицируются с помощью куки `sessionId`, которая валидируется через хранилище сессий в памяти. Предполагается, что именно это поведение обеспечивает `authMiddleware()` (которая вызывает `auth.CheckCookieAuth`) для группы маршрутов `/admin/`.

Время жизни сессии (TTL) составляет **24 часа** и обновляется ("скользящая пролонгация") при каждом успешном вызове `CheckCookieAuth` — срок действия куки и время ее истечения на стороне сервера продлеваются при каждом аутентифицированном запросе.

---

## Публичные эндпоинты

### `GET /health`

Эндпоинт проверки работоспособности (Health check).

**Ответ — 200 OK**

```json
{ "message": "alive" }

```

---

### `GET /post/get`

Получение одного поста по его ID.

**Query-параметры**

| Имя | Тип | Обязательный | Описание |
| --- | --- | --- | --- |
| `id` | UUID | Да | Идентификатор поста |

**Ответы**

* `200 OK`
```json
{
  "id": "e2e5a1b0-...",
  "title": "Заголовок поста",
  "poster_id": "author-id",
  "created_at": "2026-01-01T12:00:00Z",
  "updated_at": "2026-01-02T09:00:00Z",
  "body": "Содержимое поста..."
}

```


Поле `updated_at` имеет значение `null`, если пост никогда не редактировался (время создания == времени обновления).
* `400 Bad Request`
* `{ "message": "query id not set" }`
* `{ "message": "invalid UUID format" }`


---

### `GET /post/all`

Получение пагинированного списка постов.

**Query-параметры**

| Имя | Тип | Обязательный | Описание |
| --- | --- | --- | --- |
| `page` | int | Да | Номер страницы, должен быть `> 0` |

**Ответы**

* `200 OK`
```json
{
  "posts": [ /* массив объектов Post, аналогично GET /post/get */ ],
  "pages": 5
}

```

* `400 Bad Request`
* `{ "message": "query page not set" }`
* `{ "message": "query page should be int" }`
* `{ "message": "query page must be greater than 0" }`
* `{ "message": "page not found" }` — запрашиваемая страница превышает общее количество страниц

---

### `GET /file/get`

Скачивание / потоковая передача файла по ID.

**Query-параметры**

| Имя | Тип | Обязательный | Описание |
| --- | --- | --- | --- |
| `id` | UUID | Да | Идентификатор файла |

**Ответы**

* `200 OK` — сырое содержимое файла (отдается через `c.File`), `Content-Type` определяется из файла
* `400 Bad Request`
* `{ "message": "query id not set" }`
* `{ "message": "invalid UUID format" }`



---

### `POST /auth/login`

Аутентификация в качестве администратора через **HTTP Basic Auth**.

**Запрос**

Тело JSON отсутствует. Учетные данные отправляются через стандартный заголовок `Authorization`:

```
Authorization: Basic base64(admin_login:admin_password)

```

Единственная валидная пара логин/пароль настраивается на стороне сервера через переменные окружения `admin_login` и `admin_password`.

**Ответы**

* `202 Accepted`
```json
{ "message": "login" }

```


В случае успеха сервер устанавливает куку сессии:
```
Set-Cookie: sessionId=<random-256-bit-token>; Path=/; Max-Age=86400; HttpOnly; Secure

```

---

### `POST /auth/logout`

Инвалидация текущей сессии.

**Запрос**

Тело отсутствует. Требуется кука `sessionId`, установленная при входе (отправляется браузером автоматически).

**Ответы**

* `202 Accepted`
```json
{ "message": "logout" }

```

---

## Эндпоинты администратора

Все приведенные ниже маршруты находятся по пути `/admin/` и защищены `authMiddleware()`.

### `POST /admin/post/create`

Создание нового поста.

**Тело запроса**

```json
{
  "title": "Заголовок поста",
  "poster_id": "author-id",
  "body": "Содержимое поста..."
}

```

*(`id`, `created_at`, `updated_at` игнорируются / не требуются при создании.)*

**Ответы**

* `201 Created`
```json
{ "message": "post created" }

```

* `400 Bad Request` — `{ "message": "invalid JSON body" }`

---

### `PUT /admin/post/edit`

Редактирование существующего поста.

**Тело запроса**

```json
{
  "id": "e2e5a1b0-...",
  "title": "Обновленный заголовок",
  "poster_id": "author-id",
  "body": "Обновленное содержимое..."
}

```

**Ответы**

* `200 OK` — `{ "message": "post edited" }`
* `400 Bad Request`
* `{ "message": "invalid JSON body" }`
* `{ "message": "invalid UUID" }` — `id` отсутствует или равен нулевому UUID

---

### `DELETE /admin/post/delete`

Удаление поста по ID.

**Query-параметры**

| Имя | Тип | Обязательный | Описание |
| --- | --- | --- | --- |
| `id` | UUID | Да | Идентификатор поста |

**Ответы**

* `200 OK` — `{ "message": "post deleted" }`
* `400 Bad Request`
* `{ "message": "query id not set" }`
* `{ "message": "invalid UUID format" }`

---

### `GET /admin/file/list`

Список загруженных файлов с пагинацией.

**Query-параметры**

| Имя | Тип | Обязательный | Описание |
| --- | --- | --- | --- |
| `page` | int | Да | Номер страницы, должен быть `> 0` |

**Ответы**

* `200 OK`
```json
{
  "files": [
    { "id": "f1a2...", "ext": ".png" }
  ],
  "pages": 3
}

```

* `400 Bad Request`
* `{ "message": "query page not set" }`
* `{ "message": "query page should be int" }`
* `{ "message": "query page must be greater than 0" }`
* `{ "message": "page not found" }`

---

### `POST /admin/file/upload`

Загрузка файла.

**Запрос:** `multipart/form-data` с единственным полем под названием `file`.

**Разрешенные MIME-типы**:

* `image/jpeg`, `image/png`, `image/webp`, `image/gif`
* `text/markdown`, `text/plain`
* `application/pdf`
* `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet` (.xlsx)
* `application/vnd.openxmlformats-officedocument.wordprocessingml.document` (.docx)
* `application/vnd.openxmlformats-officedocument.presentationml.presentation` (.pptx)

**Ответы**

* `200 OK`
```json
{ "message": "file uploaded. File id: <uuid>" }

```

* `400 Bad Request` — `{ "message": "no file uploaded" }` / `{ "message": "empty file uploaded" }`
* `415 Unsupported Media Type` — `{ "message": "invalid file uploaded" }`
* `500 Internal Server Error` — `{ "message": "failed to process file" }`


---

### `DELETE /admin/file/delete`

Удаление файла по ID.

**Query-параметры**

| Имя | Тип | Обязательный | Описание |
| --- | --- | --- | --- |
| `id` | UUID | Да | Идентификатор файла |

**Ответы**

* `200 OK` — `{ "message": "file deleted" }`
* `400 Bad Request`
* `{ "message": "query id not set" }`
* `{ "message": "invalid UUID format" }`

---

## Сводная таблица

| Метод | Путь | С аутентификацией | Описание |
| --- | --- | --- | --- |
| GET | `/health` | Нет | Проверка работоспособности |
| GET | `/post/get` | Нет | Получить пост по ID |
| GET | `/post/all` | Нет | Список постов (пагинация) |
| GET | `/file/get` | Нет | Скачать файл по ID |
| POST | `/auth/login` | Нет (Basic Auth) | Вход, выдает куку `sessionId` |
| POST | `/auth/logout` | Требует куку `sessionId` | Выход, инвалидирует сессию |
| POST | `/admin/post/create` | Да | Создать пост |
| PUT | `/admin/post/edit` | Да | Редактировать пост |
| DELETE | `/admin/post/delete` | Да | Удалить пост |
| GET | `/admin/file/list` | Да | Список файлов (пагинация) |
| POST | `/admin/file/upload` | Да | Загрузить файл |
| DELETE | `/admin/file/delete` | Да | Удалить файл |