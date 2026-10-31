package cinematheque

import (
	"errors"
	"strings"
)

var ErrValidation = errors.New("validation error")

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service                  { return &Service{repo: repo} }
func (s *Service) GetMovies(tags []string) ([]Movie, error) { return s.repo.ListMovies(tags) }
func (s *Service) GetMovie(id int) (Movie, error)           { return s.repo.GetMovie(id) }
func (s *Service) CreateMovie(movie Movie) (Movie, error) {
	if err := validateMovie(movie); err != nil {
		return Movie{}, err
	}
	return s.repo.CreateMovie(movie)
}
func (s *Service) UpdateMovie(id int, movie Movie) (Movie, error) {
	if err := validateMovie(movie); err != nil {
		return Movie{}, err
	}
	return s.repo.UpdateMovie(id, movie)
}
func (s *Service) DeleteMovie(id int) error { return s.repo.DeleteMovie(id) }

func validateMovie(movie Movie) error {
	if strings.TrimSpace(movie.Title) == "" || strings.TrimSpace(movie.Director) == "" {
		return ErrValidation
	}
	if movie.Year < 0 || movie.Year > 3000 || movie.Rating < 0 || movie.Rating > 10 {
		return ErrValidation
	}
	return nil
}
