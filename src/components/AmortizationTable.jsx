import { useMemo, useState } from 'react'

function formatAmount(value) {
  return `${value.toFixed(2)} 万元`
}

export function AmortizationTable({ data }) {
  const [showAll, setShowAll] = useState(false)
  const rows = useMemo(() => {
    if (!data) return []
    if (showAll) return data
    return data.slice(0, 24)
  }, [data, showAll])

  if (!data || data.length === 0) return null

  return (
    <div className="table-wrap">
      <div className="table-actions">
        <button type="button" className="btn btn-table" onClick={() => setShowAll((prev) => !prev)}>
          {showAll ? '收起' : `显示全部 ${data.length} 期`}
        </button>
      </div>
      <div className="table-scroller">
        <table className="schedule-table">
          <thead>
            <tr>
              <th>期数</th>
              <th>月供</th>
              <th>本金</th>
              <th>利息</th>
              <th>剩余本金</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((item) => (
              <tr key={item.month}>
                <td>{item.month}</td>
                <td>{formatAmount(item.payment)}</td>
                <td>{formatAmount(item.principal)}</td>
                <td>{formatAmount(item.interest)}</td>
                <td>{formatAmount(item.remainingPrincipal)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
