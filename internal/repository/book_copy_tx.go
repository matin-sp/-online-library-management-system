package repository

import (
	"database/sql"
	"errors"

	"github.com/matin-sp/-online-library-management-system/internal/model"
)

// LockByID reads a copy and locks its row until the surrounding transaction
// ends. If another transaction is lending the same copy right now, this call
// waits until that transaction is finished and then returns the NEW status.
// That is why the status check in the service can be trusted.
func (r *BookCopyRepository) LockByID(q Queryer, id int) (*model.BookCopy, error) {
	var c model.BookCopy
	err := q.QueryRow(
		"SELECT "+copyColumns+" FROM book_copies WHERE id = $1 FOR NO KEY UPDATE", id,
	).Scan(&c.ID, &c.BookID, &c.CopyNumber, &c.Status, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCopyNotFound
		}
		return nil, err
	}
	return &c, nil
}

// SetStatus changes the status of a copy. It has no rules of its own: the
// caller (the service) already decided that the change is allowed.
// UpdateStatus (the librarian endpoint) is a different method and is untouched.
func (r *BookCopyRepository) SetStatus(q Queryer, id int, status string) error {
	result, err := q.Exec("UPDATE book_copies SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return ErrCopyNotFound
	}
	return nil
}