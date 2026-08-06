'use client'

import { useCallback, useEffect, useRef, useState } from 'react'

import { toast } from 'sonner'

import type { ApprovalAction, ApprovalRiskLevel } from '@/lib/api'
import { approveApproval, listPendingApprovals, listRuns, rejectApproval } from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'

const riskVariant = (
  risk: ApprovalRiskLevel | string
): 'destructive' | 'default' | 'secondary' | 'outline' => {
  switch (risk) {
    case 'critical':
    case 'high':
      return 'destructive'
    case 'medium':
      return 'default'
    case 'low':
      return 'secondary'
    default:
      return 'outline'
  }
}

const stateVariant = (
  state: string
): 'destructive' | 'default' | 'secondary' | 'outline' => {
  switch (state) {
    case 'pending':
      return 'default'
    case 'applied':
      return 'secondary'
    case 'rejected':
    case 'failed':
      return 'destructive'
    default:
      return 'outline'
  }
}

const formatArgs = (args: ApprovalAction['args']): string => {
  if (typeof args === 'string') {
    try {
      return JSON.stringify(JSON.parse(args), null, 2)
    } catch {
      return args
    }
  }
  return JSON.stringify(args, null, 2)
}

const ApprovalsView = () => {
  const [runId, setRunId] = useState('')
  const [committedRunId, setCommittedRunId] = useState('')
  const [recentRuns, setRecentRuns] = useState<{ run_id: string; target: string }[]>([])
  const [items, setItems] = useState<ApprovalAction[] | null>(null)
  const [expanded, setExpanded] = useState<Set<string>>(new Set())
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState<Set<string>>(new Set())
  const initialLoaded = useRef(false)

  const loadPending = useCallback(async (rid: string) => {
    try {
      const res = await listPendingApprovals(rid)
      setItems(res ?? [])
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err))
      setItems([])
    }
  }, [])

  // Load recent runs once for the quick-select.
  useEffect(() => {
    if (initialLoaded.current) return
    initialLoaded.current = true
    listRuns(20)
      .then(r => setRecentRuns((r.runs ?? []).map(x => ({ run_id: x.run_id, target: x.target }))))
      .catch(() => {
        // non-critical
      })
  }, [])

  // Poll pending actions every 5s while a run is selected.
  useEffect(() => {
    if (!committedRunId) return
    void loadPending(committedRunId)
    const timer = setInterval(() => void loadPending(committedRunId), 5000)
    return () => clearInterval(timer)
  }, [committedRunId, loadPending])

  const toggleExpand = (id: string) => {
    setExpanded(prev => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id)
      else next.add(id)
      return next
    })
  }

  const markBusy = (id: string) =>
    setBusy(prev => new Set(prev).add(id))
  const clearBusy = (id: string) =>
    setBusy(prev => {
      const next = new Set(prev)
      next.delete(id)
      return next
    })

  const onApprove = async (id: string) => {
    markBusy(id)
    try {
      await approveApproval(id)
      toast.success('Action approved')
      if (committedRunId) await loadPending(committedRunId)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'approve failed')
    } finally {
      clearBusy(id)
    }
  }

  const onReject = async (id: string) => {
    markBusy(id)
    try {
      await rejectApproval(id, 'Rejected by operator')
      toast.success('Action rejected')
      if (committedRunId) await loadPending(committedRunId)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'reject failed')
    } finally {
      clearBusy(id)
    }
  }

  const submitRunId = (e: React.FormEvent) => {
    e.preventDefault()
    const rid = runId.trim()
    if (!rid) return
    setCommittedRunId(rid)
  }

  return (
    <div className='flex flex-col gap-6'>
      <div>
        <h1 className='font-mono text-xl font-semibold tracking-widest'>APPROVALS</h1>
        <p className='micro-label mt-1'>HUMAN-IN-THE-LOOP QUEUE — DANGEROUS TOOL ACTIONS AWAITING REVIEW · POLLS 5S</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className='micro-label'>SELECT RUN</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={submitRunId} className='flex flex-wrap gap-2'>
            <Input
              placeholder='run_id'
              className='font-mono text-xs'
              value={runId}
              onChange={e => setRunId(e.target.value)}
            />
            <Button type='submit' size='sm' className='font-mono text-[10px] tracking-widest uppercase'>
              Load queue
            </Button>
          </form>
          {recentRuns.length > 0 && (
            <div className='mt-3 flex flex-wrap gap-1'>
              {recentRuns.slice(0, 8).map(r => (
                <button
                  key={r.run_id}
                  type='button'
                  onClick={() => {
                    setRunId(r.run_id)
                    setCommittedRunId(r.run_id)
                  }}
                  className={`rounded border px-2 py-1 font-mono text-[10px] tracking-widest uppercase ${
                    committedRunId === r.run_id
                      ? 'border-primary bg-primary/10 text-primary'
                      : 'border-border text-muted-foreground'
                  }`}
                >
                  {r.target} · {r.run_id.slice(0, 8)}
                </button>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {error && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs'>
          {error}
        </div>
      )}

      {!committedRunId ? (
        <p className='micro-label py-12 text-center'>SELECT A RUN TO VIEW ITS PENDING APPROVAL QUEUE</p>
      ) : !items ? (
        <Skeleton className='h-40 w-full' />
      ) : items.length === 0 ? (
        <p className='micro-label py-12 text-center'>NO PENDING ACTIONS — QUEUE CLEAR</p>
      ) : (
        <div className='flex flex-col gap-3'>
          {items.map(action => {
            const isOpen = expanded.has(action.id)
            return (
              <Card key={action.id}>
                <CardHeader className='pb-2'>
                  <div className='flex flex-wrap items-center gap-2'>
                    <Badge variant={riskVariant(action.risk_level)} className='font-mono text-[10px] uppercase'>
                      {action.risk_level}
                    </Badge>
                    <Badge variant={stateVariant(action.state)} className='font-mono text-[10px] uppercase'>
                      {action.state}
                    </Badge>
                    <span className='text-muted-foreground font-mono text-[10px]'>{action.tool_name}</span>
                    <span className='text-muted-foreground ml-auto font-mono text-[10px]'>
                      {action.run_id.slice(0, 8)} · {new Date(action.created_at).toLocaleTimeString()}
                    </span>
                  </div>
                  <button
                    type='button'
                    onClick={() => toggleExpand(action.id)}
                    className='mt-1 flex items-center gap-2 text-left'
                  >
                    <CardTitle className='font-mono text-sm'>{action.summary || action.tool_name}</CardTitle>
                    <span className='text-muted-foreground font-mono text-[10px]'>{isOpen ? '▲' : '▼'}</span>
                  </button>
                </CardHeader>
                <CardContent className='space-y-3'>
                  {isOpen && (
                    <div className='rounded border border-border/50 bg-black/40 p-3'>
                      <div className='micro-label mb-1'>ARGS</div>
                      <pre className='whitespace-pre-wrap font-mono text-[11px] text-foreground/70'>
                        {formatArgs(action.args)}
                      </pre>
                      <div className='mt-2 micro-label'>ID</div>
                      <div className='text-muted-foreground font-mono text-[11px]'>{action.id}</div>
                    </div>
                  )}
                  {action.state === 'pending' && (
                    <div className='flex gap-2'>
                      <Button
                        size='sm'
                        disabled={busy.has(action.id)}
                        onClick={() => onApprove(action.id)}
                        className='font-mono text-[10px] tracking-widest uppercase'
                      >
                        Approve
                      </Button>
                      <Button
                        size='sm'
                        variant='destructive'
                        disabled={busy.has(action.id)}
                        onClick={() => onReject(action.id)}
                        className='font-mono text-[10px] tracking-widest uppercase'
                      >
                        Reject
                      </Button>
                    </div>
                  )}
                </CardContent>
              </Card>
            )
          })}
        </div>
      )}
    </div>
  )
}

export default ApprovalsView
