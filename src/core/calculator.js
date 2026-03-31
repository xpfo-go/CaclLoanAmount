const METHOD_EPI = 'epi'
const METHOD_EP = 'ep'

const EMPTY_RESULT = {
  monthlyPayment: 0,
  totalPayment: 0,
  totalInterest: 0,
  schedule: [],
}

export function parseRepricing(input) {
  const text = String(input ?? '').trim()
  if (!text) return []

  return text.split(',').map((part) => {
    const [monthRaw, rateRaw] = part.trim().split(':')
    if (!monthRaw || !rateRaw) {
      throw new Error(`重定价格式错误: ${part}`)
    }

    const month = Number.parseInt(monthRaw.trim(), 10)
    const annualRate = Number.parseFloat(rateRaw.trim()) / 100

    if (!Number.isFinite(month) || month <= 0) {
      throw new Error(`重定价月份非法: ${monthRaw}`)
    }
    if (!Number.isFinite(annualRate) || annualRate < 0) {
      throw new Error(`重定价利率非法: ${rateRaw}`)
    }

    return { month, annualRate }
  })
}

export function buildSegments(totalMonths, initialRate, changes = []) {
  if (!Number.isInteger(totalMonths) || totalMonths <= 0) {
    throw new Error('还款总月数必须大于 0')
  }
  if (!Number.isFinite(initialRate) || initialRate < 0) {
    throw new Error('初始年利率非法')
  }

  const sorted = [...changes].sort((a, b) => a.month - b.month)
  const seen = new Set()
  let currentStart = 1
  let currentRate = initialRate
  const segments = []

  for (const change of sorted) {
    if (!Number.isInteger(change.month) || change.month < 1 || change.month > totalMonths) {
      throw new Error(`重定价月份超出范围: ${change.month}`)
    }
    if (!Number.isFinite(change.annualRate) || change.annualRate < 0) {
      throw new Error(`重定价利率非法: ${change.annualRate}`)
    }
    if (seen.has(change.month)) {
      throw new Error(`重定价月份重复: ${change.month}`)
    }
    seen.add(change.month)

    if (change.month === currentStart) {
      currentRate = change.annualRate
      continue
    }

    segments.push({
      startMonth: currentStart,
      endMonth: change.month - 1,
      annualRate: currentRate,
    })

    currentStart = change.month
    currentRate = change.annualRate
  }

  segments.push({
    startMonth: currentStart,
    endMonth: totalMonths,
    annualRate: currentRate,
  })

  return segments
}

export function prepareScenario(input) {
  const houseAmount = toNumber(input.houseAmount)
  const principal = toNumber(input.principal)
  const fundAmount = toNumber(input.fundAmount)
  const method = normalizeMethod(input.method)

  if (houseAmount <= 0) throw new Error('房屋总价必须大于 0')
  if (principal < 0) throw new Error('本金不能小于 0')
  if (principal >= houseAmount) throw new Error('本金已经覆盖房价，不需要贷款')
  if (fundAmount < 0) throw new Error('公积金贷款金额不能小于 0')

  const commercialAmount = houseAmount - principal - fundAmount
  if (commercialAmount < 0) throw new Error('公积金贷款金额超过所需贷款')

  const fundLoan = {
    principal: 0,
    years: 0,
    annualRate: 0,
    method,
  }

  if (fundAmount > 0) {
    fundLoan.principal = fundAmount
    fundLoan.years = toInteger(input.fundYears)
    fundLoan.annualRate = toNumber(input.fundRatePercent) / 100

    if (fundLoan.years <= 0) throw new Error('公积金贷款年限必须大于 0')
    if (fundLoan.annualRate < 0) throw new Error('公积金贷款利率不能小于 0')
  }

  const commercialLoan = {
    principal: 0,
    years: 0,
    annualRate: 0,
    method,
  }
  let commercialSegments = []

  if (commercialAmount > 0) {
    commercialLoan.principal = commercialAmount
    commercialLoan.years = toInteger(input.commercialYears)
    commercialLoan.annualRate = toNumber(input.commercialRatePercent) / 100

    if (commercialLoan.years <= 0) throw new Error('商业贷款年限必须大于 0')
    if (commercialLoan.annualRate < 0) throw new Error('商业贷款利率不能小于 0')

    const changes = parseRepricing(input.commercialChangesText)
    commercialSegments = buildSegments(commercialLoan.years * 12, commercialLoan.annualRate, changes)
  }

  return {
    method,
    commercialAmount,
    fundLoan,
    commercialLoan,
    commercialSegments,
  }
}

