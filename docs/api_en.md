# 🚀 API Documentation (v1.0)

Interface for managing posts and media files. All timestamps are provided in **ISO 8601** format.

---

## 📝 Post Management (`/api/post`)

Base methods for users to read content.

### Get Post by ID
`GET /api/post/get`

| Parameter | Type | Required | Description |
| :--- | :--- | :---: | :--- |
| `id` | `UUID` | ✅ | Unique post identifier |

**Success Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Post Title",
  "poster_id": "330e8400-e29b-41d4-a716-446655440000",
  "created_at": "2026-04-24T21:41:48.123Z",
  "updated_at": null,
  "body": "Post content in Markdown format"
}
```

---

### List All Posts (Pagination)
`GET /api/post/all`

| Parameter | Type | Required | Description |
| :--- | :--- | :---: | :--- |
| `page` | `Integer` | ✅ | Page number (starting from 1) |

**Success Response (200 OK):**
```json
{
  "posts": [ { ... }, { ... } ],
  "pages": 10
}
```

---

## 🛠 Post Administration (`/api/admin/post`)

> **Note:** These methods require administrator privileges.

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **POST** | `/create` | Create a new post |
| **PUT** | `/edit` | Edit an existing post |
| **DELETE** | `/delete` | Permanently delete a post |

### Create and Edit
A JSON body is used for both `/create` and `/edit`:

```json
{
  "id": "UUID", // Required for /edit only
  "title": "New Title",
  "poster_id": "UUID",
  "body": "Post content"
}
```

---

## 📂 File Handling (`/api/file`)

Storage system for images and documents.

### 📤 Upload File
`POST /api/admin/file/upload`

**Content-Type:** `multipart/form-data`

**Supported MIME Types:**
* **Images:** `jpg`, `png`, `webp`, `gif`
* **Text:** `markdown`, `plain`
* **Documents:** `pdf`, `docx`, `xlsx`, `pptx`

**Response:**
```json
{
  "message": "file uploaded",
  "file_id": "d290f1ee-6c54-4b01-90e6-d701748f0851"
}
```

---

### 📥 Get File
`GET /api/file/get?id={UUID}`

Returns the file as a binary stream (**Bytestream**). The response headers will specify the corresponding `Content-Type`.

---

## ⚠️ Error Handling

All errors are returned in JSON format with a description of the cause.

| Code (HTTP) | Message | Cause |
| :--- | :--- | :--- |
| **400** | `query id not set` | Missing required URL parameter |
| **400** | `invalid UUID format` | Invalid identifier format provided |
| **400** | `query page must be greater than 0` | Pagination error |
| **404** | `page not found` | The requested page does not exist |
| **415** | `invalid file uploaded` | Unsupported file format |
| **500** | `internal server error` | Critical server-side error |