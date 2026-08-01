'use client'

import { useState } from 'react'

import { toast } from 'sonner'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { compareRuns } from '@/lib/api'

const CompareView = () => {
  const [a, setA] = useState('')
  const [b, setB] = useState('')
  const [result, setResult] = useState<Record<string, unknown> | null>(null)
  const [busy, setBusy] = useState(false)

  const run = async () => {
    if (!a.trim() || !b.trim()) {
      toast.error('Both run IDs required')
      return
    }
    setBusy(true)
    try {
      const res = await compareRuns(a.trim(), b.trim())
      setResult(res)
      toast.success('Compare ready')
    } catch (e) {
      toast.error(e instanceof Error ? e.message : 'Compare failed')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className='mx-auto flex max-w-3xl flex-col gap-6'>
      <div>
        <h1 className='font-mono text-xl font-semibold tracking-widest'>COMPARE RUNS</h1>
        <p className='micro-label mt-1'>SIDE-BY-SIDE FINDINGS · VERDICT DELTA · SHARED TITLES</p>
      </div>
      <Card className='hud-corners'>
        <CardHeader>
          <CardTitle className='micro-label'>SELECT RUNS</CardTitle>
        </CardHeader>
        <CardContent className='flex flex-col gap-4'>
          <div>
            <Label className='micro-label'>RUN A</Label>
            <Input value={a} onChange={e => setA(e.target.value)} className='font-mono' placeholder='uuid…' />
          </div>
          <div>
            <Label className='micro-label'>RUN B</Label>
            <Input value={b} onChange={e => setB(e.target.value)} className='font-mono' placeholder='uuid…' />
          </div>
          <Button onClick={run} disabled={busy} className='font-mono text-xs tracking-widest uppercase'>
            {busy ? 'Comparing…' : 'Compare'}
          </Button>
        </CardContent>
      </Card>
      {result && (
        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>
              RESULT — {(result.verdict_change as string) || 'delta'}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <pre className='max-h-[32rem] overflow-auto font-mono text-xs whitespace-pre-wrap'>
              {(result.markdown as string) || JSON.stringify(result, null, 2)}
            </pre>
          </CardContent>
        </Card>
      )}
    </div>
  )
}

export default CompareView
