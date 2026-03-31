const MARGIN = { top: 18, right: 22, bottom: 30, left: 56 }
const WIDTH = 760
const HEIGHT = 260

function createPathPoints(data, minValue, maxValue) {
  const plotWidth = WIDTH - MARGIN.left - MARGIN.right
  const plotHeight = HEIGHT - MARGIN.top - MARGIN.bottom
  const range = Math.max(maxValue - minValue, 0.000001)

  return data.map((item, index) => {
    const ratioX = data.length > 1 ? index / (data.length - 1) : 0
    const x = MARGIN.left + ratioX * plotWidth
    const y = MARGIN.top + ((maxValue - item.payment) / range) * plotHeight
    return [x, y]
  })
}

function valueLabel(value) {
  return value.toFixed(2)
}

export function LineChart({ data, stroke }) {
  if (!data || data.length === 0) return null

  const payments = data.map((item) => item.payment)
  const rawMin = Math.min(...payments)
  const rawMax = Math.max(...payments)
  const range = rawMax - rawMin
  const minValue = rawMin - (range === 0 ? rawMin * 0.1 : range * 0.12)
  const maxValue = rawMax + (range === 0 ? rawMax * 0.1 : range * 0.12)

  const points = createPathPoints(data, minValue, maxValue)
  const polyline = points.map(([x, y]) => `${x},${y}`).join(' ')

  const yTicks = [minValue, (minValue + maxValue) / 2, maxValue]
  const xLastMonth = data.at(-1)?.month ?? 1
  const xTicks = [1, Math.ceil(xLastMonth / 2), xLastMonth]

  return (
    <svg className="line-chart" viewBox={`0 0 ${WIDTH} ${HEIGHT}`} role="img" aria-label="月供折线图">
      {yTicks.map((tick) => {
        const y =
          MARGIN.top +
          ((maxValue - tick) / Math.max(maxValue - minValue, 0.000001)) *
            (HEIGHT - MARGIN.top - MARGIN.bottom)

        return (
          <g key={tick}>
            <line x1={MARGIN.left} y1={y} x2={WIDTH - MARGIN.right} y2={y} className="chart-grid" />
            <text x={MARGIN.left - 8} y={y + 4} className="chart-label chart-label-y">
              {valueLabel(tick)}
            </text>
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
        x2={WIDTH - MARGIN.right}
        y2={HEIGHT - MARGIN.bottom}
        className="chart-axis"
      />

      {xTicks.map((tick) => {
        const ratio = xLastMonth > 1 ? (tick - 1) / (xLastMonth - 1) : 0
        const x = MARGIN.left + ratio * (WIDTH - MARGIN.left - MARGIN.right)
        return (
          <text key={tick} x={x} y={HEIGHT - 8} className="chart-label chart-label-x">
            {tick}月
          </text>
        )
      })}

      <polyline points={polyline} fill="none" stroke={stroke} strokeWidth="3" strokeLinecap="round" />
    </svg>
  )
}
