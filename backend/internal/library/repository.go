package library

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

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

func (r *Repository) ListBooks(tags []string) ([]Book, error) {
	query := `SELECT b.id, b.title, b.author, COALESCE(b.isbn, ''), COALESCE(b.year, 0),
		COALESCE(b.rating, 0), COALESCE(b.cover, '')
		FROM books b`
	args := make([]any, 0, len(tags))
	if len(tags) > 0 {
		placeholders := make([]string, len(tags))
		for i, tag := range tags {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			args = append(args, tag)
		}
		query += ` JOIN book_tags bt ON bt.book_id = b.id
			JOIN library_tags t ON t.id = bt.tag_id
			WHERE t.name IN (` + strings.Join(placeholders, ",") + `)
			GROUP BY b.id
			HAVING COUNT(DISTINCT t.name) = ` + fmt.Sprint(len(tags))
	}
	query += " ORDER BY b.id"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]Book, 0)
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Rating, &book.Cover); err != nil {
			return nil, err
		}
		book.Tags, err = r.tagsForBook(book.ID)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, rows.Err()
}

func (r *Repository) GetBook(id int) (Book, error) {
	var book Book
	err := r.db.QueryRow(`SELECT id, title, author, COALESCE(isbn, ''), COALESCE(year, 0),
		COALESCE(rating, 0), COALESCE(cover, '') FROM books WHERE id = $1`, id).
		Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Rating, &book.Cover)
	if err != nil {
		return Book{}, err
	}
	book.Tags, err = r.tagsForBook(id)
	return book, err
}

func (r *Repository) CreateBook(book Book) (Book, error) {
	err := r.db.QueryRow(`INSERT INTO books (title, author, isbn, year, rating, cover)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		book.Title, book.Author, book.ISBN, book.Year, book.Rating, book.Cover).Scan(&book.ID)
	return book, err
}

func (r *Repository) UpdateBook(id int, book Book) (Book, error) {
	_, err := r.db.Exec(`UPDATE books SET title=$1, author=$2, isbn=$3, year=$4, rating=$5, cover=$6 WHERE id=$7`,
		book.Title, book.Author, book.ISBN, book.Year, book.Rating, book.Cover, id)
	if err != nil {
		return Book{}, err
	}
	return r.GetBook(id)
}

func (r *Repository) DeleteBook(id int) error {
	result, err := r.db.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) tagsForBook(id int) ([]string, error) {
	rows, err := r.db.Query(`SELECT t.name FROM library_tags t
		JOIN book_tags bt ON bt.tag_id = t.id WHERE bt.book_id = $1 ORDER BY t.name`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}
