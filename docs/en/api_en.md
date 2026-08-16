# Amai API Documentation

Base URL: `http://<host>:<port>/`

All responses are JSON. Error responses generally follow the shape:

```json
{ "message": "human readable error description" }
```

---

## Authentication

Authentication is a two-step scheme:

1. **`POST /auth/login`** — authenticate via **HTTP Basic Auth**, using a
   single hardcoded admin account (credentials from the `admin_login` /
   `admin_password` environment variables — there is no per-user account
   system). On success, an opaque, randomly generated **session ID** is
   stored server-side in an in-memory map and returned to the client as an
   `HttpOnly`, `Secure` cookie named `sessionId`.
2. **Subsequent admin requests** — authenticated via that `sessionId` cookie,
   validated against the in-memory session store. This is presumably what
   `authMiddleware()` (not shown, but implied to call `auth.CheckCookieAuth`)
   enforces on the `/admin/` route group.

Session TTL is **24 hours**, refreshed ("sliding expiration") on every
successful `CheckCookieAuth` call — the cookie and the server-side expiry are
both extended on each authenticated request.

---

## Public Endpoints

### `GET /health`

Health check endpoint.

**Response — 200 OK**
```json
{ "message": "alive" }
```

---

### `GET /post/get`

Fetch a single post by ID.

**Query Parameters**

| Name | Type | Required | Description        |
|------|------|----------|---------------------|
| `id` | UUID | Yes      | Post identifier      |

**Responses**

- `200 OK`
  ```json
  {
    "id": "e2e5a1b0-...",
    "title": "Post title",
    "poster_id": "author-id",
    "created_at": "2026-01-01T12:00:00Z",
    "updated_at": "2026-01-02T09:00:00Z",
    "body": "Post content..."
  }
  ```
  `updated_at` is `null` if the post has never been edited (created == updated timestamp).

- `400 Bad Request`
  - `{ "message": "query id not set" }`
  - `{ "message": "invalid UUID format" }`

---

### `GET /post/all`

Fetch a paginated list of posts.

**Query Parameters**

| Name   | Type | Required | Description                  |
|--------|------|----------|-------------------------------|
| `page` | int  | Yes      | Page number, must be `> 0`    |

**Responses**

- `200 OK`
  ```json
  {
    "posts": [ /* array of Post objects, same shape as GET /post/get */ ],
    "pages": 5
  }
  ```

- `400 Bad Request`
  - `{ "message": "query page not set" }`
  - `{ "message": "query page should be int" }`
  - `{ "message": "query page must be greater than 0" }`
  - `{ "message": "page not found" }` — requested page exceeds total pages

---

### `GET /file/get`

Download/stream a file by ID.

**Query Parameters**

| Name | Type | Required | Description        |
|------|------|----------|---------------------|
| `id` | UUID | Yes      | File identifier      |

**Responses**

- `200 OK` — raw file contents (served via `c.File`), `Content-Type` inferred from file
- `400 Bad Request`
  - `{ "message": "query id not set" }`
  - `{ "message": "invalid UUID format" }`

---

### `POST /auth/login`

Authenticate as the admin user via **HTTP Basic Auth**.

**Request**

No JSON body. Credentials are sent via the standard `Authorization` header:

```
Authorization: Basic base64(admin_login:admin_password)
```

The single valid username/password pair is configured server-side via the
`admin_login` and `admin_password` environment variables.

**Responses**

- `202 Accepted`
  ```json
  { "message": "login" }
  ```
  On success, the server sets a session cookie:
  ```
  Set-Cookie: sessionId=<random-256-bit-token>; Path=/; Max-Age=86400; HttpOnly; Secure
  ```

---

### `POST /auth/logout`

Invalidate the current session.

**Request**

No body. Requires the `sessionId` cookie set at login (sent automatically by
the browser).

**Responses**

- `202 Accepted`
  ```json
  { "message": "logout" }
  ```

---

## Admin Endpoints

All routes below are mounted under `/admin/` and protected by `authMiddleware()`.

### `POST /admin/post/create`

Create a new post.

