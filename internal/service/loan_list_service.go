package service

import (
	"errors"
	"time"

	"github.com/matin-sp/-online-library-management-system/internal/model"
)

var ErrInvalidLoanStatus = errors.New("invalid loan status")


type MyFines struct {
	Fines       *model.Page[model.Fine]
	TotalUnpaid int64
}

// MyLoans returns one page of the loans of one user. status is "" or one of
// active, overdue, returned. These are plain reads, so they use s.db directly
// and no transaction is needed.
func (s *LoanService) MyLoans(userID int, status string, page, limit int) (*model.Page[model.LoanDetail], error) {
	if status != "" && !model.IsValidLoanStatus(status) {
		return nil, ErrInvalidLoanStatus
	}
	total, err := s.loans.CountByUser(s.db, userID, status)
	if err != nil {
		return nil, err
	}
	loans, err := s.loans.ListByUser(s.db, userID, status, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	return model.NewPage(loans, page, limit, total), nil
}

// OverdueLoans returns one page of all overdue loans, for the librarian.
func (s *LoanService) OverdueLoans(page, limit int) (*model.Page[model.OverdueLoan], error) {
	total, err := s.loans.CountOverdue(s.db)
	if err != nil {
		return nil, err
	}
	loans, err := s.loans.ListOverdue(s.db, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	for i := range loans {
		// Same rule as the fine: a started day counts as a full day.
		loans[i].DaysOverdue = lateDays(loans[i].DueDate, now)
	}
	return model.NewPage(loans, page, limit, total), nil
}

// MyFines returns one page of the fines of one user, plus the unpaid total.
func (s *LoanService) MyFines(userID, page, limit int) (*MyFines, error) {
	total, err := s.fines.CountByUser(s.db, userID)
	if err != nil {
		return nil, err
	}
	fines, err := s.fines.ListByUser(s.db, userID, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	unpaid, err := s.fines.SumUnpaidByUser(s.db, userID)
	if err != nil {
		return nil, err
	}
	return &MyFines{Fines: model.NewPage(fines, page, limit, total), TotalUnpaid: unpaid}, nil
}