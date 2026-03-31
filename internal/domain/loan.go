package domain

type RepaymentMethod string

const (
	MethodEqualPrincipalInterest RepaymentMethod = "equal_principal_interest"
	MethodEqualPrincipal         RepaymentMethod = "equal_principal"
)

type Loan struct {
	Principal  float64
	Years      int
	AnnualRate float64
	Method     RepaymentMethod
}

type Payment struct {
	Month              int
	Payment            float64
	Principal          float64
	Interest           float64
	RemainingPrincipal float64
	AnnualRate         float64
}

type LoanResult struct {
	MonthlyPayment float64
	TotalPayment   float64
	TotalInterest  float64
	Schedule       []Payment
}
