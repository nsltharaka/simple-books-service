# 📚 Simple Book Service

A simple REST API for managing books, built with Go, Fiber, and SQLite.

---

## 🚀 Features

- CRUD operations for books
- Request validation using `validator.v10`
- Pagination support with `?page=1&limit=10`
- Consistent JSON response format
- Structured logging with `slog`
- Unit-tested service and handler layers

---

## 🛠 Tech Stack

- [Go](https://golang.org/)
- [Fiber](https://gofiber.io/)
- [GORM](https://gorm.io/)
- [Validator.v10](https://pkg.go.dev/github.com/go-playground/validator/v10)
- [slog](https://pkg.go.dev/log/slog)

---

## 📦 Getting Started

### 1. Clone the repo

```bash
$ git clone https://github.com/nsltharaka/book-api.git
$ cd book-api
```

### 2. Set environment variables

```bash
$ mv .env.example .env
```

### 3. Install Dependencies

- if you have make tool installed in your system,

```bash
$ make install
```

- or run

```bash
$ go mod tidy
```

### 4. Start the server

- if you have make tool installed in your system,

```bash
$ make run
```

- or,

```bash
$ go run .
```

## 📘 API Endpoints

### Create a Book

_POST /books_

- 201 on success
- 400 if request body is invalid or missing required fields
- 500 if something unexpected happens on the server
- example request body

```json
{
  "title": "The Da Vinci Code",
  "author": "Dan Brown",
  "year": 2003
}
```

- example response

```json
{
  "message": "success",
  "data": {
    "id": 1,
    "title": "The Da Vinci Code",
    "author": "Dan Brown",
    "year": 2003
  }
}
```

- Test with curl

```bash
curl -X POST http://localhost:3030/books \
  -H "Content-Type: application/json" \
  -d '{"title": "The Da Vinci Code", "author": "Dan Brown", "year": 2003}'
```

---

### Get all books

_GET /books_

_GET /books?page=1&limit=10_

- default page = 1, limit = 10
- default values are used if malformed values are passed
- 200 on success
- 500 if an unexpected error occurs
- example response

```json
{
  "message": "success",
  "data": [
    {
      "id": 1,
      "title": "The Da Vinci Code",
      "author": "Dan Brown",
      "year": 2003
    },
    {
      "id": 2,
      "title": "Angels & Demons",
      "author": "Dan Brown",
      "year": 2000
    }
  ]
}
```

- Test with curl

```bash
curl http://localhost:3030/books?page=1&limit=10
```

---

### Get Book by ID

_GET /books/:id_

- 200 on success
- 400 if ID is not a valid number
- 404 if book not found
- 500 if an unexpected error occurs
- unsuccessful request returns 404 status code
  - eg: no book found with the given ID
- example response

```json
{
  "message": "success",
  "data": {
    "id": 1,
    "title": "The Da Vinci Code",
    "author": "Dan Brown",
    "year": 2003
  }
}
```

- Test with curl

```bash
curl http://localhost:3030/books/1
```

---

### Update a Book

_PUT /books/:id_

- 200 on success
- 400 if ID is invalid, request body is malformed, or fails validation
- 404 if book with given ID does not exist
- 500 if an unexpected error occurs
- example request body

```json
{
  "title": "Angels & Demons",
  "author": "Dan Brown",
  "year": 2000
}
```

- example response

```json
{
  "message": "success",
  "data": {
    "id": 1,
    "title": "Angels & Demons",
    "author": "Dan Brown",
    "year": 2000
  }
}
```

- Test with curl

```bash
curl -X PUT http://localhost:3030/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Angels & Demons", "author": "Dan Brown", "year": 2000}'
```

---

### Delete a Book

_DELETE /books/:id_

- 200 on success
- 400 if ID is not a valid number
- 404 if book not found
- 500 if an unexpected error occurs
- example response

```json
{
  "message": "success",
  "data": {
    "id": 1,
    "title": "Angels & Demons",
    "author": "Dan Brown",
    "year": 2000
  }
}
```

- Test with curl

```bash
curl -X DELETE http://localhost:3030/books/1
```

## ✅ API Response format

```json
{
  "data": {},
  "error": null,
  "message": "success"
}
```

## 🧪 Running Tests

- if you have make tool installed in your system,

```bash
$ make test
```

- or run

```bash
$ go test ./...
```
