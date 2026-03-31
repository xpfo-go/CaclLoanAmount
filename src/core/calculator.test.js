import { describe, expect, it } from 'vitest'
import {
  buildSegments,
  calculateReport,
  calculateFixed,
  calculateVariable,
  combineResults,
  parseRepricing,
  prepareScenario,
} from './calculator'

describe('calculator core', () => {
  it('calculates fixed equal principal and interest', () => {
    const result = calculateFixed({
      principal: 120000,
      years: 1,
      annualRate: 0.12,
      method: 'epi',
    })

    expect(result.schedule).toHaveLength(12)
    expect(result.monthlyPayment).toBeCloseTo(10661.85, 2)
    expect(result.totalInterest).toBeCloseTo(7942.26, 2)
    expect(result.schedule[0].interest).toBeCloseTo(1200, 2)
  })

  it('calculates fixed equal principal', () => {
    const result = calculateFixed({
      principal: 120000,
      years: 1,
      annualRate: 0.12,
      method: 'ep',
    })

    expect(result.schedule[0].payment).toBeCloseTo(11200, 2)
    expect(result.schedule.at(-1).payment).toBeCloseTo(10100, 2)
    expect(result.totalInterest).toBeCloseTo(7800, 2)
  })

  it('parses repricing changes and builds segments', () => {
    const changes = parseRepricing('7:6,13:5.5')
    const segments = buildSegments(24, 0.12, changes)

    expect(segments).toEqual([
      { startMonth: 1, endMonth: 6, annualRate: 0.12 },
      { startMonth: 7, endMonth: 12, annualRate: 0.06 },
      { startMonth: 13, endMonth: 24, annualRate: 0.055 },
    ])
  })

  it('calculates variable equal principal and interest with repricing', () => {
    const segments = buildSegments(12, 0.12, [{ month: 7, annualRate: 0.06 }])
    const result = calculateVariable(
      {
        principal: 120000,
        years: 1,
        annualRate: 0.12,
        method: 'epi',
      },
      segments,
    )

    expect(result.schedule).toHaveLength(12)
    expect(result.schedule[0].payment).toBeCloseTo(10661.85, 2)
    expect(result.schedule[6].payment).toBeCloseTo(10479.39, 2)
    expect(result.totalInterest).toBeCloseTo(6847.48, 2)
  })

  it('combines loan results by month', () => {
    const fund = calculateFixed({ principal: 60000, years: 1, annualRate: 0.06, method: 'epi' })
    const commercial = calculateVariable(
      { principal: 60000, years: 1, annualRate: 0.12, method: 'epi' },
      buildSegments(12, 0.12, [{ month: 7, annualRate: 0.06 }]),
    )
    const combo = combineResults(fund, commercial)

    expect(combo.totalPayment).toBeCloseTo(125391.57, 2)
    expect(combo.totalInterest).toBeCloseTo(5391.57, 2)
    expect(combo.schedule).toHaveLength(12)
  })

  it('prepares scenario and allows only fund loan', () => {
    const scenario = prepareScenario({
      houseAmount: 160,
      principal: 60,
      fundAmount: 100,
      fundYears: 30,
      fundRatePercent: 2.6,
      commercialYears: 0,
      commercialRatePercent: 0,
      method: 'epi',
      commercialChangesText: '',
    })

    expect(scenario.commercialAmount).toBe(0)
    expect(scenario.commercialLoan.principal).toBe(0)
    expect(scenario.commercialSegments).toHaveLength(0)
  })

  it('keeps old fixed-rate combo output in scenario flow', () => {
    const prepared = prepareScenario({
      houseAmount: 30,
      principal: 10,
      fundAmount: 8,
      fundYears: 1,
      fundRatePercent: 6,
      commercialYears: 1,
      commercialRatePercent: 12,
      method: 'epi',
      commercialChangesText: '',
    })
    const report = calculateReport(prepared)

    expect(report.commercialAmount).toBeCloseTo(12, 2)
    expect(report.fund.totalInterest).toBeCloseTo(0.26, 2)
    expect(report.commercial.totalInterest).toBeCloseTo(0.79, 2)
    expect(report.combo.totalInterest).toBeCloseTo(1.06, 2)
  })

  it('supports repricing in scenario flow', () => {
    const prepared = prepareScenario({
      houseAmount: 30,
      principal: 10,
      fundAmount: 8,
      fundYears: 1,
      fundRatePercent: 6,
      commercialYears: 1,
      commercialRatePercent: 12,
      method: 'epi',
      commercialChangesText: '7:6',
    })
    const report = calculateReport(prepared)

    expect(report.commercial.totalInterest).toBeCloseTo(0.68, 2)
    expect(report.combo.totalInterest).toBeCloseTo(0.95, 2)
  })
})
