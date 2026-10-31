# AK In Web (Work in Progress)

This is a full-stack personal website project intended to teach practical
application delivery, DevOps, and later DevSecOps skills.

## Project goal

The public website will contain personal content in two sections:

- Library: books and tags
- Cinematheque: movies and tags

The future admin area will allow the owner to manage that content. There is no
public user registration planned for the first version.

## Current backend features

- Health check: `GET /health`
- Automatic database migrations on backend startup
- PostgreSQL readiness checks in Docker
- Book CRUD: `/api/library/books`
- Movie CRUD: `/api/cinematheque/movies`
- Basic input validation and JSON errors
- AND filtering by existing tag names: `?tags=tag-one,tag-two`

Authentication, tag management, image uploads, and the frontend are later
milestones.

## Run locally

Copy `.env.example` to `.env` if you want to customize local values, then run:

```bash
docker compose up --build
```

The backend waits for PostgreSQL, applies pending files in `backend/migrations`,
and listens on `http://localhost:8080`.

If your Docker installation does not include the Compose plugin, install Docker
Compose or use the equivalent individual `docker run` commands.

Verify the service:

```bash
curl http://localhost:8080/health
```

## API examples

Create a book:

```bash
curl -X POST http://localhost:8080/api/library/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Dune","author":"Frank Herbert","year":1965,"rating":9.5}'
```

List books, optionally using AND tag filtering:

```bash
curl "http://localhost:8080/api/library/books?tags=favorite,science-fiction"
```

The same CRUD routes are available for movies under
`/api/cinematheque/movies`.

## Architecture

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).
