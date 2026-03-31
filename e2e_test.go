package main_test

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestE2E_FixedRateFlow(t *testing.T) {
	input := "30\n10\n8\n1\n6\n1\n12\nepi\n\n"
	out, err := runCLI(input)
	if err != nil {
		t.Fatalf("expected success, got error: %v\noutput:\n%s", err, out)
	}

	assertContains(t, out, "商业贷款金额: 12.00 万元")
	assertContains(t, out, "公积金总利息: 0.26 万元")
	assertContains(t, out, "商业贷款总利息: 0.79 万元")
	assertContains(t, out, "组合贷款总利息: 1.06 万元")
	assertContains(t, out, "公积金月供折线图")
	assertContains(t, out, "商业贷款月供折线图")
}

func TestE2E_RepricingFlow(t *testing.T) {
	input := "30\n10\n8\n1\n6\n1\n12\nepi\n7:6\n"
	out, err := runCLI(input)
	if err != nil {
		t.Fatalf("expected success, got error: %v\noutput:\n%s", err, out)
	}

	assertContains(t, out, "商业贷款金额: 12.00 万元")
	assertContains(t, out, "商业贷款总利息: 0.68 万元")
	assertContains(t, out, "组合贷款总利息: 0.95 万元")
}

func TestE2E_InvalidInputFlow(t *testing.T) {
	input := "100\n100\n0\nepi\n"
	out, err := runCLI(input)
	if err == nil {
		t.Fatalf("expected error for invalid input, output:\n%s", out)
	}

	assertContains(t, out, "错误:")
}

func TestE2E_ZeroFundSkipsFundPrompts(t *testing.T) {
	input := "160\n60\n0\n30\n3.6\nepi\n\n"
	out, err := runCLI(input)
	if err != nil {
		t.Fatalf("expected success, got error: %v\noutput:\n%s", err, out)
	}

	assertContains(t, out, "商业贷款金额: 100.00 万元")
	assertContains(t, out, "公积金月供: 0.00 万元")
	assertContains(t, out, "商业贷款总利息: 63.67 万元")
	assertContains(t, out, "商业贷款月供折线图")
	assertNotContains(t, out, "公积金月供折线图")
}

func TestE2E_OnlyFundHasSingleChart(t *testing.T) {
	input := "160\n60\n100\n30\n2.6\nepi\n"
	out, err := runCLI(input)
	if err != nil {
		t.Fatalf("expected success, got error: %v\noutput:\n%s", err, out)
	}

	assertContains(t, out, "商业贷款金额: 0.00 万元")
	assertContains(t, out, "公积金月供折线图")
	assertNotContains(t, out, "商业贷款月供折线图")
}

func runCLI(input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", ".")
	cmd.Dir = "."
	cmd.Stdin = strings.NewReader(input)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func assertContains(t *testing.T, got string, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("expected output to contain %q, got:\n%s", want, got)
	}
}

func assertNotContains(t *testing.T, got string, notWant string) {
	t.Helper()
	if strings.Contains(got, notWant) {
		t.Fatalf("expected output not to contain %q, got:\n%s", notWant, got)
	}
}
