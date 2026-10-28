package service

import (
    "log"
    "strings"
    "time"

    "github.com/matin-sp/-online-library-management-system/internal/cache"
    "github.com/matin-sp/-online-library-management-system/internal/model"
    "github.com/matin-sp/-online-library-management-system/internal/repository"
)

const authorCacheTTL = 10 * time.Minute

var (
    ErrAuthorNotFound = repository.ErrAuthorNotFound
    ErrAuthorHasBooks = repository.ErrAuthorHasBooks
)

type AuthorService struct {
    repo  *repository.AuthorRepository
    cache *cache.Cache
}

func NewAuthorService(repo *repository.AuthorRepository, appCache *cache.Cache) *AuthorService {
    return &AuthorService{repo: repo, cache: appCache}
}

func (s *AuthorService) Create(author model.Author) (*model.Author, error) {
    return s.repo.Create(author) // a brand-new id has no cached copy, nothing to invalidate
}

func (s *AuthorService) GetByID(id int) (*model.Author, error) {
    key := cache.AuthorKey(id)

    var cached model.Author
    if s.cache.GetJSON(key, &cached) {
        log.Println("cache HIT ", key)
        return &cached, nil
    }
    log.Println("cache MISS", key)

    author, err := s.repo.GetByID(id)
    if err != nil {
        return nil, err // "not found" is NOT cached
    }
    s.cache.SetJSON(key, author, authorCacheTTL)
    return author, nil
}


func (s *AuthorService) List(f model.AuthorFilter) (*model.Page[model.Author], error) {
   
    f.Name = strings.TrimSpace(f.Name)

    if f.Page < 1 {
        f.Page = model.DefaultPage
    }
    if f.Limit < 1 {
        f.Limit = model.DefaultLimit
    }
    if f.Limit > model.MaxLimit {
        f.Limit = model.MaxLimit
    }

    authors, total, err := s.repo.List(f)
    if err != nil {
        return nil, err
    }
    return model.NewPage(authors, f.Page, f.Limit, total), nil
}

func (s *AuthorService) Update(author model.Author) error {
    if err := s.repo.Update(author); err != nil {
        return err
    }
    s.cache.Delete(cache.AuthorKey(author.ID)) // database first, then invalidate
    return nil
}

func (s *AuthorService) Delete(id int) error {
    if err := s.repo.Delete(id); err != nil {
        return err
    }
    s.cache.Delete(cache.AuthorKey(id))
    return nil
}