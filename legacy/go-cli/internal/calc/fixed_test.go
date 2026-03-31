package calc_test

import (
	"math"
	"testing"

	"calcLoanAmount/internal/calc"
	"calcLoanAmount/internal/domain"
)

func TestCalculateFixed_EqualPrincipalAndInterest(t *testing.T) {
	loan := domain.Loan{
		Principal:  120000,
		Years:      1,
		AnnualRate: 0.12,
		Method:     domain.MethodEqualPrincipalInterest,
	}

	result, err := calc.CalculateFixed(loan)
	if err != nil {
		t.Fatalf("CalculateFixed returned error: %v", err)
	}

	if len(result.Schedule) != 12 {
		t.Fatalf("expected 12 months, got %d", len(result.Schedule))
	}

	if !almostEqual(result.MonthlyPayment, 10661.85, 0.01) {
		t.Fatalf("expected monthly payment 10661.85, got %.2f", result.MonthlyPayment)
	}

	if !almostEqual(result.TotalInterest, 7942.26, 0.02) {
		t.Fatalf("expected total interest 7942.26, got %.2f", result.TotalInterest)
	}

	first := result.Schedule[0]
	if !almostEqual(first.Interest, 1200, 0.01) {
		t.Fatalf("expected first month interest 1200, got %.2f", first.Interest)
	}
}

func TestCalculateFixed_EqualPrincipal(t *testing.T) {
	loan := domain.Loan{
		Principal:  120000,
		Years:      1,
		AnnualRate: 0.12,
		Method:     domain.MethodEqualPrincipal,
	}

	result, err := calc.CalculateFixed(loan)
	if err != nil {
		t.Fatalf("CalculateFixed returned error: %v", err)
	}

	if len(result.Schedule) != 12 {
		t.Fatalf("expected 12 months, got %d", len(result.Schedule))
	}

	first := result.Schedule[0]
	last := result.Schedule[len(result.Schedule)-1]

	if !almostEqual(first.Payment, 11200, 0.01) {
		t.Fatalf("expected first month payment 11200, got %.2f", first.Payment)
	}

	if !almostEqual(last.Payment, 10100, 0.01) {
		t.Fatalf("expected last month payment 10100, got %.2f", last.Payment)
	}

	if !almostEqual(result.TotalInterest, 7800, 0.01) {
		t.Fatalf("expected total interest 7800, got %.2f", result.TotalInterest)
	}
}

func TestCalculateFixed_InvalidLoan(t *testing.T) {
	_, err := calc.CalculateFixed(domain.Loan{})
	if err == nil {
		t.Fatalf("expected validation error for empty loan")
	}
}

func almostEqual(a float64, b float64, tol float64) bool {
	return math.Abs(a-b) <= tol
}
