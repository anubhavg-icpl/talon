'use client'

/** Root crash boundary (must include html/body). */
export default function GlobalError({
  error,
  reset
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  return (
    <html lang='en'>
      <body style={{ background: '#0a0a0b', color: '#fafafa', fontFamily: 'ui-monospace, monospace', padding: 48 }}>
        <h1 style={{ color: '#ef4444', letterSpacing: '0.2em' }}>TALON // FAULT</h1>
        <p style={{ opacity: 0.7, maxWidth: 480 }}>
          The operator shell crashed. Retry, or open a lighter page.
        </p>
        <pre style={{ opacity: 0.5, fontSize: 12, marginTop: 16 }}>{error.message}</pre>
        <div style={{ display: 'flex', gap: 12, marginTop: 24 }}>
          <button
            type='button'
            onClick={reset}
            style={{ padding: '8px 16px', background: '#ef4444', border: 0, color: '#fff', cursor: 'pointer' }}
          >
            RETRY
          </button>
          <a href='/login' style={{ padding: '8px 16px', border: '1px solid #333', color: '#fff', textDecoration: 'none' }}>
            LOGIN
          </a>
          <a href='/assist' style={{ padding: '8px 16px', border: '1px solid #333', color: '#fff', textDecoration: 'none' }}>
            ASSIST
          </a>
        </div>
      </body>
    </html>
  )
}
