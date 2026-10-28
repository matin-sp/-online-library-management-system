package service

import (
	"log"
	"time"

	"github.com/matin-sp/-online-library-management-system/internal/cache"
	"github.com/matin-sp/-online-library-management-system/internal/model"
	"github.com/matin-sp/-online-library-management-system/internal/repository"
)

const bookCacheTTL = 10 * time.Minute

var (
	ErrBookNotFound  = repository.ErrBookNotFound
	ErrInvalidAuthor = repository.ErrInvalidAuthor
)

type BookService struct {
	repo  *repository.BookRepository
	cache *cache.Cache
}

func NewBookService(repo *repository.BookRepository, appCache *cache.Cache) *BookService {
	return &BookService{repo: repo, cache: appCache}
}

func (s *BookService) Create(book model.Book) (*model.Book, error) {
	return s.repo.Create(book) // a brand-new id has no cached copy, nothing to invalidate
}

// GetByID uses the cache-aside pattern: Redis first, then Postgres.
func (s *BookService) GetByID(id int) (*model.Book, error) {
	key := cache.BookKey(id)

	var cached model.Book
	if s.cache.GetJSON(key, &cached) {
		log.Println("cache HIT ", key)
		return &cached, nil
	}
	log.Println("cache MISS", key)

	book, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err // "not found" is NOT cached
	}
	s.cache.SetJSON(key, book, bookCacheTTL)
	return book, nil
}

// GetByIDWithAuthor is not cached yet. It contains the author's name, so an
// author rename would make the cached copy wrong. See the invalidation
// discussion before adding a cache here.
func (s *BookService) GetByIDWithAuthor(id int) (*repository.BookWithAuthor, error) {
	return s.repo.GetByIDWithAuthor(id)
}

func (s *BookService) GetAll() ([]model.Book, error) {
	return s.repo.GetAll()
}

func (s *BookService) GetByAuthorID(authorID int) ([]model.Book, error) {
	return s.repo.GetByAuthorID(authorID)
}

// Update changes the database first, then removes the cached copy.
func (s *BookService) Update(book model.Book) error {
	if err := s.repo.Update(book); err != nil {
		return err
	}
	s.cache.Delete(cache.BookKey(book.ID))
	return nil
}

func (s *BookService) Delete(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.cache.Delete(cache.BookKey(id))
	return nil
}