package repository

// FineRepository holds the SQL for the fines table.
// It keeps no database handle: every method receives a Queryer.
type FineRepository struct{}

func NewFineRepository() *FineRepository {
	return &FineRepository{}
}

// HasUnpaidByUser tells if the user has an unpaid fine from ANY past loan.
func (r *FineRepository) HasUnpaidByUser(q Queryer, userID int) (bool, error) {
	return queryExists(q, `
		SELECT EXISTS (
			SELECT 1 FROM fines f
			JOIN loans l ON l.id = f.loan_id
			WHERE l.user_id = $1 AND f.status = 'unpaid')`, userID)
}