**Request Body**
```json
{
  "title": "Post title",
  "poster_id": "author-id",
  "body": "Post content..."
}
```

**Responses**

- `201 Created`
  ```json
  { "message": "post created" }
  ```

- `400 Bad Request` — `{ "message": "invalid JSON body" }`

---

### `PUT /admin/post/edit`

Edit an existing post.

**Request Body**
```json
{
  "id": "e2e5a1b0-...",
  "title": "Updated title",
  "poster_id": "author-id",
  "body": "Updated content..."
}
```

**Responses**

- `200 OK` — `{ "message": "post edited" }`
- `400 Bad Request`
  - `{ "message": "invalid JSON body" }`
  - `{ "message": "invalid UUID" }` — `id` missing/zero UUID

---

### `DELETE /admin/post/delete`

Delete a post by ID.

**Query Parameters**

| Name | Type | Required | Description        |
|------|------|----------|---------------------|
| `id` | UUID | Yes      | Post identifier      |

**Responses**

- `200 OK` — `{ "message": "post deleted" }`
- `400 Bad Request`
  - `{ "message": "query id not set" }`
  - `{ "message": "invalid UUID format" }`

---

### `GET /admin/file/list`

List uploaded files, paginated.

**Query Parameters**

| Name   | Type | Required | Description                |
|--------|------|----------|------------------------------|
| `page` | int  | Yes      | Page number, must be `> 0`   |

**Responses**

- `200 OK`
  ```json
  {
    "files": [
      { "id": "f1a2...", "ext": ".png" }
    ],
    "pages": 3
  }
  ```

- `400 Bad Request`
  - `{ "message": "query page not set" }`
  - `{ "message": "query page should be int" }`
  - `{ "message": "query page must be greater than 0" }`
  - `{ "message": "page not found" }`

---

### `POST /admin/file/upload`

Upload a file.

**Request:** `multipart/form-data` with a single field named `file`.

**Allowed MIME types**:

- `image/jpeg`, `image/png`, `image/webp`, `image/gif`
- `text/markdown`, `text/plain`
- `application/pdf`
- `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet` (.xlsx)
- `application/vnd.openxmlformats-officedocument.wordprocessingml.document` (.docx)
- `application/vnd.openxmlformats-officedocument.presentationml.presentation` (.pptx)

**Responses**

- `200 OK`
  ```json
  { "message": "file uploaded. File id: <uuid>" }
  ```

- `400 Bad Request` — `{ "message": "no file uploaded" }` / `{ "message": "empty file uploaded" }`
- `415 Unsupported Media Type` — `{ "message": "invalid file uploaded" }`
- `500 Internal Server Error` — `{ "message": "failed to process file" }`

---

### `DELETE /admin/file/delete`

Delete a file by ID.

**Query Parameters**

| Name | Type | Required | Description        |
|------|------|----------|---------------------|
| `id` | UUID | Yes      | File identifier       |

**Responses**

- `200 OK` — `{ "message": "file deleted" }`
- `400 Bad Request`
  - `{ "message": "query id not set" }`
  - `{ "message": "invalid UUID format" }`

---

## Summary Table

| Method | Path                  | Auth  | Description               |
|--------|------------------------|-------|----------------------------|
| GET    | `/health`               | No    | Health check                |
| GET    | `/post/get`             | No    | Get post by ID              |
| GET    | `/post/all`             | No    | List posts (paginated)      |
| GET    | `/file/get`             | No    | Download file by ID         |
| POST   | `/auth/login`           | No (Basic Auth) | Log in, issues `sessionId` cookie |
| POST   | `/auth/logout`          | Requires `sessionId` cookie | Log out, invalidates session |
| POST   | `/admin/post/create`    | Yes   | Create post                 |
| PUT    | `/admin/post/edit`      | Yes   | Edit post                   |
| DELETE | `/admin/post/delete`    | Yes   | Delete post                 |
| GET    | `/admin/file/list`      | Yes   | List files (paginated)      |
| POST   | `/admin/file/upload`    | Yes   | Upload file                 |
| DELETE | `/admin/file/delete`    | Yes   | Delete file                 |