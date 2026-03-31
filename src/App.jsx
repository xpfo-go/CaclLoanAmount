import { useMemo, useState } from 'react'
import { calculateReport, prepareScenario } from './core/calculator'
import { LineChart } from './components/LineChart'
import { MetricCard } from './components/MetricCard'
import './App.css'

const DEFAULT_FORM = {
  houseAmount: '160',
  principal: '60',
  fundAmount: '0',
  fundYears: '30',
  fundRatePercent: '2.6',
  commercialYears: '30',
  commercialRatePercent: '3.6',
  method: 'epi',
  commercialChangesText: '',
}

const EXAMPLE_COMBO = {
  houseAmount: '160',
  principal: '60',
  fundAmount: '40',
  fundYears: '30',
  fundRatePercent: '2.6',
  commercialYears: '30',
  commercialRatePercent: '3.6',
  method: 'epi',
  commercialChangesText: '61:3.3,121:3.1',
}

function toNumber(input) {
  const value = Number.parseFloat(input)
  if (!Number.isFinite(value)) return 0
  return value
}

function formatWan(value) {
  return `${value.toFixed(2)} 万元`
}

function App() {
  const [form, setForm] = useState(DEFAULT_FORM)
  const [report, setReport] = useState(null)
  const [errors, setErrors] = useState([])

  const estimateCommercial = useMemo(() => {
    const house = toNumber(form.houseAmount)
    const principal = toNumber(form.principal)
    const fund = toNumber(form.fundAmount)
    return house - principal - fund
  }, [form.fundAmount, form.houseAmount, form.principal])

  const needFundLoan = toNumber(form.fundAmount) > 0
  const needCommercialLoan = estimateCommercial > 0

  function updateField(event) {
    const { name, value } = event.target
    setForm((prev) => ({ ...prev, [name]: value }))
  }

  function applyExample() {
    setForm(EXAMPLE_COMBO)
    setErrors([])
    setReport(null)
  }

  function clearForm() {
    setForm(DEFAULT_FORM)
    setErrors([])
    setReport(null)
  }

  function submit(event) {
    event.preventDefault()
    try {
      const prepared = prepareScenario(form)
      const nextReport = calculateReport(prepared)
      setReport(nextReport)
      setErrors([])
    } catch (error) {
      setReport(null)
      setErrors([error.message || '计算失败，请检查输入'])
    }
  }

  return (
    <main className="page">
      <header className="hero">
        <h1>房贷计算器</h1>
        <p className="subtitle">
          浏览器本地计算，无后端依赖。支持公积金贷、商贷、组合贷与重定价。
        </p>
        <a className="link" href="https://github.com/xpfo-go/CaclLoanAmount" target="_blank" rel="noreferrer">
          GitHub 仓库
        </a>
      </header>

      <section className="workspace">
        <form className="panel form-panel" onSubmit={submit}>
          <h2>输入参数</h2>

          <div className="field-grid">
            <label className="field">
              <span>房屋总价（万元）</span>
              <input name="houseAmount" value={form.houseAmount} onChange={updateField} />
            </label>

            <label className="field">
              <span>首付/本金（万元）</span>
              <input name="principal" value={form.principal} onChange={updateField} />
            </label>

            <label className="field">
              <span>公积金贷款金额（万元）</span>
              <input name="fundAmount" value={form.fundAmount} onChange={updateField} />
            </label>

            <label className="field">
              <span>公积金贷款年限（年）</span>
              <input
                name="fundYears"
                value={form.fundYears}
                onChange={updateField}
                disabled={!needFundLoan}
              />
            </label>

            <label className="field">
              <span>公积金贷款年利率（%）</span>
              <input
                name="fundRatePercent"
                value={form.fundRatePercent}
                onChange={updateField}
                disabled={!needFundLoan}
              />
            </label>

            <label className="field">
              <span>商业贷款年限（年）</span>
              <input
                name="commercialYears"
                value={form.commercialYears}
                onChange={updateField}
                disabled={!needCommercialLoan}
              />
            </label>

            <label className="field">
              <span>商业贷款年利率（%）</span>
              <input
                name="commercialRatePercent"
                value={form.commercialRatePercent}
                onChange={updateField}
                disabled={!needCommercialLoan}
              />
            </label>

            <label className="field">
              <span>还款方式</span>
              <select name="method" value={form.method} onChange={updateField}>
                <option value="epi">等额本息</option>
                <option value="ep">等额本金</option>
              </select>
            </label>

            <label className="field field-full">
              <span>商贷重定价（示例：13:3.2,25:3.1）</span>
              <input
                name="commercialChangesText"
                value={form.commercialChangesText}
                onChange={updateField}
                placeholder="留空表示固定利率"
                disabled={!needCommercialLoan}
              />
            </label>
          </div>

          <p className="hint">
            预估商业贷款金额：{formatWan(Math.max(estimateCommercial, 0))}
          </p>
          <p className="hint">
            当贷款金额为 0 时，对应年限和利率字段可以不填写。
          </p>

          {errors.length > 0 ? (
            <ul className="error-list" role="alert">
              {errors.map((item) => (
                <li key={item}>{item}</li>
              ))}
            </ul>
          ) : null}

          <div className="actions">
            <button type="submit" className="btn btn-primary">
              开始计算
            </button>
            <button type="button" className="btn" onClick={applyExample}>
              填充组合贷示例
            </button>
            <button type="button" className="btn" onClick={clearForm}>
              清空
            </button>
          </div>
        </form>

        <section className="panel result-panel">
          <h2>结果总览</h2>
          {report ? (
            <>
              <div className="metrics">
                <MetricCard label="商业贷款金额" value={formatWan(report.commercialAmount)} />
                <MetricCard label="公积金月供" value={formatWan(report.fund.monthlyPayment)} />
                <MetricCard label="商业贷款月供" value={formatWan(report.commercial.monthlyPayment)} />
                <MetricCard label="组合贷款月供" value={formatWan(report.combo.monthlyPayment)} />
                <MetricCard label="公积金总利息" value={formatWan(report.fund.totalInterest)} />
                <MetricCard label="商业贷款总利息" value={formatWan(report.commercial.totalInterest)} />
                <MetricCard label="组合贷款总利息" value={formatWan(report.combo.totalInterest)} />
              </div>

              <div className="charts">
                {report.fund.schedule.length > 0 ? (
                  <article className="chart-card">
                    <h3>公积金月供趋势</h3>
                    <LineChart data={report.fund.schedule} stroke="#106f5f" />
                  </article>
                ) : null}
                {report.commercial.schedule.length > 0 ? (
                  <article className="chart-card">
                    <h3>商业贷款月供趋势</h3>
                    <LineChart data={report.commercial.schedule} stroke="#ce6a1d" />
                  </article>
                ) : null}
              </div>
            </>
          ) : (
            <p className="empty">填写参数后点击“开始计算”查看结果与折线图。</p>
          )}
        </section>
      </section>
    </main>
  )
}

export default App