export function calculateReport(prepared) {
  const fund = prepared.fundLoan.principal > 0 ? calculateFixed(prepared.fundLoan) : emptyResult()

  let commercial = EMPTY_RESULT
  if (prepared.commercialLoan.principal > 0) {
    commercial =
      prepared.commercialSegments.length > 1
        ? calculateVariable(prepared.commercialLoan, prepared.commercialSegments)
        : calculateFixed(prepared.commercialLoan)
  }

  const combo = combineResults(fund, commercial)

  return {
    commercialAmount: prepared.commercialAmount,
    fund,
    commercial,
    combo,
  }
}

export function calculateFixed(loan) {
  validateLoan(loan)

  if (loan.method === METHOD_EPI) {
    return calculateFixedEPI(loan)
  }
  if (loan.method === METHOD_EP) {
    return calculateFixedEP(loan)
  }
  throw new Error('不支持的还款方式')
}

export function calculateVariable(loan, segments) {
  validateLoan(loan)
  validateSegments(loan.years * 12, segments)

  if (loan.method === METHOD_EPI) {
    return calculateVariableEPI(loan, segments)
  }
  if (loan.method === METHOD_EP) {
    return calculateVariableEP(loan, segments)
  }

  throw new Error('不支持的还款方式')
}

export function combineResults(...results) {
  const maxMonths = results.reduce((m, item) => Math.max(m, item.schedule.length), 0)
  const schedule = []

  for (let i = 0; i < maxMonths; i += 1) {
    const item = {
      month: i + 1,
      payment: 0,
      principal: 0,
      interest: 0,
      remainingPrincipal: 0,
    }

    for (const result of results) {
      const monthItem = result.schedule[i]
      if (!monthItem) continue
      item.payment += monthItem.payment
      item.principal += monthItem.principal
      item.interest += monthItem.interest
      item.remainingPrincipal += monthItem.remainingPrincipal
    }

    schedule.push(item)
  }

  const totalPayment = results.reduce((sum, item) => sum + item.totalPayment, 0)
  const totalInterest = results.reduce((sum, item) => sum + item.totalInterest, 0)

  return {
    monthlyPayment: schedule[0]?.payment ?? 0,
    totalPayment,
    totalInterest,
    schedule,
  }
}

function calculateFixedEPI(loan) {
  const months = loan.years * 12
  const monthlyRate = loan.annualRate / 12
  const schedule = []

  let remaining = loan.principal
  const monthlyPayment =
    monthlyRate === 0
      ? loan.principal / months
      : (loan.principal * monthlyRate * (1 + monthlyRate) ** months) / ((1 + monthlyRate) ** months - 1)

  let totalPayment = 0
  let totalInterest = 0

  for (let month = 1; month <= months; month += 1) {
    const interest = remaining * monthlyRate
    let principal = monthlyPayment - interest
    let payment = monthlyPayment

    if (month === months) {
      principal = remaining
      payment = principal + interest
    }

    remaining -= principal
    if (remaining < 0) remaining = 0

    schedule.push({
      month,
      payment,
      principal,
      interest,
      remainingPrincipal: remaining,
      annualRate: loan.annualRate,
    })

    totalPayment += payment
    totalInterest += interest
  }

  return {
    monthlyPayment,
    totalPayment,
    totalInterest,
    schedule,
  }
}

function calculateFixedEP(loan) {
  const months = loan.years * 12
  const monthlyRate = loan.annualRate / 12
  const monthlyPrincipal = loan.principal / months
  const schedule = []

  let remaining = loan.principal
  let totalPayment = 0
  let totalInterest = 0

  for (let month = 1; month <= months; month += 1) {
    const interest = remaining * monthlyRate
    let principal = monthlyPrincipal
    let payment = principal + interest

    if (month === months) {
      principal = remaining
      payment = principal + interest
    }

    remaining -= principal
    if (remaining < 0) remaining = 0

    schedule.push({
      month,
      payment,
      principal,
      interest,
      remainingPrincipal: remaining,
      annualRate: loan.annualRate,
    })

    totalPayment += payment
    totalInterest += interest
  }

  return {
    monthlyPayment: schedule[0]?.payment ?? 0,
    totalPayment,
    totalInterest,
    schedule,
  }
}

