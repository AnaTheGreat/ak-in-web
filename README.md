# AK In Web (Work in Progress)

This project is a full-stack, personal portfolio, website built by using the following skills & tools:

- Go (backend)
- React (frontend)
- PostgreSQL (database)
- Docker (containerization)
- REST API architecture

It is currently under active development. The goal is not just to make it work, but to build a system with proper and clean backend structure, database design and DevOps practices rather than rushing features. Personal goal is preparation for DevOps/Backend roles.

---

## 🧠 Project Goal

The application will support a personal website with two main data-driven sections:

- Library (books + tags)
- Cinematheque (movies + tags)

Each section is independent and will allow:

- Creating entries (books/movies)
- Tagging entries
- Filtering by tags
- Managing content through an admin interface

---

## 🏗️ Architecture

The backend follows a layered structure:

Handler (HTTP layer)  
↓  
Service (business logic)  
↓  
Repository (database access)  

Each domain is separated:

backend/internal/
  library/
    handler.go
    service.go
    repository.go

  cinematheque/
    handler.go
    service.go
    repository.go

This keeps the system modular and scalable.

---

## 📦 Current Features

- Health check endpoint (`/health`)
- List books (`GET /library/books`)
- Create book (`POST /library/books`)
- List movies (`GET /cinematheque/movies`)
- PostgreSQL integration
- Dockerized environment

---

## 🚧 Work in Progress

The following features are currently being implemented:

- Migrations
- Tag system (many-to-many relationships)
- Filtering by tags
- Admin authentication
- Update/Delete operations
- Input validation
- Better error handling

---

## 🗄️ Database Design

Planned structure:

Library:
- books
- library_tags
- book_tags (many-to-many)

Cinematheque:
- movies
- cinematheque_tags
- movie_tags (many-to-many)

---

## ▶️ How to Run

Start services:

docker-compose up -d

Run backend locally (optional):

cd backend  
go run ./cmd/api  

Test health endpoint:

curl http://localhost:8080/health

---

## 🔍 Example Request

Create a book:

curl -X POST http://localhost:8080/library/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Book",
    "author": "AK",
    "year": 2024,
    "rating": 8.5
  }'
