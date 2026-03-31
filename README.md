# CaclLoanAmount (React Web)

在线体验：https://xpfo-go.github.io/CaclLoanAmount/

一个适合 React 新手学习的房贷计算器前端项目，支持：
- 公积金贷款
- 商业贷款
- 组合贷
- 等额本息 / 等额本金
- 商贷重定价（如 `13:3.2,25:3.1`）
- 月供折线图（按贷款类型自动显示 1 张或 2 张）

## 项目结构

```text
.
├── src/
│   ├── components/         # 通用 UI 组件（指标卡片、折线图）
│   ├── core/               # 纯计算逻辑 + 单元测试
│   ├── App.jsx             # 页面编排与交互状态
│   └── index.css/App.css   # 全局与页面样式
├── docs/
│   ├── requirements.md     # 需求说明
│   └── design.md           # 设计说明
└── legacy/go-cli/          # 旧 Go CLI 归档备份（不再主维护）
```

## 本地开发

```bash
npm install
npm run dev
```

默认访问：http://localhost:5173/

## 测试与检查

```bash
npm run lint
npm test
npm run e2e
npm run build
npm run check
```

## GitHub Actions

仓库包含两条工作流：
- `CI`：`lint + unit test + e2e + build`
- `Deploy Pages`：构建并发布到 GitHub Pages

## 计算口径

- 等额本息：

`M = P * i * (1+i)^N / ((1+i)^N - 1)`

- 等额本金（第 `k` 期）：

`A_k = P/N + (P - (k-1)*P/N) * i`

其中：
- `P` 为贷款本金
- `i` 为月利率（年利率 / 12）
- `N` 为总期数（月）

商贷重定价采用“分段重算”：
- 到达重定价起点后，按“剩余本金 + 剩余期数 + 新利率”重算后续月供。
