package library

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetBooks() ([]Book, error) {
	return s.repo.ListBooks()
}

func (s *Service) CreateBook(b Book) (Book, error) {
	return s.repo.CreateBook(b)
}
