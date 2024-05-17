package main

import (
	"fmt"
	"math"
	"os"
)

const (
	// 首套住房5年以内含5年公积金贷款利率
	FIRST_HOUSE_ACCUMULATION_FUND_LOAN_RATE_WITHIN_FIVE_YEARS = 2.35
	// 首套住房5年以上公积金贷款利率
	FIRST_HOUSE_ACCUMULATION_FUND_LOAN_RATE_MORE_FIVE_YEARS = 2.85
)

const (
	// 二套房5年以内公积金贷款利率
	NOT_FIRST_HOUSE_ACCUMULATION_FUND_LOAN_RATE_WITHIN_FIVE_YEARS = 2.775
	// 二套房5年以上公积金贷款利率
	NOT_FIRST_HOUSE_ACCUMULATION_FUND_LOAN_RATE_MORE_FIVE_YEARS = 3.325
)

var (
	// 5年以内商业贷款利率
	COMMERCIAL_LOAN_RATE_WITHIN_FIVE_YEARS = 3.45
	// 5年以上商业贷款利率
	COMMERCIAL_LOAN_RATE_MORE_FIVE_YEARS = 3.95
)

type LoanAmount struct {
	// 公积金贷款还款金额
	AccumulationFundRepaymentAmount float64
	// 公积金贷款每年还款金额
	AccumulationFundAnnualRepaymentAmount float64
	// 商业贷款还款金额
	CommercialRepaymentAmount float64
	// 商业贷款每年还款金额
	CommercialAnnualRepaymentAmount float64
	// 公积金贷款超出金额
	AccumulationFundLoanExceedAmount float64
	// 商业贷款超出金额
	CommercialLoanExceedAmount float64
	// 总共贷款超出金额
	TotalExceedLoanAmount float64
}

func (l *LoanAmount) String() {
	fmt.Println("公积金贷款还款金额：", Decimal(l.AccumulationFundRepaymentAmount), "万元")
	fmt.Println("公积金贷款每年还款金额：", Decimal(l.AccumulationFundAnnualRepaymentAmount), "万元")
	fmt.Println("商业贷款还款金额：", Decimal(l.CommercialRepaymentAmount), "万元")
	fmt.Println("商业贷款每年还款金额：", Decimal(l.CommercialAnnualRepaymentAmount), "万元")
	fmt.Println("公积金贷款超出金额：", Decimal(l.AccumulationFundLoanExceedAmount), "万元")
	fmt.Println("商业贷款超出金额：", Decimal(l.CommercialLoanExceedAmount), "万元")
	fmt.Println("总共贷款超出金额：", Decimal(l.TotalExceedLoanAmount), "万元")
}

func main() {
	var (
		houseAmount                float64
		principal                  float64
		accumulationFundLoanAmount float64
		commercialLoanAmount       float64
		accumulationFundLoanYears  int
		commercialLoanAmountYears  int
		isFirstHouse               int

		commercialLoanRateLte5 float64
		commercialLoanRateGt5  float64
	)

	fmt.Println("请输入预计买房金额(单位:万，房子总价)：")
	fmt.Scanln(&houseAmount)

	fmt.Println("请输入您的本金(单位:万)：")
	fmt.Scanln(&principal)

	if principal >= houseAmount {
		fmt.Println("全款买房，崇拜大佬！")

		fmt.Printf("Press Enter to exit ...")
		endKey := make([]byte, 1)
		os.Stdin.Read(endKey)
		return
	}

	fmt.Println("请输入公积金贷款金额(单位:万)：")
	fmt.Scanln(&accumulationFundLoanAmount)

	commercialLoanAmount = houseAmount - principal - accumulationFundLoanAmount
	if commercialLoanAmount <= 0 {
		fmt.Println("您的本金+公积金贷款已经够买房啦，恭喜大佬！")
		fmt.Printf("Press Enter to exit ...")
		endKey := make([]byte, 1)
		os.Stdin.Read(endKey)
		return
	} else {
		fmt.Println("您的本金+公积金贷款还不够买房哦，需要商业贷款金额为(单位:万)：", commercialLoanAmount)
	}

	fmt.Println("请输入公积金贷款年份(单位:年)：")
	fmt.Scanln(&accumulationFundLoanYears)

	fmt.Println("请输入商业贷款年份(单位:年)：")
	fmt.Scanln(&commercialLoanAmountYears)

	if commercialLoanAmountYears <= 5 {
		fmt.Println("请输入商业贷款5年以内利率(示例: 3.45)，默认3.45，因为2024.05.18之前最低就是3.45：")
		fmt.Scanln(&commercialLoanRateLte5)
		if commercialLoanRateLte5 > 0 {
			COMMERCIAL_LOAN_RATE_WITHIN_FIVE_YEARS = commercialLoanRateLte5
		}
	} else {
		fmt.Println("请输入商业贷款5年以上利率(示例: 3.95)，默认3.95，因为2024.05.18之前最低就是3.95：")
		fmt.Scanln(&commercialLoanRateGt5)
		if commercialLoanRateGt5 > 0 {
			COMMERCIAL_LOAN_RATE_MORE_FIVE_YEARS = commercialLoanRateGt5
		}
	}

	fmt.Println("请输入是否首房(0:否，1：是)：")
	fmt.Scanln(&isFirstHouse)

	fmt.Println("")
	fmt.Println("------------------")
	amount := calculateLoanAmount(
		accumulationFundLoanAmount, accumulationFundLoanYears,
		commercialLoanAmount, commercialLoanAmountYears, isFirstHouse)
	amount.String()

	fmt.Printf("Press Enter to exit ...")
	endKey := make([]byte, 1)
	os.Stdin.Read(endKey)
}

