import { useState } from 'react'

const HEIGHT = 280
const MARGIN = { top: 20, right: 20, bottom: 34, left: 56 }
const BAR_WIDTH = 9
const BAR_GAP = 3

function getLabelIndexes(length) {
  if (length <= 1) return [0]
  const indexes = [0, Math.floor((length - 1) / 3), Math.floor(((length - 1) * 2) / 3), length - 1]
  return [...new Set(indexes)]
}

function formatAmount(value) {
  return value.toFixed(2)
}

export function PaymentCompositionChart({ data }) {
  const [hover, setHover] = useState(null)

  if (!data || data.length === 0) return null

  const plotHeight = HEIGHT - MARGIN.top - MARGIN.bottom
  const width = MARGIN.left + MARGIN.right + data.length * (BAR_WIDTH + BAR_GAP)
  const maxPayment = Math.max(...data.map((item) => item.payment), 0.000001)
  const yTicks = [0, maxPayment / 3, (maxPayment * 2) / 3, maxPayment]
  const labelIndexes = getLabelIndexes(data.length)

  return (
    <div className="chart-scroller">
      <svg
        width={width}
        height={HEIGHT}
        viewBox={`0 0 ${width} ${HEIGHT}`}
        role="img"
        aria-label="每月本金和利息构成图"
      >
        {yTicks.map((tick) => {
          const y = MARGIN.top + plotHeight - (tick / maxPayment) * plotHeight
          return (
            <g key={tick}>
              <line x1={MARGIN.left} y1={y} x2={width - MARGIN.right} y2={y} className="chart-grid" />
              <text x={MARGIN.left - 8} y={y + 4} className="chart-label chart-label-y">
                {formatAmount(tick)}
              </text>
            </g>
          )
        })}

        <text x={MARGIN.left + 4} y={MARGIN.top - 6} className="chart-unit">
          单位：万元
        </text>

        {data.map((item, index) => {
          const x = MARGIN.left + index * (BAR_WIDTH + BAR_GAP)
          const principalHeight = (item.principal / maxPayment) * plotHeight
          const interestHeight = (item.interest / maxPayment) * plotHeight
          const yPrincipal = MARGIN.top + plotHeight - principalHeight
          const yInterest = yPrincipal - interestHeight

          return (
            <g key={item.month}>
              <rect x={x} y={yPrincipal} width={BAR_WIDTH} height={principalHeight} className="bar-principal" />
              <rect x={x} y={yInterest} width={BAR_WIDTH} height={interestHeight} className="bar-interest" />
              <rect
                x={x - 1}
                y={MARGIN.top}
                width={BAR_WIDTH + 2}
                height={plotHeight}
                className="bar-hit"
                data-testid={`bar-hit-${item.month}`}
                onMouseEnter={() =>
                  setHover({
                    x: x + BAR_WIDTH + 6,
                    y: Math.max(yInterest - 10, MARGIN.top + 8),
                    item,
                  })
                }
                onMouseMove={() =>
                  setHover({
                    x: x + BAR_WIDTH + 6,
                    y: Math.max(yInterest - 10, MARGIN.top + 8),
                    item,
                  })
                }
                onMouseLeave={() => setHover(null)}
              />
            </g>
          )
        })}

        <line
          x1={MARGIN.left}
          y1={MARGIN.top}
          x2={MARGIN.left}
          y2={HEIGHT - MARGIN.bottom}
          className="chart-axis"
        />
        <line
          x1={MARGIN.left}
          y1={HEIGHT - MARGIN.bottom}
          x2={width - MARGIN.right}
          y2={HEIGHT - MARGIN.bottom}
          className="chart-axis"
        />

        {labelIndexes.map((index) => {
          const x = MARGIN.left + index * (BAR_WIDTH + BAR_GAP) + BAR_WIDTH / 2
          const item = data[index]
          return (
            <text key={item.month} x={x} y={HEIGHT - 10} className="chart-label chart-label-x">
              {item.month}月
            </text>
          )
        })}
      </svg>
      <div className="chart-legend">
        <span className="legend-item">
          <i className="legend-dot principal" />
          本金
        </span>
        <span className="legend-item">
          <i className="legend-dot interest" />
          利息
        </span>
      </div>
      {hover ? (
        <div className="chart-tooltip" style={{ left: `${hover.x}px`, top: `${hover.y}px` }}>
          <p>第 {hover.item.month} 月</p>
          <p>本金：{formatAmount(hover.item.principal)} 万元</p>
          <p>利息：{formatAmount(hover.item.interest)} 万元</p>
          <p>剩余本金：{formatAmount(hover.item.remainingPrincipal)} 万元</p>
        </div>
      ) : null}
    </div>
  )
}
