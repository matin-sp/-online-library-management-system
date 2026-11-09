package service

import (
	"errors"

	"github.com/matin-sp/-online-library-management-system/internal/model"
	"github.com/matin-sp/-online-library-management-system/internal/repository"
)

var (
	ErrCopyNotFound       = repository.ErrCopyNotFound
	ErrCopyOnLoan         = repository.ErrCopyOnLoan
	ErrCopyNumberConflict = repository.ErrCopyNumberConflict
)

var ErrInvalidCopyStatus = errors.New("status must be one of: available, under_repair, lost, damaged")

// settableCopyStatuses are the statuses a librarian may set by hand.

var settableCopyStatuses = map[string]bool{
	model.CopyStatusAvailable:   true,
	model.CopyStatusUnderRepair: true,
	model.CopyStatusLost:        true,
	model.CopyStatusDamaged:     true,
}

// BookCopyService has no cache on purpose: the status of a copy changes
// every time a book is borrowed or returned, so a cached copy would be wrong.
type BookCopyService struct {
	repo *repository.BookCopyRepository
}

func NewBookCopyService(repo *repository.BookCopyRepository) *BookCopyService {
	return &BookCopyService{repo: repo}
}

// Create adds one new copy of a book; the repository picks the copy_number.
func (s *BookCopyService) Create(bookID int) (*model.BookCopy, error) {
	return s.repo.Create(bookID)
}

func (s *BookCopyService) ListByBook(bookID int) ([]model.BookCopy, error) {
	return s.repo.ListByBook(bookID)
}

// UpdateStatus checks the business rule (which statuses may be set by hand)
// and then asks the repository to change the status.
func (s *BookCopyService) UpdateStatus(id int, status string) (*model.BookCopy, error) {
	if !settableCopyStatuses[status] {
		return nil, ErrInvalidCopyStatus
	}
	return s.repo.UpdateStatus(id, status)
}