function calculateVariableEPI(loan, segments) {
  const totalMonths = loan.years * 12
  const schedule = []
  let remaining = loan.principal
  let month = 1
  let totalPayment = 0
  let totalInterest = 0

  for (const segment of segments) {
    const remainingMonths = totalMonths - (month - 1)
    const monthlyRate = segment.annualRate / 12
    const monthlyPayment =
      monthlyRate === 0
        ? remaining / remainingMonths
        : (remaining * monthlyRate * (1 + monthlyRate) ** remainingMonths) /
          ((1 + monthlyRate) ** remainingMonths - 1)

    for (; month <= segment.endMonth; month += 1) {
      const interest = remaining * monthlyRate
      let principal = monthlyPayment - interest
      let payment = monthlyPayment

      if (month === totalMonths) {
        principal = remaining
        payment = principal + interest
      }

      remaining -= principal
      if (remaining < 0) remaining = 0

      schedule.push({
        month,
        payment,
        principal,
        interest,
        remainingPrincipal: remaining,
        annualRate: segment.annualRate,
      })

      totalPayment += payment
      totalInterest += interest
    }
  }

  return {
    monthlyPayment: schedule[0]?.payment ?? 0,
    totalPayment,
    totalInterest,
    schedule,
  }
}

function calculateVariableEP(loan, segments) {
  const totalMonths = loan.years * 12
  const monthlyPrincipal = loan.principal / totalMonths
  const schedule = []

  let remaining = loan.principal
  let month = 1
  let totalPayment = 0
  let totalInterest = 0

  for (const segment of segments) {
    const monthlyRate = segment.annualRate / 12

    for (; month <= segment.endMonth; month += 1) {
      const interest = remaining * monthlyRate
      let principal = monthlyPrincipal
      let payment = principal + interest

      if (month === totalMonths) {
        principal = remaining
        payment = principal + interest
      }

      remaining -= principal
      if (remaining < 0) remaining = 0

      schedule.push({
        month,
        payment,
        principal,
        interest,
        remainingPrincipal: remaining,
        annualRate: segment.annualRate,
      })

      totalPayment += payment
      totalInterest += interest
    }
  }

  return {
    monthlyPayment: schedule[0]?.payment ?? 0,
    totalPayment,
    totalInterest,
    schedule,
  }
}

function validateLoan(loan) {
  if (!loan || typeof loan !== 'object') throw new Error('贷款参数不能为空')
  if (!Number.isFinite(loan.principal) || loan.principal <= 0) throw new Error('贷款本金必须大于 0')
  if (!Number.isInteger(loan.years) || loan.years <= 0) throw new Error('贷款年限必须大于 0')
  if (!Number.isFinite(loan.annualRate) || loan.annualRate < 0) throw new Error('贷款利率不能小于 0')
  if (![METHOD_EPI, METHOD_EP].includes(normalizeMethod(loan.method))) throw new Error('还款方式非法')
}

function validateSegments(totalMonths, segments) {
  if (!Array.isArray(segments) || segments.length === 0) throw new Error('重定价分段不能为空')
  if (segments[0].startMonth !== 1) throw new Error('首段必须从第 1 月开始')

  let expectStart = 1
  for (const segment of segments) {
    if (segment.startMonth !== expectStart) throw new Error('重定价分段不连续')
    if (segment.endMonth < segment.startMonth || segment.endMonth > totalMonths) {
      throw new Error('重定价分段范围非法')
    }
    if (!Number.isFinite(segment.annualRate) || segment.annualRate < 0) {
      throw new Error('重定价分段利率非法')
    }
    expectStart = segment.endMonth + 1
  }

  if (expectStart !== totalMonths + 1) throw new Error('重定价分段未覆盖整个周期')
}

function normalizeMethod(method) {
  const value = String(method ?? '').trim().toLowerCase()
  if (value === '1' || value === METHOD_EPI) return METHOD_EPI
  if (value === '2' || value === METHOD_EP) return METHOD_EP
  return value
}

function toNumber(value) {
  const n = Number.parseFloat(value)
  if (!Number.isFinite(n)) return 0
  return n
}

function toInteger(value) {
  const n = Number.parseInt(value, 10)
  if (!Number.isFinite(n)) return 0
  return n
}

function emptyResult() {
  return {
    ...EMPTY_RESULT,
    schedule: [],
  }
}
