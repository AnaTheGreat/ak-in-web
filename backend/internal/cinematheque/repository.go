package cinematheque

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

type Movie struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Director string   `json:"director"`
	Year     int      `json:"year"`
	Rating   float64  `json:"rating"`
	Cover    string   `json:"cover"`
	Tags     []string `json:"tags"`
}

func (r *Repository) ListMovies(tags []string) ([]Movie, error) {
	query := `SELECT m.id, m.title, m.director, COALESCE(m.year, 0),
		COALESCE(m.rating, 0), COALESCE(m.cover, '') FROM movies m`
	args := make([]any, 0, len(tags))
	if len(tags) > 0 {
		placeholders := make([]string, len(tags))
		for i, tag := range tags {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			args = append(args, tag)
		}
		query += ` JOIN movie_tags mt ON mt.movie_id = m.id
			JOIN cinematheque_tags t ON t.id = mt.tag_id
			WHERE t.name IN (` + strings.Join(placeholders, ",") + `)
			GROUP BY m.id
			HAVING COUNT(DISTINCT t.name) = ` + fmt.Sprint(len(tags))
	}
	query += " ORDER BY m.id"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	movies := make([]Movie, 0)
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Director, &movie.Year, &movie.Rating, &movie.Cover); err != nil {
			return nil, err
		}
		movie.Tags, err = r.tagsForMovie(movie.ID)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, rows.Err()
}

func (r *Repository) GetMovie(id int) (Movie, error) {
	var movie Movie
	err := r.db.QueryRow(`SELECT id, title, director, COALESCE(year, 0),
		COALESCE(rating, 0), COALESCE(cover, '') FROM movies WHERE id = $1`, id).
		Scan(&movie.ID, &movie.Title, &movie.Director, &movie.Year, &movie.Rating, &movie.Cover)
	if err != nil {
		return Movie{}, err
	}
	movie.Tags, err = r.tagsForMovie(id)
	return movie, err
}

func (r *Repository) CreateMovie(movie Movie) (Movie, error) {
	err := r.db.QueryRow(`INSERT INTO movies (title, director, year, rating, cover)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		movie.Title, movie.Director, movie.Year, movie.Rating, movie.Cover).Scan(&movie.ID)
	return movie, err
}

func (r *Repository) UpdateMovie(id int, movie Movie) (Movie, error) {
	_, err := r.db.Exec(`UPDATE movies SET title=$1, director=$2, year=$3, rating=$4, cover=$5 WHERE id=$6`,
		movie.Title, movie.Director, movie.Year, movie.Rating, movie.Cover, id)
	if err != nil {
		return Movie{}, err
	}
	return r.GetMovie(id)
}

func (r *Repository) DeleteMovie(id int) error {
	result, err := r.db.Exec("DELETE FROM movies WHERE id = $1", id)
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

func (r *Repository) tagsForMovie(id int) ([]string, error) {
	rows, err := r.db.Query(`SELECT t.name FROM cinematheque_tags t
		JOIN movie_tags mt ON mt.tag_id = t.id WHERE mt.movie_id = $1 ORDER BY t.name`, id)
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
