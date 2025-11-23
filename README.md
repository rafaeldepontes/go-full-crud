# Go Full CRUD

Simple RESTful CRUD API for user management written in Go.

## Overview

This repository contains a small HTTP API that manages users backed by PostgreSQL. It was built using **Go 1.25** and assumes a PostgreSQL instance is available, if not... I've provided a Docker command to run one locally.

## Technologies used

- Go 1.25
- PostgreSQL (run with Docker)
- chi (router)
- jackc/pgx driver
- logrus for logging
- gorilla/schema for query decoding
- godotenv for environment variables

## Requirements

- Go 1.25 or later installed on your machine
- Docker (for running PostgreSQL locally) or an accessible Postgres instance
- Make sure the `DATABASE_URL` environment variable points to a reachable Postgres database

## Environment variables

Create a `.env` file (example provided in the project). The service expects at least:

```
ADDRESS="localhost:8080"           # HTTP listen address (overrides default 8000)
DATABASE="Postgres"                # Informational
DATABASE_URL="postgres://root:example@localhost:5432/golang-database"
```

## Database schema

Run this SQL in your PostgreSQL database to create the `users` table used by the application:

```SQL
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  email TEXT,
  birthdate TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP
);
```

## How to use

1. Clone the repo:

```bash
git clone <repo-url>
cd go-full-crud
```

2. Create a `.env` file (or use the provided `.env.example`) and set `DATABASE_URL` and `ADDRESS` as needed.

3. Start PostgreSQL with Docker:

```bash
docker run --name go-postgres -e POSTGRES_PASSWORD=example -e POSTGRES_USER=root -e POSTGRES_DB=golang-database -p 5432:5432 -d postgres:15
```

Adjust user/password/db name to match your `DATABASE_URL` if necessary.

4. Apply the database schema (run the SQL above) using `psql` or a GUI tool.

5. Run the service:

```bash
# from the project root
go run cmd/go-crud/main.go
```

By default the service will listen on the address set in `ADDRESS` env var (default `localhost:8000` if not set). If you set `ADDRESS="localhost:8080"` in `.env`, the service will listen on `localhost:8080` (main.go reads the env variable).

## Endpoints

All endpoints are prefixed with `/api/v1`.

### Get user by username (query param)

- **Route:** `GET /api/v1/users?username=<username>`
- **Success:** `200 OK` with user JSON
- **Errors:** `400 Bad Request` if username blank, `404 Not Found` if no user

```bash
curl -G "http://localhost:8080/api/v1/users" --data-urlencode "username=johndoe"
```

### Get user by id

- **Route:** `GET /api/v1/users/{id}`
- **Success:** `200 OK` with user JSON
- **Errors:** `400 Bad Request` if id blank/invalid, `404 Not Found` if not found

```bash
curl "http://localhost:8080/api/v1/users/1"
```

### Register user

- **Route:** `POST /api/v1/users`
- **Body:**
  - JSON
  ```json
  {
    "username": "johndoe",
    "password": "s3cret"
  }
  ```
- **Success:** `201 Created`
- **Errors:** `400 Bad Request` if missing fields, `500 Internal Server Error` for server errors

```bash
curl -X POST "http://localhost:8080/api/v1/users" \
  -H "Content-Type: application/json" \
  -d '{"username":"johndoe","password":"s3cret"}'
```

### Update user info (email, birthdate)

- **Route:** `PUT /api/v1/users/{id}`
- **Body:**
  -JSON with optional fields (use `null` to clear), e.g.
  ```json
  {
    "email": "user@example.com",
    "birthdate": "1990-01-01"
  }
  ```
- **Success:** `200 OK`
- **Errors:** `400 Bad Request` if id blank/invalid, `500 Internal Server Error` for server errors

```bash
curl -X PUT "http://localhost:8080/api/v1/users/1" \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","birthdate":"1990-01-01"}'
```

_Note:_ The code uses pointer fields for `email` and `birthdate` so you can omit them from the JSON to leave values unchanged or send them as `null` to explicitly clear them (depending on your repo/sql behavior).

### Delete user by id

- **Route:** `DELETE /api/v1/users/{id}`
- **Success:** `204 No Content`
- **Errors:** `400 Bad Request` if id blank/invalid, `500 Internal Server Error` for server errors

```bash
curl -X DELETE "http://localhost:8080/api/v1/users/1"
```

## Error responses

The API returns JSON error objects with this shape:

```json
{
  "status": 400,
  "message": "Description of the error",
  "timestamp": "02/01/2006 15:04:05"
}
```

(Format uses the BrazilianDateTimeFormat constant in the code: `DD/MM/YYYY HH:MM:SS`.)

## Notes & tips

- Passwords are stored as-is because hashing or providing an authentication flow is outside the purpose of this project. The focus is on showcasing CRUD operations and clean handler/repository structure.
- Database migrations are intentionally not included to keep the setup simpler and more direct for learning.
- Additional concerns such as authentication, authorization, password hashing, request validation, and production-grade logging are intentionally left out to keep the example focused on CRUD fundamentals.

---
