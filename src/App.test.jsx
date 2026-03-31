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
    expect(screen.getByRole('heading', { name: '商业贷款月供趋势' })).toBeInTheDocument()
    expect(screen.queryByRole('heading', { name: '公积金月供趋势' })).toBeNull()
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

    expect(screen.getByRole('heading', { name: '公积金月供趋势' })).toBeInTheDocument()
    expect(screen.queryByRole('heading', { name: '商业贷款月供趋势' })).toBeNull()
  })

  it('shows two charts for combo loan', async () => {
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

    expect(screen.getByRole('heading', { name: '公积金月供趋势' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: '商业贷款月供趋势' })).toBeInTheDocument()
  })
})
