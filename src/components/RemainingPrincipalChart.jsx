const HEIGHT = 260
const MARGIN = { top: 20, right: 20, bottom: 34, left: 56 }
const POINT_GAP = 7

function getLabelIndexes(length) {
  if (length <= 1) return [0]
  const indexes = [0, Math.floor((length - 1) / 3), Math.floor(((length - 1) * 2) / 3), length - 1]
  return [...new Set(indexes)]
}

function formatAmount(value) {
  return value.toFixed(2)
}

export function RemainingPrincipalChart({ data }) {
  if (!data || data.length === 0) return null

  const width = MARGIN.left + MARGIN.right + Math.max(1, data.length - 1) * POINT_GAP
  const plotHeight = HEIGHT - MARGIN.top - MARGIN.bottom
  const values = data.map((item) => item.remainingPrincipal)
  const maxValue = Math.max(...values, 0.000001)
  const minValue = Math.min(...values, 0)
  const range = Math.max(maxValue - minValue, 0.000001)
  const labelIndexes = getLabelIndexes(data.length)
  const yTicks = [minValue, minValue + range / 2, maxValue]

  const path = data
    .map((item, index) => {
      const x = MARGIN.left + index * POINT_GAP
      const y = MARGIN.top + ((maxValue - item.remainingPrincipal) / range) * plotHeight
      return `${x},${y}`
    })
    .join(' ')

  return (
    <div className="chart-scroller">
      <svg width={width} height={HEIGHT} viewBox={`0 0 ${width} ${HEIGHT}`} role="img" aria-label="剩余本金趋势图">
        {yTicks.map((tick) => {
          const y = MARGIN.top + ((maxValue - tick) / range) * plotHeight
          return (
            <g key={tick}>
              <line x1={MARGIN.left} y1={y} x2={width - MARGIN.right} y2={y} className="chart-grid" />
              <text x={MARGIN.left - 8} y={y + 4} className="chart-label chart-label-y">
                {formatAmount(tick)}
              </text>
            </g>
          )
        })}

        <polyline points={path} fill="none" className="remaining-line" />

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
          const x = MARGIN.left + index * POINT_GAP
          const month = data[index].month
          return (
            <text key={month} x={x} y={HEIGHT - 10} className="chart-label chart-label-x">
              {month}月
            </text>
          )
        })}
      </svg>
    </div>
  )
}