func calculateLoanAmount(
	accumulationFundLoanAmount float64, accumulationFundLoanYears int,
	commercialLoanAmount float64, commercialLoanAmountYears int,
	isFirstHouse int) LoanAmount {
	loanAmount := LoanAmount{}

	// 公积金贷款还款金额
	var accumulationFundRepaymentAmount float64
	if isFirstHouse == 1 {
		if accumulationFundLoanYears <= 5 {
			accumulationFundRepaymentAmount = accumulationFundLoanAmount * math.Pow((100+FIRST_HOUSE_ACCUMULATION_FUND_LOAN_RATE_WITHIN_FIVE_YEARS)/100, float64(accumulationFundLoanYears))
		} else {
			accumulationFundRepaymentAmount = accumulationFundLoanAmount * math.Pow((100+FIRST_HOUSE_ACCUMULATION_FUND_LOAN_RATE_MORE_FIVE_YEARS)/100, float64(accumulationFundLoanYears))
		}
	} else {
		if accumulationFundLoanYears <= 5 {
			accumulationFundRepaymentAmount = accumulationFundLoanAmount * math.Pow((100+NOT_FIRST_HOUSE_ACCUMULATION_FUND_LOAN_RATE_WITHIN_FIVE_YEARS)/100, float64(accumulationFundLoanYears))
		} else {
			accumulationFundRepaymentAmount = accumulationFundLoanAmount * math.Pow((100+NOT_FIRST_HOUSE_ACCUMULATION_FUND_LOAN_RATE_MORE_FIVE_YEARS)/100, float64(accumulationFundLoanYears))
		}
	}
	loanAmount.AccumulationFundRepaymentAmount = accumulationFundRepaymentAmount
	//fmt.Println("公积金贷款还款金额：", accumulationFundRepaymentAmount)

	// 公积金贷款每年还款金额
	loanAmount.AccumulationFundAnnualRepaymentAmount = accumulationFundRepaymentAmount / float64(accumulationFundLoanYears)
	//fmt.Println("公积金贷款每年还款金额：", accumulationFundRepaymentAmount/float64(accumulationFundLoanYears))

	// 商业贷款还款金额
	var commercialRepaymentAmount float64
	if commercialLoanAmountYears <= 5 {
		commercialRepaymentAmount = commercialLoanAmount * math.Pow((100+COMMERCIAL_LOAN_RATE_WITHIN_FIVE_YEARS)/100, float64(commercialLoanAmountYears))
	} else {
		commercialRepaymentAmount = commercialLoanAmount * math.Pow((100+COMMERCIAL_LOAN_RATE_MORE_FIVE_YEARS)/100, float64(commercialLoanAmountYears))
	}
	loanAmount.CommercialRepaymentAmount = commercialRepaymentAmount
	//fmt.Println("商业贷款还款金额：", commercialRepaymentAmount)

	// 商业贷款每年还款金额
	loanAmount.CommercialAnnualRepaymentAmount = commercialRepaymentAmount / float64(commercialLoanAmountYears)
	//fmt.Println("商业贷款每年还款金额：", commercialRepaymentAmount/float64(commercialLoanAmountYears))

	// 公积金贷款超出金额
	loanAmount.AccumulationFundLoanExceedAmount = accumulationFundRepaymentAmount - accumulationFundLoanAmount
	//fmt.Println("公积金贷款超出金额：", accumulationFundRepaymentAmount-accumulationFundLoanAmount)

	// 商业贷款超出金额
	loanAmount.CommercialLoanExceedAmount = commercialRepaymentAmount - commercialLoanAmount
	//fmt.Println("商业贷款超出金额：", commercialRepaymentAmount-commercialLoanAmount)

	// 总共贷款超出金额
	loanAmount.TotalExceedLoanAmount = accumulationFundRepaymentAmount - accumulationFundLoanAmount + commercialRepaymentAmount - commercialLoanAmount
	//fmt.Println("总共贷款超出金额：", accumulationFundRepaymentAmount-accumulationFundLoanAmount+commercialRepaymentAmount-commercialLoanAmount)

	return loanAmount
}

func Decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}
