package cinematheque

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type Movie struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Director string   `json:"director"`
	Year     int      `json:"year"`
	Rating   float64  `json:"rating"`
	Cover    string   `json:"cover"`
	Tags     []string `json:"tags"`
}

func (r *Repository) ListMovies() ([]Movie, error) {
	rows, err := r.db.Query(`
		SELECT id, title, director, year, rating, cover
		FROM movies
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		rows.Scan(&m.ID, &m.Title, &m.Director, &m.Year, &m.Rating, &m.Cover)
		movies = append(movies, m)
	}

	return movies, nil
}
