package calc

import (
	"errors"
	"math"

	"calcLoanAmount/internal/domain"
)

var (
	errInvalidPrincipal = errors.New("principal must be greater than 0")
	errInvalidYears     = errors.New("years must be greater than 0")
	errInvalidRate      = errors.New("annual rate must be >= 0")
	errInvalidMethod    = errors.New("unsupported repayment method")
)

func CalculateFixed(loan domain.Loan) (domain.LoanResult, error) {
	if loan.Principal <= 0 {
		return domain.LoanResult{}, errInvalidPrincipal
	}
	if loan.Years <= 0 {
		return domain.LoanResult{}, errInvalidYears
	}
	if loan.AnnualRate < 0 {
		return domain.LoanResult{}, errInvalidRate
	}

	switch loan.Method {
	case domain.MethodEqualPrincipalInterest:
		return calculateEqualPrincipalInterest(loan), nil
	case domain.MethodEqualPrincipal:
		return calculateEqualPrincipal(loan), nil
	default:
		return domain.LoanResult{}, errInvalidMethod
	}
}

func calculateEqualPrincipalInterest(loan domain.Loan) domain.LoanResult {
	months := loan.Years * 12
	monthlyRate := loan.AnnualRate / 12
	remaining := loan.Principal

	result := domain.LoanResult{
		Schedule: make([]domain.Payment, 0, months),
	}

	monthlyPayment := loan.Principal / float64(months)
	if monthlyRate > 0 {
		pow := math.Pow(1+monthlyRate, float64(months))
		monthlyPayment = loan.Principal * monthlyRate * pow / (pow - 1)
	}
	result.MonthlyPayment = monthlyPayment

	for month := 1; month <= months; month++ {
		interest := remaining * monthlyRate
		principal := monthlyPayment - interest
		payment := monthlyPayment

		if month == months {
			principal = remaining
			payment = principal + interest
		}

		remaining -= principal
		if remaining < 0 {
			remaining = 0
		}

		result.Schedule = append(result.Schedule, domain.Payment{
			Month:              month,
			Payment:            payment,
			Principal:          principal,
			Interest:           interest,
			RemainingPrincipal: remaining,
			AnnualRate:         loan.AnnualRate,
		})

		result.TotalPayment += payment
		result.TotalInterest += interest
	}

	return result
}

func calculateEqualPrincipal(loan domain.Loan) domain.LoanResult {
	months := loan.Years * 12
	monthlyRate := loan.AnnualRate / 12
	monthlyPrincipal := loan.Principal / float64(months)
	remaining := loan.Principal

	result := domain.LoanResult{
		Schedule: make([]domain.Payment, 0, months),
	}

	for month := 1; month <= months; month++ {
		interest := remaining * monthlyRate
		principal := monthlyPrincipal
		payment := principal + interest

		if month == months {
			principal = remaining
			payment = principal + interest
		}

		remaining -= principal
		if remaining < 0 {
			remaining = 0
		}

		result.Schedule = append(result.Schedule, domain.Payment{
			Month:              month,
			Payment:            payment,
			Principal:          principal,
			Interest:           interest,
			RemainingPrincipal: remaining,
			AnnualRate:         loan.AnnualRate,
		})

		result.TotalPayment += payment
		result.TotalInterest += interest
	}

	if len(result.Schedule) > 0 {
		result.MonthlyPayment = result.Schedule[0].Payment
	}

	return result
}
