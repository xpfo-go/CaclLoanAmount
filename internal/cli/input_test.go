package cli_test

import (
	"testing"

	"calcLoanAmount/internal/cli"
	"calcLoanAmount/internal/domain"
)

func TestBuildPrepared_ValidScenario(t *testing.T) {
	scenario := cli.Scenario{
		HouseAmount:           300,
		Principal:             100,
		FundAmount:            80,
		FundYears:             20,
		CommercialYears:       30,
		FundRatePercent:       2.6,
		CommercialRatePercent: 3.6,
		Method:                "epi",
		CommercialChanges:     "13:3.2,25:3.1",
	}

	prepared, err := cli.BuildPrepared(scenario)
	if err != nil {
		t.Fatalf("BuildPrepared returned error: %v", err)
	}

	if prepared.CommercialLoan.Principal != 120 {
		t.Fatalf("expected commercial principal 120, got %.2f", prepared.CommercialLoan.Principal)
	}
	if prepared.CommercialLoan.Method != domain.MethodEqualPrincipalInterest {
		t.Fatalf("unexpected method: %v", prepared.CommercialLoan.Method)
	}
	if len(prepared.CommercialSegments) != 3 {
		t.Fatalf("expected 3 commercial segments, got %d", len(prepared.CommercialSegments))
	}
	if prepared.CommercialSegments[0].StartMonth != 1 || prepared.CommercialSegments[0].EndMonth != 12 {
		t.Fatalf("unexpected segment[0]: %+v", prepared.CommercialSegments[0])
	}
}

func TestBuildPrepared_InvalidScenario(t *testing.T) {
	_, err := cli.BuildPrepared(cli.Scenario{HouseAmount: 100, Principal: 100})
	if err == nil {
		t.Fatalf("expected error when principal >= house amount")
	}
}

func TestParseChanges_InvalidFormat(t *testing.T) {
	_, err := cli.ParseChanges("13-3.2")
	if err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestFormatReport_Deterministic(t *testing.T) {
	report := cli.Report{
		CommercialAmount: 120,
		Fund: domain.LoanResult{
			MonthlyPayment: 1000,
			TotalInterest:  2000,
		},
		Commercial: domain.LoanResult{
			MonthlyPayment: 2000,
			TotalInterest:  3000,
		},
		Combo: domain.LoanResult{
			MonthlyPayment: 3000,
			TotalInterest:  5000,
		},
	}

	got := cli.FormatReport(report)
	want := "商业贷款金额: 120.00 万元\n" +
		"公积金月供: 1000.00 万元\n" +
		"商业贷款月供: 2000.00 万元\n" +
		"组合贷款月供: 3000.00 万元\n" +
		"公积金总利息: 2000.00 万元\n" +
		"商业贷款总利息: 3000.00 万元\n" +
		"组合贷款总利息: 5000.00 万元\n"
	if got != want {
		t.Fatalf("unexpected report output:\nwant:\n%s\ngot:\n%s", want, got)
	}
}
