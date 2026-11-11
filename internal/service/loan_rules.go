package service

import "time"

// The standard limits. Everyone is a "standard" user until the membership
// payment exists. Change a number here and it changes everywhere.
const (
	standardMaxActiveLoans = 2  // books one user may hold at the same time
	standardLoanDays       = 14 // length of one loan
)

// finePerDay is the fine for ONE late day, in rial.
// PLACEHOLDER: the real amount is not decided yet (open item 1 in the rules
// document). Changing this one number is all that is needed.
const finePerDay int64 = 10000

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

// lateDays says how many days late a return is. A day that has started counts
// as a full day: 1 second late is 1 day, exactly 24 hours late is 1 day,
// 24 hours and 1 second late is 2 days. On time or early is 0.
func lateDays(due, returned time.Time) int {
	late := returned.Sub(due)
	if late <= 0 {
		return 0
	}
	const day = 24 * time.Hour
	return int((late + day - 1) / day) // division that rounds up
}

// fineFor is the amount (in rial) for a number of late days.
func fineFor(days int) int64 {
	return int64(days) * finePerDay
}