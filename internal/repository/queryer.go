package repository

import "database/sql"

// Queryer is the small set of database methods the repositories need.
//
// Both *sql.DB (one statement on its own) and *sql.Tx (a statement inside a
// transaction) have these three methods. So a repository method that takes a
// Queryer works in both situations and does not know which one it got.
// Starting, committing and rolling back a transaction is the job of the
// service, never of a repository.
type Queryer interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

// queryExists runs a "SELECT EXISTS (...)" query and returns its true/false answer.
func queryExists(q Queryer, query string, args ...any) (bool, error) {
	var found bool
	if err := q.QueryRow(query, args...).Scan(&found); err != nil {
		return false, err
	}
	return found, nil
}