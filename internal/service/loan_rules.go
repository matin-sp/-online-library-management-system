package service

import "time"

const (
	standardMaxActiveLoans = 2 
	standardLoanDays       = 14 
	standardMaxRenewals    = 1 
	standardRenewalDays    = 3  
)


const finePerDay int64 = 10000

type borrowLimits struct {
	maxActiveLoans int
	loanDays       int
	maxRenewals    int
	renewalDays    int
}


func borrowLimitsFor(userID int) borrowLimits {
	return borrowLimits{
		maxActiveLoans: standardMaxActiveLoans,
		loanDays:       standardLoanDays,
		maxRenewals:    standardMaxRenewals,
		renewalDays:    standardRenewalDays,
	}
}


func lateDays(due, returned time.Time) int {
	late := returned.Sub(due)
	if late <= 0 {
		return 0
	}
	const day = 24 * time.Hour
	return int((late + day - 1) / day) 
}


func fineFor(days int) int64 {
	return int64(days) * finePerDay
}