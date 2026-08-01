/**
 * Optional in-browser SLM path (WASM / Transformers.js).
 *
 * Production chat uses Go → POST /llm/assist (tools) or /llm/stream.
 * This module is only for offline demos (SmolLM2-135M in the browser).
 *
 * @huggingface/transformers is intentionally NOT a hard dependency — load it
 * at runtime if installed. The dashboard build must not require it.
 */

export type WASMStreamHandlers = {
  onToken?: (token: string) => void
  onDone?: (text: string) => void
  onError?: (error: Error) => void
  onStatus?: (status: string) => void
}

const DEFAULT_WASM_MODEL = 'HuggingFaceTB/SmolLM2-135M-Instruct'

type TransformersMod = {
  pipeline: (
    task: string,
    model: string,
    opts?: Record<string, unknown>
  ) => Promise<(messages: unknown, genOpts?: Record<string, unknown>) => Promise<unknown>>
  env: { allowLocalModels: boolean }
}

async function loadTransformers(): Promise<TransformersMod> {
  // Dynamic string keeps TypeScript from resolving a missing package at build time.
  const spec = '@huggingface/' + 'transformers'
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const mod = await (Function('s', 'return import(s)') as (s: string) => Promise<any>)(spec)

  return mod as TransformersMod
}

/** Stream a reply from in-browser SmolLM (Transformers.js + ORT-Web WASM). */
export async function streamSLMWASM(
  messages: { role: string; content: string }[],
  handlers: WASMStreamHandlers = {},
  modelId: string = DEFAULT_WASM_MODEL
): Promise<() => void> {
  let cancelled = false

  try {
    handlers.onStatus?.('loading transformers.js + ORT-Web WASM…')
    const { pipeline, env } = await loadTransformers()

    env.allowLocalModels = false
    handlers.onStatus?.(`loading ${modelId}…`)
    const generator = await pipeline('text-generation', modelId, {
      dtype: 'q4',
      device: 'wasm'
    })

    if (cancelled) return () => undefined

    handlers.onStatus?.('generating…')
    const out = (await generator(messages, { max_new_tokens: 128, do_sample: false })) as
      | { generated_text?: unknown }[]
      | null
    const generated = out?.[0]?.generated_text
    let text = ''

    if (Array.isArray(generated)) {
      const last = generated[generated.length - 1] as { content?: string }

      text = typeof last?.content === 'string' ? last.content : String(last ?? '')
    } else if (typeof generated === 'string') {
      text = generated
    } else {
      text = String(generated ?? '')
    }

    if (!cancelled) {
      handlers.onToken?.(text)
      handlers.onDone?.(text)
    }
  } catch (err) {
    if (!cancelled) {
      handlers.onError?.(err instanceof Error ? err : new Error(String(err)))
    }
  }

  return () => {
    cancelled = true
  }
}

export function wasmSLMSupported(): boolean {
  if (typeof window === 'undefined') return false

  return typeof WebAssembly !== 'undefined'
}
