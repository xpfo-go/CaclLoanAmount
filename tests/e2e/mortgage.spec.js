import { expect, test } from '@playwright/test'

async function fill(page, label, value) {
  const input = page.getByLabel(label)
  await input.fill('')
  await input.fill(value)
}

test.describe('房贷计算器 e2e', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('仅商贷时显示商业贷款视图和结构化图表', async ({ page }) => {
    await fill(page, '房屋总价（万元）', '160')
    await fill(page, '首付/本金（万元）', '60')
    await fill(page, '公积金贷款金额（万元）', '0')
    await fill(page, '商业贷款年限（年）', '30')
    await fill(page, '商业贷款年利率（%）', '3.6')
    await page.getByRole('button', { name: '开始计算' }).click()

    await expect(page.getByText('商业贷款金额', { exact: true })).toBeVisible()
    await expect(page.getByRole('button', { name: '商业贷款' })).toBeVisible()
    await expect(page.getByRole('heading', { name: '月度还款构成（本金 + 利息）' })).toBeVisible()
    await expect(page.getByRole('heading', { name: '剩余本金趋势' })).toBeVisible()
    await expect(page.getByRole('heading', { name: '摊还明细（按月）' })).toBeVisible()
    await expect(page.getByText('单位：万元')).toHaveCount(2)

    await page.locator('[data-testid="bar-hit-1"]').hover()
    await expect(page.locator('.chart-tooltip')).toContainText('本金：')
    await expect(page.locator('.chart-tooltip')).toContainText('利息：')
    await expect(page.locator('.chart-tooltip')).toContainText('剩余本金：')

    await page.locator('[data-testid="line-point-1"]').dispatchEvent('mouseenter')
    await expect(page.locator('.chart-tooltip')).toContainText('第 1 月')
  })

  test('仅公积金时只显示公积金视图', async ({ page }) => {
    await fill(page, '房屋总价（万元）', '160')
    await fill(page, '首付/本金（万元）', '60')
    await fill(page, '公积金贷款金额（万元）', '100')
    await fill(page, '公积金贷款年限（年）', '30')
    await fill(page, '公积金贷款年利率（%）', '2.6')
    await page.getByRole('button', { name: '开始计算' }).click()

    await expect(page.getByRole('button', { name: '公积金贷款' })).toBeVisible()
    await expect(page.getByRole('button', { name: '商业贷款' })).toHaveCount(0)
  })

  test('组合贷时可切换合并/公积金/商贷视图', async ({ page }) => {
    await fill(page, '房屋总价（万元）', '200')
    await fill(page, '首付/本金（万元）', '50')
    await fill(page, '公积金贷款金额（万元）', '70')
    await fill(page, '公积金贷款年限（年）', '30')
    await fill(page, '公积金贷款年利率（%）', '2.6')
    await fill(page, '商业贷款年限（年）', '30')
    await fill(page, '商业贷款年利率（%）', '3.6')
    await page.getByRole('button', { name: '开始计算' }).click()

    await expect(page.getByRole('button', { name: '合并视图' })).toBeVisible()
    await expect(page.getByRole('button', { name: '公积金贷款' })).toBeVisible()
    await expect(page.getByRole('button', { name: '商业贷款' })).toBeVisible()
    await expect(page.getByText('当前视图：合并视图')).toBeVisible()

    await page.getByRole('button', { name: '商业贷款' }).click()
    await expect(page.getByText('当前视图：商业贷款')).toBeVisible()
  })

  test('非法输入时提示错误信息', async ({ page }) => {
    await fill(page, '房屋总价（万元）', '100')
    await fill(page, '首付/本金（万元）', '100')
    await fill(page, '公积金贷款金额（万元）', '0')
    await page.getByRole('button', { name: '开始计算' }).click()

    await expect(page.getByRole('alert')).toContainText('本金已经覆盖房价，不需要贷款')
  })
})
