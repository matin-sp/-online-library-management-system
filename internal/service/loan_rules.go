package service

// The standard limits. Everyone is a "standard" user until the membership
// payment exists. Change a number here and it changes everywhere.
const (
	standardMaxActiveLoans = 2  // books one user may hold at the same time
	standardLoanDays       = 14 // length of one loan
)

// borrowLimits are the limits that apply to one user.
type borrowLimits struct {
	maxActiveLoans int
	loanDays       int
}

// borrowLimitsFor returns the limits of a user. Today it is the same for
// everybody. When the member tier (5 loans, 30 days) is added, this is the
// ONLY function that has to change, because the user id is already passed in.
func borrowLimitsFor(userID int) borrowLimits {
	return borrowLimits{
		maxActiveLoans: standardMaxActiveLoans,
		loanDays:       standardLoanDays,
	}
}