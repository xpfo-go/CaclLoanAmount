package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"calcLoanAmount/internal/calc"
	"calcLoanAmount/internal/cli"
	"calcLoanAmount/internal/domain"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("房贷计算器（支持固定利率与重定价分段利率）")

	houseAmount, err := readFloat(reader, "请输入预计买房金额(单位:万，房子总价): ")
	if err != nil {
		exitWithError(err)
	}
	principal, err := readFloat(reader, "请输入您的本金(单位:万): ")
	if err != nil {
		exitWithError(err)
	}
	fundAmount, err := readFloat(reader, "请输入公积金贷款金额(单位:万，若无填0): ")
	if err != nil {
		exitWithError(err)
	}

	fundYears := 0
	if fundAmount > 0 {
		fundYears, err = readInt(reader, "请输入公积金贷款年份(单位:年): ")
		if err != nil {
			exitWithError(err)
		}
	}

	commercialYears, err := readInt(reader, "请输入商业贷款年份(单位:年): ")
	if err != nil {
		exitWithError(err)
	}

	fundRate := 0.0
	if fundAmount > 0 {
		fundRate, err = readFloat(reader, "请输入公积金年利率(%，示例2.6): ")
		if err != nil {
			exitWithError(err)
		}
	}

	commercialRate, err := readFloat(reader, "请输入商业贷款初始年利率(%，示例3.6): ")
	if err != nil {
		exitWithError(err)
	}
	method, err := readString(reader, "请输入还款方式(epi=等额本息, ep=等额本金): ")
	if err != nil {
		exitWithError(err)
	}
	changes, err := readString(reader, "请输入商业贷款重定价(格式: 月份:年利率%, 例如13:3.2,25:3.1；无则回车): ")
	if err != nil {
		exitWithError(err)
	}

	prepared, err := cli.BuildPrepared(cli.Scenario{
		HouseAmount:           houseAmount,
		Principal:             principal,
		FundAmount:            fundAmount,
		FundYears:             fundYears,
		CommercialYears:       commercialYears,
		FundRatePercent:       fundRate,
		CommercialRatePercent: commercialRate,
		Method:                method,
		CommercialChanges:     changes,
	})
	if err != nil {
		exitWithError(err)
	}

	fundResult := domain.LoanResult{}
	if prepared.FundLoan.Principal > 0 {
		fundResult, err = calc.CalculateFixed(prepared.FundLoan)
		if err != nil {
			exitWithError(err)
		}
	}

	commercialResult, err := calc.CalculateFixed(prepared.CommercialLoan)
	if err != nil {
		exitWithError(err)
	}
	if len(prepared.CommercialSegments) > 1 {
		commercialResult, err = calc.CalculateVariable(prepared.CommercialLoan, prepared.CommercialSegments)
		if err != nil {
			exitWithError(err)
		}
	}

	combo := calc.CombineResults(fundResult, commercialResult)
	report := cli.Report{
		CommercialAmount: prepared.CommercialAmount,
		Fund:             fundResult,
		Commercial:       commercialResult,
		Combo:            combo,
	}

	fmt.Println()
	fmt.Print(cli.FormatReport(report))
}

func readFloat(reader *bufio.Reader, prompt string) (float64, error) {
	for {
		text, err := readString(reader, prompt)
		if err != nil {
			return 0, err
		}
		value, err := strconv.ParseFloat(text, 64)
		if err == nil {
			return value, nil
		}
		fmt.Println("输入无效，请输入数字。")
	}
}

func readInt(reader *bufio.Reader, prompt string) (int, error) {
	for {
		text, err := readString(reader, prompt)
		if err != nil {
			return 0, err
		}
		value, err := strconv.Atoi(text)
		if err == nil {
			return value, nil
		}
		fmt.Println("输入无效，请输入整数。")
	}
}

func readString(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func exitWithError(err error) {
	fmt.Printf("错误: %v\n", err)
	os.Exit(1)
}
