package service

import (
	"github.com/matin-sp/-online-library-management-system/internal/model"
	"github.com/matin-sp/-online-library-management-system/internal/repository"
)

// Re-exported so the handler layer can check errors without importing the
// repository package (handler -> service -> repository stays one-directional).
var (
	ErrAuthorNotFound = repository.ErrAuthorNotFound
	ErrAuthorHasBooks = repository.ErrAuthorHasBooks
)

// AuthorService holds the business logic for authors. Today it just passes
// calls through; rules like "no duplicate author names" would be added here,
// not in the handler (HTTP) or the repository (SQL).
type AuthorService struct {
	repo *repository.AuthorRepository
}

func NewAuthorService(repo *repository.AuthorRepository) *AuthorService {
	return &AuthorService{repo: repo}
}

func (s *AuthorService) Create(author model.Author) (*model.Author, error) {
	return s.repo.Create(author)
}

func (s *AuthorService) GetByID(id int) (*model.Author, error) {
	return s.repo.GetByID(id)
}

func (s *AuthorService) GetAll() ([]model.Author, error) {
	return s.repo.GetAll()
}

func (s *AuthorService) Update(author model.Author) error {
	return s.repo.Update(author)
}

func (s *AuthorService) Delete(id int) error {
	return s.repo.Delete(id)
}