package cli

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"calcLoanAmount/internal/domain"
	"calcLoanAmount/internal/rate"
)

var (
	errInvalidHouseAmount     = errors.New("house amount must be greater than 0")
	errInvalidPrincipal       = errors.New("principal must be >= 0")
	errNoLoanNeeded           = errors.New("principal is greater than or equal to house amount")
	errInvalidFundAmount      = errors.New("fund loan amount must be >= 0")
	errInvalidFundYears       = errors.New("fund loan years must be greater than 0 when fund loan amount > 0")
	errInvalidFundRate        = errors.New("fund loan rate must be >= 0")
	errInvalidCommercialYears = errors.New("commercial loan years must be greater than 0")
	errInvalidCommercialRate  = errors.New("commercial loan rate must be >= 0")
	errInvalidMethod          = errors.New("repayment method must be epi or ep")
	errFundAmountTooLarge     = errors.New("fund loan amount exceeds required loan amount")
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
	method, err := parseMethod(s.Method)
	if err != nil {
		return Prepared{}, err
	}

	commercialAmount := s.HouseAmount - s.Principal - s.FundAmount
	if commercialAmount < 0 {
		return Prepared{}, errFundAmountTooLarge
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

	commercialLoan := domain.Loan{Method: method}
	var segments []rate.Segment
	if commercialAmount > 0 {
		if s.CommercialRatePercent < 0 {
			return Prepared{}, errInvalidCommercialRate
		}
		if s.CommercialYears <= 0 {
			return Prepared{}, errInvalidCommercialYears
		}

		changes, err := ParseChanges(s.CommercialChanges)
		if err != nil {
			return Prepared{}, err
		}
		segments, err = rate.BuildSegments(s.CommercialYears*12, s.CommercialRatePercent/100, changes)
		if err != nil {
			return Prepared{}, err
		}

		commercialLoan = domain.Loan{
			Principal:  commercialAmount,
			Years:      s.CommercialYears,
			AnnualRate: s.CommercialRatePercent / 100,
			Method:     method,
		}
	}

	return Prepared{
		FundLoan:           fundLoan,
		CommercialLoan:     commercialLoan,
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
	summary := fmt.Sprintf(
		"商业贷款金额: %.2f 万元\n公积金月供: %.2f 万元\n商业贷款月供: %.2f 万元\n组合贷款月供: %.2f 万元\n公积金总利息: %.2f 万元\n商业贷款总利息: %.2f 万元\n组合贷款总利息: %.2f 万元\n",
		report.CommercialAmount,
		report.Fund.MonthlyPayment,
		report.Commercial.MonthlyPayment,
		report.Combo.MonthlyPayment,
		report.Fund.TotalInterest,
		report.Commercial.TotalInterest,
		report.Combo.TotalInterest,
	)

	return summary + renderLoanLineCharts(report)
}

func renderLoanLineCharts(report Report) string {
	type chart struct {
		name     string
		schedule []domain.Payment
	}

	charts := make([]chart, 0, 2)
	if len(report.Fund.Schedule) > 0 {
		charts = append(charts, chart{
			name:     "公积金",
			schedule: report.Fund.Schedule,
		})
	}
	if len(report.Commercial.Schedule) > 0 {
		charts = append(charts, chart{
			name:     "商业贷款",
			schedule: report.Commercial.Schedule,
		})
	}
	if len(charts) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("\n")
	for i, c := range charts {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(renderLineChart(c.name, c.schedule))
	}

	return b.String()
}

func renderLineChart(name string, schedule []domain.Payment) string {
	const (
		chartHeight = 10
		maxPoints   = 48
	)

	points := samplePayments(schedule, maxPoints)
	if len(points) == 0 {
		return ""
	}

	maxValue := 0.0
	for _, v := range points {
		if v > maxValue {
			maxValue = v
		}
	}
	if maxValue <= 0 {
		maxValue = 1
	}

	cols := len(points)
	grid := make([][]rune, chartHeight)
	for i := 0; i < chartHeight; i++ {
		grid[i] = make([]rune, cols)
		for j := 0; j < cols; j++ {
			grid[i][j] = ' '
		}
	}

	toRow := func(value float64) int {
		norm := value / maxValue
		row := int(math.Round((1 - norm) * float64(chartHeight-1)))
		if row < 0 {
			row = 0
		}
		if row >= chartHeight {
			row = chartHeight - 1
		}
		return row
	}

	prevY := toRow(points[0])
	grid[prevY][0] = '*'
	for x := 1; x < cols; x++ {
		y := toRow(points[x])
		grid[y][x] = '*'

		if y != prevY {
			start := prevY
			end := y
			if start > end {
				start, end = end, start
			}
			for r := start; r <= end; r++ {
				grid[r][x] = '*'
			}
		}
		prevY = y
	}

	totalMonths := schedule[len(schedule)-1].Month

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s月供折线图(横轴:时间(月), 纵轴:钱(万元))\n", name))
	for r := 0; r < chartHeight; r++ {
		yValue := maxValue * (1 - float64(r)/float64(chartHeight-1))
		b.WriteString(fmt.Sprintf("%6.2f |%s\n", yValue, string(grid[r])))
	}
	b.WriteString(fmt.Sprintf("      +%s\n", strings.Repeat("-", cols)))

	startLabel := "1月"
	endLabel := fmt.Sprintf("%d月", totalMonths)
	if totalMonths > len(points) {
		endLabel = fmt.Sprintf("%d月(采样%d点)", totalMonths, len(points))
	}
	padding := cols - len(startLabel) - len(endLabel)
	if padding < 1 {
		padding = 1
	}
	b.WriteString(fmt.Sprintf("      %s%s%s\n", startLabel, strings.Repeat(" ", padding), endLabel))

	return b.String()
}

func samplePayments(schedule []domain.Payment, maxPoints int) []float64 {
	if len(schedule) == 0 {
		return nil
	}
	if maxPoints <= 0 {
		maxPoints = len(schedule)
	}
	if len(schedule) <= maxPoints {
		out := make([]float64, 0, len(schedule))
		for _, p := range schedule {
			out = append(out, p.Payment)
		}
		return out
	}

	out := make([]float64, 0, maxPoints)
	last := len(schedule) - 1
	for i := 0; i < maxPoints; i++ {
		idx := int(math.Round(float64(i) * float64(last) / float64(maxPoints-1)))
		out = append(out, schedule[idx].Payment)
	}
	return out
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
