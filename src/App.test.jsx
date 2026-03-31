/* @vitest-environment jsdom */
import { afterEach, describe, expect, it } from 'vitest'
import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import '@testing-library/jest-dom/vitest'
import App from './App'

afterEach(() => {
  cleanup()
})

describe('App', () => {
  it('supports only commercial loan when fund amount is zero', async () => {
    const user = userEvent.setup()
    render(<App />)

    await user.clear(screen.getByLabelText('房屋总价（万元）'))
    await user.type(screen.getByLabelText('房屋总价（万元）'), '160')
    await user.clear(screen.getByLabelText('首付/本金（万元）'))
    await user.type(screen.getByLabelText('首付/本金（万元）'), '60')
    await user.clear(screen.getByLabelText('公积金贷款金额（万元）'))
    await user.type(screen.getByLabelText('公积金贷款金额（万元）'), '0')
    await user.clear(screen.getByLabelText('商业贷款年限（年）'))
    await user.type(screen.getByLabelText('商业贷款年限（年）'), '30')
    await user.clear(screen.getByLabelText('商业贷款年利率（%）'))
    await user.type(screen.getByLabelText('商业贷款年利率（%）'), '3.6')

    await user.click(screen.getByRole('button', { name: '开始计算' }))

    expect(screen.getByText('商业贷款金额')).toBeInTheDocument()
    expect(screen.getByText('100.00 万元')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '商业贷款' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '公积金贷款' })).toBeNull()
    expect(screen.getByRole('heading', { name: '月度还款构成（本金 + 利息）' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: '剩余本金趋势' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: '摊还明细（按月）' })).toBeInTheDocument()
  })

  it('shows only fund chart for pure fund loan', async () => {
    const user = userEvent.setup()
    render(<App />)

    await user.clear(screen.getByLabelText('房屋总价（万元）'))
    await user.type(screen.getByLabelText('房屋总价（万元）'), '160')
    await user.clear(screen.getByLabelText('首付/本金（万元）'))
    await user.type(screen.getByLabelText('首付/本金（万元）'), '60')
    await user.clear(screen.getByLabelText('公积金贷款金额（万元）'))
    await user.type(screen.getByLabelText('公积金贷款金额（万元）'), '100')
    await user.clear(screen.getByLabelText('公积金贷款年限（年）'))
    await user.type(screen.getByLabelText('公积金贷款年限（年）'), '30')
    await user.clear(screen.getByLabelText('公积金贷款年利率（%）'))
    await user.type(screen.getByLabelText('公积金贷款年利率（%）'), '2.6')

    await user.click(screen.getByRole('button', { name: '开始计算' }))

    expect(screen.getByRole('button', { name: '公积金贷款' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '商业贷款' })).toBeNull()
    expect(screen.getByRole('heading', { name: '月度还款构成（本金 + 利息）' })).toBeInTheDocument()
  })

  it('shows combined/fund/commercial tabs for combo loan', async () => {
    const user = userEvent.setup()
    render(<App />)

    await user.clear(screen.getByLabelText('房屋总价（万元）'))
    await user.type(screen.getByLabelText('房屋总价（万元）'), '200')
    await user.clear(screen.getByLabelText('首付/本金（万元）'))
    await user.type(screen.getByLabelText('首付/本金（万元）'), '50')
    await user.clear(screen.getByLabelText('公积金贷款金额（万元）'))
    await user.type(screen.getByLabelText('公积金贷款金额（万元）'), '70')
    await user.clear(screen.getByLabelText('公积金贷款年限（年）'))
    await user.type(screen.getByLabelText('公积金贷款年限（年）'), '30')
    await user.clear(screen.getByLabelText('公积金贷款年利率（%）'))
    await user.type(screen.getByLabelText('公积金贷款年利率（%）'), '2.6')
    await user.clear(screen.getByLabelText('商业贷款年限（年）'))
    await user.type(screen.getByLabelText('商业贷款年限（年）'), '30')
    await user.clear(screen.getByLabelText('商业贷款年利率（%）'))
    await user.type(screen.getByLabelText('商业贷款年利率（%）'), '3.6')

    await user.click(screen.getByRole('button', { name: '开始计算' }))

    const comboTab = screen.getByRole('button', { name: '合并视图' })
    expect(comboTab).toBeInTheDocument()
    expect(comboTab).toHaveAttribute('aria-pressed', 'true')
    expect(screen.getByRole('button', { name: '公积金贷款' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '商业贷款' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: '商业贷款' }))
    expect(screen.getByText('当前视图：商业贷款')).toBeInTheDocument()
  })
})
