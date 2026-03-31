package calc

import (
	"errors"
	"math"

	"calcLoanAmount/internal/domain"
	"calcLoanAmount/internal/rate"
)

var errInvalidSegments = errors.New("invalid rate segments")

func CalculateVariable(loan domain.Loan, segments []rate.Segment) (domain.LoanResult, error) {
	if loan.Principal <= 0 {
		return domain.LoanResult{}, errInvalidPrincipal
	}
	if loan.Years <= 0 {
		return domain.LoanResult{}, errInvalidYears
	}
	if len(segments) == 0 {
		return domain.LoanResult{}, errInvalidSegments
	}

	totalMonths := loan.Years * 12
	if !validSegments(totalMonths, segments) {
		return domain.LoanResult{}, errInvalidSegments
	}

	switch loan.Method {
	case domain.MethodEqualPrincipalInterest:
		return calculateVariableEqualPrincipalInterest(loan.Principal, totalMonths, segments), nil
	case domain.MethodEqualPrincipal:
		return calculateVariableEqualPrincipal(loan.Principal, totalMonths, segments), nil
	default:
		return domain.LoanResult{}, errInvalidMethod
	}
}

func calculateVariableEqualPrincipalInterest(principal float64, totalMonths int, segments []rate.Segment) domain.LoanResult {
	result := domain.LoanResult{
		Schedule: make([]domain.Payment, 0, totalMonths),
	}
	remaining := principal
	month := 1

	for _, seg := range segments {
		remainingMonths := totalMonths - (month - 1)
		monthlyRate := seg.AnnualRate / 12
		monthlyPayment := remaining / float64(remainingMonths)
		if monthlyRate > 0 {
			pow := math.Pow(1+monthlyRate, float64(remainingMonths))
			monthlyPayment = remaining * monthlyRate * pow / (pow - 1)
		}

		for ; month <= seg.EndMonth; month++ {
			interest := remaining * monthlyRate
			principalPaid := monthlyPayment - interest
			payment := monthlyPayment
			if month == totalMonths {
				principalPaid = remaining
				payment = principalPaid + interest
			}

			remaining -= principalPaid
			if remaining < 0 {
				remaining = 0
			}

			result.Schedule = append(result.Schedule, domain.Payment{
				Month:              month,
				Payment:            payment,
				Principal:          principalPaid,
				Interest:           interest,
				RemainingPrincipal: remaining,
				AnnualRate:         seg.AnnualRate,
			})

			result.TotalPayment += payment
			result.TotalInterest += interest
		}
	}

	if len(result.Schedule) > 0 {
		result.MonthlyPayment = result.Schedule[0].Payment
	}

	return result
}

func calculateVariableEqualPrincipal(principal float64, totalMonths int, segments []rate.Segment) domain.LoanResult {
	result := domain.LoanResult{
		Schedule: make([]domain.Payment, 0, totalMonths),
	}
	monthlyPrincipal := principal / float64(totalMonths)
	remaining := principal
	month := 1

	for _, seg := range segments {
		monthlyRate := seg.AnnualRate / 12
		for ; month <= seg.EndMonth; month++ {
			interest := remaining * monthlyRate
			principalPaid := monthlyPrincipal
			payment := principalPaid + interest
			if month == totalMonths {
				principalPaid = remaining
				payment = principalPaid + interest
			}

			remaining -= principalPaid
			if remaining < 0 {
				remaining = 0
			}

			result.Schedule = append(result.Schedule, domain.Payment{
				Month:              month,
				Payment:            payment,
				Principal:          principalPaid,
				Interest:           interest,
				RemainingPrincipal: remaining,
				AnnualRate:         seg.AnnualRate,
			})

			result.TotalPayment += payment
			result.TotalInterest += interest
		}
	}

	if len(result.Schedule) > 0 {
		result.MonthlyPayment = result.Schedule[0].Payment
	}

	return result
}

func CombineResults(results ...domain.LoanResult) domain.LoanResult {
	out := domain.LoanResult{}
	maxMonths := 0
	for _, result := range results {
		if len(result.Schedule) > maxMonths {
			maxMonths = len(result.Schedule)
		}
		out.TotalPayment += result.TotalPayment
		out.TotalInterest += result.TotalInterest
	}

	out.Schedule = make([]domain.Payment, maxMonths)
	for month := 0; month < maxMonths; month++ {
		item := domain.Payment{Month: month + 1}
		for _, result := range results {
			if month >= len(result.Schedule) {
				continue
			}
			monthItem := result.Schedule[month]
			item.Payment += monthItem.Payment
			item.Principal += monthItem.Principal
			item.Interest += monthItem.Interest
			item.RemainingPrincipal += monthItem.RemainingPrincipal
		}
		out.Schedule[month] = item
	}

	if len(out.Schedule) > 0 {
		out.MonthlyPayment = out.Schedule[0].Payment
	}

	return out
}

func validSegments(totalMonths int, segments []rate.Segment) bool {
	if len(segments) == 0 {
		return false
	}
	if segments[0].StartMonth != 1 {
		return false
	}

	expectedStart := 1
	for _, seg := range segments {
		if seg.StartMonth != expectedStart {
			return false
		}
		if seg.EndMonth < seg.StartMonth || seg.EndMonth > totalMonths {
			return false
		}
		if seg.AnnualRate < 0 {
			return false
		}
		expectedStart = seg.EndMonth + 1
	}

	return expectedStart == totalMonths+1
}
