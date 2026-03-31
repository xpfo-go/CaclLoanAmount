package calc_test

import (
	"testing"

	"calcLoanAmount/internal/calc"
	"calcLoanAmount/internal/domain"
	"calcLoanAmount/internal/rate"
)

func TestCalculateVariable_EqualPrincipalAndInterest(t *testing.T) {
	segments, err := rate.BuildSegments(12, 0.12, []rate.Change{
		{Month: 7, AnnualRate: 0.06},
	})
	if err != nil {
		t.Fatalf("BuildSegments returned error: %v", err)
	}

	loan := domain.Loan{
		Principal:  120000,
		Years:      1,
		AnnualRate: 0.12,
		Method:     domain.MethodEqualPrincipalInterest,
	}

	result, err := calc.CalculateVariable(loan, segments)
	if err != nil {
		t.Fatalf("CalculateVariable returned error: %v", err)
	}

	if len(result.Schedule) != 12 {
		t.Fatalf("expected 12 months, got %d", len(result.Schedule))
	}

	if !almostEqual(result.Schedule[0].Payment, 10661.85, 0.01) {
		t.Fatalf("expected month1 payment 10661.85, got %.2f", result.Schedule[0].Payment)
	}

	if !almostEqual(result.Schedule[6].Payment, 10479.39, 0.01) {
		t.Fatalf("expected month7 payment 10479.39, got %.2f", result.Schedule[6].Payment)
	}

	if !almostEqual(result.TotalInterest, 6847.48, 0.02) {
		t.Fatalf("expected total interest 6847.48, got %.2f", result.TotalInterest)
	}
}

func TestCombineResults(t *testing.T) {
	fundLoan := domain.Loan{
		Principal:  60000,
		Years:      1,
		AnnualRate: 0.06,
		Method:     domain.MethodEqualPrincipalInterest,
	}
	fundResult, err := calc.CalculateFixed(fundLoan)
	if err != nil {
		t.Fatalf("CalculateFixed returned error: %v", err)
	}

	segments, err := rate.BuildSegments(12, 0.12, []rate.Change{
		{Month: 7, AnnualRate: 0.06},
	})
	if err != nil {
		t.Fatalf("BuildSegments returned error: %v", err)
	}
	commercialLoan := domain.Loan{
		Principal:  60000,
		Years:      1,
		AnnualRate: 0.12,
		Method:     domain.MethodEqualPrincipalInterest,
	}
	commercialResult, err := calc.CalculateVariable(commercialLoan, segments)
	if err != nil {
		t.Fatalf("CalculateVariable returned error: %v", err)
	}

	combo := calc.CombineResults(fundResult, commercialResult)
	if !almostEqual(combo.TotalPayment, 125391.57, 0.02) {
		t.Fatalf("expected combo total payment 125391.57, got %.2f", combo.TotalPayment)
	}
	if !almostEqual(combo.TotalInterest, 5391.57, 0.02) {
		t.Fatalf("expected combo total interest 5391.57, got %.2f", combo.TotalInterest)
	}
	if len(combo.Schedule) != 12 {
		t.Fatalf("expected combo schedule length 12, got %d", len(combo.Schedule))
	}
	if !almostEqual(combo.Schedule[0].Payment, 10494.91, 0.02) {
		t.Fatalf("expected combo month1 payment 10494.91, got %.2f", combo.Schedule[0].Payment)
	}
}
