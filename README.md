# Go JWT Authentication Backend (Containerized)

![Go](https://img.shields.io/badge/Go-1.23-blue?style=flat&logo=go)
![Docker](https://img.shields.io/badge/Docker-Production%20Ready-2496ED?style=flat&logo=docker)
![Postgres](https://img.shields.io/badge/Postgres-15-336791?style=flat&logo=postgresql)
![JWT](https://img.shields.io/badge/Auth-JWT%20%2B%20Refresh%20Tokens-green)

A professional, production-ready REST API boilerplate written in Go (Golang). It features a complete Authentication system using **UUIDs**, **JWT Access Tokens**, and **Secure Refresh Tokens** stored in PostgreSQL. The entire stack is fully containerized with Docker, including automated database migrations running as a dedicated service.

---

## 🚀 Features

* **Production Architecture**: Follows the "Standard Go Project Layout" (`cmd`, `internal`, `db`).
* **Containerization**: Multi-stage Docker builds (Distroless/Alpine) resulting in tiny, secure images.
* **Database**: PostgreSQL 15 with `uuid-ossp` extension enabled for UUID primary keys.
* **Authentication Strategy**:
    * **Access Tokens**: Short-lived (15 min) JWTs signed with HMAC-SHA256.
    * **Refresh Tokens**: Long-lived (7 days) opaque tokens, hashed (SHA256) and stored in the DB for security and revocation.
    * **Password Security**: Passwords hashed using `bcrypt`.
* **Automated Migrations**: Uses `golang-migrate` running in a dedicated container to handle schema changes on startup.
* **Routing**: Lightweight, idiomatic routing using `go-chi`.

---

## 🛠️ Tech Stack

* **Language**: Go 1.23+
* **Database**: PostgreSQL 15
* **Router**: [Chi v5](https://github.com/go-chi/chi)
* **Database Driver**: [lib/pq](https://github.com/lib/pq)
* **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
* **Infrastructure**: Docker & Docker Compose

---

## 📂 Project Structure

```text
.
├── cmd/
│   └── server/
│       └── main.go           # Application Entry Point
├── db/
│   └── migrations/           # SQL Migration files (Up/Down)
│       ├── 000001_init_schema.up.sql
│       └── ...
├── internal/
│   ├── auth/                 # JWT generation & SHA256 hashing utilities
│   ├── config/               # Database Connection & Retry logic
│   ├── controllers/          # HTTP Handlers (Login, Register, Refresh)
│   ├── models/               # Database Structs & Repositories
│   ├── utils/                # Helper functions
│   └── routes/               # Route definitions
├── Dockerfile                # Multi-stage production build
├── docker-compose.yml        # Orchestration (DB, App, Migrate)
├── go.mod                    # Dependencies
└── .env                      # Environment Variables (Optional for local dev)
```
## ⚡ Quick Start (Docker)

You do not need Go or Postgres installed on your machine. You only need **Docker**.

### 1. Run the Project
This command builds the Go binary, starts Postgres, and runs the SQL migrations automatically.

```bash
docker compose up -d --build
```

### 2. Verify Status
Ensure all three containers (`go_jwt_docker-app`, `go_pgsql_db`, `go_db_migrate`) are running.

```bash
docker ps
```

### 3. Access
The API will be available at: `http://localhost:8080`

* **API URL**: `http://localhost:8080`
* **Database Port (External)**: `localhost:5433` (Mapped to prevent conflicts with local Postgres)

---

## 🔌 API Documentation

### 1. Register User
Creates a new user with an auto-generated UUID.

* **URL**: `/v1/user/register`
* **Method**: `POST`
* **Body**:
    ```json
    {
        "email": "user@example.com",
        "password": "strongpassword123",
        "age": 25
    }
    ```
* **Response**: `201 Created`

### 2. Login
Authenticates user and returns an Access Token (JWT) and a Refresh Token (Raw string).

* **URL**: `/v1/user/login`
* **Method**: `POST`
* **Body**:
    ```json
    {
        "email": "user@example.com",
        "password": "strongpassword123"
    }
    ```
* **Response**: `200 OK`
    ```json
    {
        "message": "Login successful",
        "access_token": "eyJhbGciOiJI...",  // Valid for 15 mins
        "refresh_token": "Mj4tpqnJgwi...",  // Valid for 7 days
        "user": {
            "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
            "email": "user@example.com",
            "age": 25
        }
    }
    ```

### 3. Refresh Token
Exchanges a valid Refresh Token for a new Access Token when the JWT expires.

* **URL**: `/v1/user/refresh`
* **Method**: `POST`
* **Body**:
    ```json
    {
        "refresh_token": "Mj4tpqnJgwi..."
    }
    ```
* **Response**: `200 OK`
    ```json
    {
        "access_token": "eyJhbGciOiJI..."
    }
    ```

---

## 🗄️ Database Migrations

This project treats database changes as version-controlled code.

### The Migration Container
A dedicated container (`go_db_migrate`) runs automatically when you start Docker. It checks the `db/migrations` folder and applies any new SQL files to the Postgres database.

### Creating a New Migration
To add a new table (e.g., `orders`), run this command (requires `migrate` CLI installed locally):

```bash
migrate create -ext sql -dir db/migrations -seq create_orders_table
```

Then edit the newly created `.up.sql` and `.down.sql` files.

### Resetting the Database (Development Only)
If you need to wipe all data and start fresh (e.g., if migration state gets dirty):

```bash
# 1. Tear down volumes (Deletes all data)
docker compose down -v
```
