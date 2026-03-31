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

  test('仅商贷时只显示商业贷款折线图', async ({ page }) => {
    await fill(page, '房屋总价（万元）', '160')
    await fill(page, '首付/本金（万元）', '60')
    await fill(page, '公积金贷款金额（万元）', '0')
    await fill(page, '商业贷款年限（年）', '30')
    await fill(page, '商业贷款年利率（%）', '3.6')
    await page.getByRole('button', { name: '开始计算' }).click()

    await expect(page.getByText('商业贷款金额', { exact: true })).toBeVisible()
    await expect(page.getByRole('heading', { name: '商业贷款月供趋势' })).toBeVisible()
    await expect(page.getByRole('heading', { name: '公积金月供趋势' })).toHaveCount(0)
  })

  test('仅公积金时只显示公积金折线图', async ({ page }) => {
    await fill(page, '房屋总价（万元）', '160')
    await fill(page, '首付/本金（万元）', '60')
    await fill(page, '公积金贷款金额（万元）', '100')
    await fill(page, '公积金贷款年限（年）', '30')
    await fill(page, '公积金贷款年利率（%）', '2.6')
    await page.getByRole('button', { name: '开始计算' }).click()

    await expect(page.getByRole('heading', { name: '公积金月供趋势' })).toBeVisible()
    await expect(page.getByRole('heading', { name: '商业贷款月供趋势' })).toHaveCount(0)
  })

  test('组合贷时显示两张折线图', async ({ page }) => {
    await fill(page, '房屋总价（万元）', '200')
    await fill(page, '首付/本金（万元）', '50')
    await fill(page, '公积金贷款金额（万元）', '70')
    await fill(page, '公积金贷款年限（年）', '30')
    await fill(page, '公积金贷款年利率（%）', '2.6')
    await fill(page, '商业贷款年限（年）', '30')
    await fill(page, '商业贷款年利率（%）', '3.6')
    await page.getByRole('button', { name: '开始计算' }).click()

    await expect(page.getByRole('heading', { name: '公积金月供趋势' })).toBeVisible()
    await expect(page.getByRole('heading', { name: '商业贷款月供趋势' })).toBeVisible()
  })

  test('非法输入时提示错误信息', async ({ page }) => {
    await fill(page, '房屋总价（万元）', '100')
    await fill(page, '首付/本金（万元）', '100')
    await fill(page, '公积金贷款金额（万元）', '0')
    await page.getByRole('button', { name: '开始计算' }).click()

    await expect(page.getByRole('alert')).toContainText('本金已经覆盖房价，不需要贷款')
  })
})
