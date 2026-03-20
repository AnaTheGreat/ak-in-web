## Backend structure:

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

## Database structure:

Library:
- books
- library_tags
- book_tags (many-to-many)

Cinematheque:
- movies
- cinematheque_tags
- movie_tags (many-to-many)
