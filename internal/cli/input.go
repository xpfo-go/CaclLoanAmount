package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"calcLoanAmount/internal/domain"
	"calcLoanAmount/internal/rate"
)

var (
	errInvalidHouseAmount      = errors.New("house amount must be greater than 0")
	errInvalidPrincipal        = errors.New("principal must be >= 0")
	errNoLoanNeeded            = errors.New("principal is greater than or equal to house amount")
	errInvalidFundAmount       = errors.New("fund loan amount must be >= 0")
	errInvalidFundYears        = errors.New("fund loan years must be greater than 0 when fund loan amount > 0")
	errInvalidFundRate         = errors.New("fund loan rate must be >= 0")
	errInvalidCommercialYears  = errors.New("commercial loan years must be greater than 0")
	errInvalidCommercialRate   = errors.New("commercial loan rate must be >= 0")
	errInvalidMethod           = errors.New("repayment method must be epi or ep")
	errCommercialAmountNotNeed = errors.New("commercial loan amount must be greater than 0")
)

type Scenario struct {
	HouseAmount           float64
	Principal             float64
	FundAmount            float64
	FundYears             int
	CommercialYears       int
	FundRatePercent       float64
	CommercialRatePercent float64
	Method                string
	CommercialChanges     string
}

type Prepared struct {
	FundLoan           domain.Loan
	CommercialLoan     domain.Loan
	CommercialAmount   float64
	CommercialSegments []rate.Segment
}

type Report struct {
	CommercialAmount float64
	Fund             domain.LoanResult
	Commercial       domain.LoanResult
	Combo            domain.LoanResult
}

func BuildPrepared(s Scenario) (Prepared, error) {
	if s.HouseAmount <= 0 {
		return Prepared{}, errInvalidHouseAmount
	}
	if s.Principal < 0 {
		return Prepared{}, errInvalidPrincipal
	}
	if s.Principal >= s.HouseAmount {
		return Prepared{}, errNoLoanNeeded
	}
	if s.FundAmount < 0 {
		return Prepared{}, errInvalidFundAmount
	}
	if s.FundRatePercent < 0 {
		return Prepared{}, errInvalidFundRate
	}
	if s.CommercialRatePercent < 0 {
		return Prepared{}, errInvalidCommercialRate
	}
	if s.CommercialYears <= 0 {
		return Prepared{}, errInvalidCommercialYears
	}

	method, err := parseMethod(s.Method)
	if err != nil {
		return Prepared{}, err
	}

	commercialAmount := s.HouseAmount - s.Principal - s.FundAmount
	if commercialAmount <= 0 {
		return Prepared{}, errCommercialAmountNotNeed
	}

	changes, err := ParseChanges(s.CommercialChanges)
	if err != nil {
		return Prepared{}, err
	}
	segments, err := rate.BuildSegments(s.CommercialYears*12, s.CommercialRatePercent/100, changes)
	if err != nil {
		return Prepared{}, err
	}

	fundLoan := domain.Loan{Method: method}
	if s.FundAmount > 0 {
		if s.FundYears <= 0 {
			return Prepared{}, errInvalidFundYears
		}
		fundLoan = domain.Loan{
			Principal:  s.FundAmount,
			Years:      s.FundYears,
			AnnualRate: s.FundRatePercent / 100,
			Method:     method,
		}
	}

	return Prepared{
		FundLoan: fundLoan,
		CommercialLoan: domain.Loan{
			Principal:  commercialAmount,
			Years:      s.CommercialYears,
			AnnualRate: s.CommercialRatePercent / 100,
			Method:     method,
		},
		CommercialAmount:   commercialAmount,
		CommercialSegments: segments,
	}, nil
}

func ParseChanges(input string) ([]rate.Change, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, nil
	}

	parts := strings.Split(input, ",")
	changes := make([]rate.Change, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		pair := strings.Split(part, ":")
		if len(pair) != 2 {
			return nil, fmt.Errorf("invalid change format: %s", part)
		}

		month, err := strconv.Atoi(strings.TrimSpace(pair[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid change month: %s", pair[0])
		}
		ratePercent, err := strconv.ParseFloat(strings.TrimSpace(pair[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid annual rate: %s", pair[1])
		}

		changes = append(changes, rate.Change{
			Month:      month,
			AnnualRate: ratePercent / 100,
		})
	}

	return changes, nil
}

func FormatReport(report Report) string {
	return fmt.Sprintf(
		"商业贷款金额: %.2f 万元\n公积金月供: %.2f 万元\n商业贷款月供: %.2f 万元\n组合贷款月供: %.2f 万元\n公积金总利息: %.2f 万元\n商业贷款总利息: %.2f 万元\n组合贷款总利息: %.2f 万元\n",
		report.CommercialAmount,
		report.Fund.MonthlyPayment,
		report.Commercial.MonthlyPayment,
		report.Combo.MonthlyPayment,
		report.Fund.TotalInterest,
		report.Commercial.TotalInterest,
		report.Combo.TotalInterest,
	)
}

func parseMethod(raw string) (domain.RepaymentMethod, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "epi", "1":
		return domain.MethodEqualPrincipalInterest, nil
	case "ep", "2":
		return domain.MethodEqualPrincipal, nil
	default:
		return "", errInvalidMethod
	}
}
