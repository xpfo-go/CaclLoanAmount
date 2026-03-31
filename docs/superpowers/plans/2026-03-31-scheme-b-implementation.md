# Scheme B Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rebuild mortgage calculation with correct monthly amortization, support variable repricing schedules, and keep CLI output deterministic and testable.

**Architecture:** Split core into `internal/domain` (loan models), `internal/calc` (amortization engine), and `internal/rate` (rate schedule/repricing). CLI in `main.go` only handles parsing, validation, orchestration, and display. Use table-driven tests to lock behavior for each iteration.

**Tech Stack:** Go 1.21, stdlib only (`testing`, `math`, `flag`, `fmt`).

---

### Task 1: Build fixed-rate core with TDD

**Files:**
- Create: `internal/domain/loan.go`
- Create: `internal/calc/fixed.go`
- Create: `internal/calc/fixed_test.go`

- [ ] **Step 1: Write failing tests for equal principal+interest and equal principal**
- [ ] **Step 2: Run tests and confirm failure due to missing implementation**
- [ ] **Step 3: Implement minimal monthly amortization logic**
- [ ] **Step 4: Run tests and confirm pass**

### Task 2: Add variable repricing schedule and combo-loan aggregation

**Files:**
- Create: `internal/rate/schedule.go`
- Create: `internal/rate/schedule_test.go`
- Create: `internal/calc/variable.go`
- Create: `internal/calc/variable_test.go`

- [ ] **Step 1: Write failing tests for repricing date/period and piecewise payment recompute**
- [ ] **Step 2: Run tests and confirm failure**
- [ ] **Step 3: Implement minimal schedule + variable calculator**
- [ ] **Step 4: Run tests and confirm pass**

### Task 3: Migrate CLI to new engine with input validation

**Files:**
- Modify: `main.go`
- Create: `internal/cli/input.go`
- Create: `internal/cli/input_test.go`

- [ ] **Step 1: Write failing tests for input validation and deterministic formatting**
- [ ] **Step 2: Run tests and confirm failure**
- [ ] **Step 3: Implement parser/validation and wire old prompts to new engine**
- [ ] **Step 4: Run tests and confirm pass**

### Task 4: Engineering baseline (CI + docs)

**Files:**
- Create: `.github/workflows/ci.yml`
- Create: `README.md`
- Create: `.gitignore`

- [ ] **Step 1: Add CI to run `go test ./...` + `go vet ./...` on push/PR**
- [ ] **Step 2: Add README with formulas, examples, and assumptions**
- [ ] **Step 3: Ignore build artifacts (`*.exe`, `bin/`)**
- [ ] **Step 4: Run full verification (`go test ./... && go vet ./...`)**
