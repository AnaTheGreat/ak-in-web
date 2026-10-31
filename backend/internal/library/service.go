package library

import (
	"errors"
	"strings"
)

var ErrValidation = errors.New("validation error")

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func (s *Service) GetBooks(tags []string) ([]Book, error) { return s.repo.ListBooks(tags) }
func (s *Service) GetBook(id int) (Book, error)           { return s.repo.GetBook(id) }
func (s *Service) CreateBook(book Book) (Book, error) {
	if err := validateBook(book); err != nil {
		return Book{}, err
	}
	return s.repo.CreateBook(book)
}
func (s *Service) UpdateBook(id int, book Book) (Book, error) {
	if err := validateBook(book); err != nil {
		return Book{}, err
	}
	return s.repo.UpdateBook(id, book)
}
func (s *Service) DeleteBook(id int) error { return s.repo.DeleteBook(id) }

func validateBook(book Book) error {
	if strings.TrimSpace(book.Title) == "" || strings.TrimSpace(book.Author) == "" {
		return ErrValidation
	}
	if book.Year < 0 || book.Year > 3000 || book.Rating < 0 || book.Rating > 10 {
		return ErrValidation
	}
	return nil
}
