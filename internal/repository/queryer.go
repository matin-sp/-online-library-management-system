package repository

import "database/sql"

// Queryer is the small set of database methods the repositories need.
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