package library

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type Book struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Author string   `json:"author"`
	ISBN   string   `json:"isbn"`
	Year   int      `json:"year"`
	Rating float64  `json:"rating"`
	Cover  string   `json:"cover"`
	Tags   []string `json:"tags"`
}

func (r *Repository) ListBooks() ([]Book, error) {
	rows, err := r.db.Query(`
		SELECT id, title, author, isbn, year, rating, cover
		FROM books
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		rows.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Year, &b.Rating, &b.Cover)
		books = append(books, b)
	}

	return books, nil
}

func (r *Repository) CreateBook(b Book) (Book, error) {
	err := r.db.QueryRow(`
		INSERT INTO books (title, author, isbn, year, rating, cover)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`,
		b.Title,
		b.Author,
		b.ISBN,
		b.Year,
		b.Rating,
		b.Cover,
	).Scan(&b.ID)

	return b, err
}
