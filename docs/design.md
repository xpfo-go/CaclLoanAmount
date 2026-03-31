# 房贷计算器 Web 版设计（React）

## 1. 技术选型
- 构建工具：Vite
- 框架：React（JavaScript）
- 测试：Vitest
- 部署：GitHub Pages + GitHub Actions

## 2. 目录设计
- `src/core/`
  - 纯函数计算引擎（等额本息、等额本金、分段利率、组合贷汇总）
  - 输入解析和格式化（重定价字符串解析、金额格式化）
- `src/features/calculator/`
  - 计算器页面业务逻辑与状态管理
- `src/components/ui/`
  - 输入卡片、结果卡片、按钮、错误提示
- `src/components/charts/`
  - SVG 折线图组件
- `src/styles/`
  - 全局主题、布局、组件样式
- `legacy/go-cli/`
  - 原 CLI 备份

## 3. 页面布局
- 顶部：项目标题 + 在线部署说明
- 左栏：参数输入区
- 右栏：结果摘要卡片 + 图表区
- 底部：说明与旧版 CLI 备份提示

## 4. 关键交互
- 输入变更后点击“开始计算”触发计算。
- 校验失败时在输入区显示错误列表。
- 成功后渲染：
  - 汇总指标卡片
  - 公积金/商贷折线图（按实际贷款类型显示 1 或 2 张）

## 5. 计算设计
- 沿用并前端化现有 Go 版本计算口径：
  - 固定利率：等额本息 / 等额本金
  - 分段利率：在分段起点按“剩余本金 + 剩余期数 + 新利率”重算
- 单位统一：万元

## 6. 部署设计
- Vite `build` 输出 `dist/`
- GitHub Actions：
  - `ci.yml`：lint + test + build
  - `deploy-pages.yml`：构建并发布 `dist/` 到 Pages
- `README` 顶部放在线体验链接。
