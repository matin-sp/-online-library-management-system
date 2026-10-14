package service

import (
	"github.com/matin-sp/-online-library-management-system/internal/model"
	"github.com/matin-sp/-online-library-management-system/internal/repository"
)

var (
	ErrBookNotFound  = repository.ErrBookNotFound
	ErrInvalidAuthor = repository.ErrInvalidAuthor
)


type BookService struct {
	repo *repository.BookRepository
}

func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) Create(book model.Book) (*model.Book, error) {
	return s.repo.Create(book)
}

func (s *BookService) GetByID(id int) (*model.Book, error) {
	return s.repo.GetByID(id)
}

func (s *BookService) GetByIDWithAuthor(id int) (*repository.BookWithAuthor, error) {
	return s.repo.GetByIDWithAuthor(id)
}

func (s *BookService) GetAll() ([]model.Book, error) {
	return s.repo.GetAll()
}

func (s *BookService) GetByAuthorID(authorID int) ([]model.Book, error) {
	return s.repo.GetByAuthorID(authorID)
}

func (s *BookService) Update(book model.Book) error {
	return s.repo.Update(book)
}

func (s *BookService) Delete(id int) error {
	return s.repo.Delete(id)
}