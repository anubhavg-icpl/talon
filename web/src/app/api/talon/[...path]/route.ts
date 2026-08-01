// Next Imports
import type { NextRequest } from 'next/server'

const TALON_CORE_URL = process.env.TALON_CORE_URL ?? 'http://localhost:8000'

const HOP_BY_HOP = new Set(['connection', 'keep-alive', 'transfer-encoding', 'upgrade', 'host', 'content-length'])

async function proxy(req: NextRequest, ctx: { params: Promise<{ path: string[] }> }) {
  const { path } = await ctx.params
  const upstreamUrl = `${TALON_CORE_URL}/${path.join('/')}${req.nextUrl.search}`

  const headers = new Headers()

  req.headers.forEach((value, key) => {
    if (!HOP_BY_HOP.has(key.toLowerCase())) headers.set(key, value)
  })

  const hasBody = req.method !== 'GET' && req.method !== 'HEAD'

  const upstream = await fetch(upstreamUrl, {
    method: req.method,
    headers,
    body: hasBody ? req.body : undefined,

    // Required by undici when forwarding a streaming request body
    // @ts-expect-error -- duplex is not in the TS lib types yet
    duplex: hasBody ? 'half' : undefined,
    cache: 'no-store'
  })

  const resHeaders = new Headers()

  upstream.headers.forEach((value, key) => {
    if (!HOP_BY_HOP.has(key.toLowerCase())) resHeaders.set(key, value)
  })

  // Stream the body straight through — keeps text/event-stream unbuffered.
  // Force no proxy buffering so /llm/assist heartbeats reach the browser.
  if ((upstream.headers.get('content-type') || '').includes('text/event-stream')) {
    resHeaders.set('Cache-Control', 'no-cache, no-transform')
    resHeaders.set('X-Accel-Buffering', 'no')
    resHeaders.set('Connection', 'keep-alive')
  }

  return new Response(upstream.body, {
    status: upstream.status,
    statusText: upstream.statusText,
    headers: resHeaders
  })
}

export const GET = proxy
export const POST = proxy
export const PUT = proxy
export const PATCH = proxy
export const DELETE = proxy

export const dynamic = 'force-dynamic'
