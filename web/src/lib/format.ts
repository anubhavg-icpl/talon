/** Relative time like "3m ago" / "2h ago" / "4d ago". */
export const relativeTime = (iso: string, now = Date.now()): string => {
  const then = new Date(iso).getTime()

  if (Number.isNaN(then)) return '—'

  const diff = Math.max(0, now - then)
  const sec = Math.floor(diff / 1000)

  if (sec < 60) return `${sec}s ago`

  const min = Math.floor(sec / 60)

  if (min < 60) return `${min}m ago`

  const hrs = Math.floor(min / 60)

  if (hrs < 24) return `${hrs}h ago`

  const days = Math.floor(hrs / 24)

  return `${days}d ago`
}

/** Elapsed duration like "01:23:45" from a start timestamp. */
export const elapsed = (iso: string, now = Date.now()): string => {
  const then = new Date(iso).getTime()

  if (Number.isNaN(then)) return '--:--:--'

  const total = Math.max(0, Math.floor((now - then) / 1000))
  const h = Math.floor(total / 3600)
  const m = Math.floor((total % 3600) / 60)
  const s = total % 60

  const pad = (n: number) => String(n).padStart(2, '0')

  return `${pad(h)}:${pad(m)}:${pad(s)}`
}

export const shortId = (id: string, head = 8): string => (id.length > head ? `${id.slice(0, head)}…` : id)
