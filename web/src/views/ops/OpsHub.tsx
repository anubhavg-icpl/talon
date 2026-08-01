'use client'

import { useEffect, useState } from 'react'

import Link from 'next/link'

import { toast } from 'sonner'

import type { BudgetStats, Credential, NotifyConfig, Schedule, ScopePolicy, Target } from '@/lib/api'
import GlobePanel from '@/components/shared/GlobePanel'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Textarea } from '@/components/ui/textarea'
import {
  addCredential,
  batchStart,
  deleteCredential,
  deleteSchedule,
  deleteTarget,
  getBudget,
  getNotify,
  getScope,
  listCredentials,
  listSchedules,
  listTargets,
  putNotify,
  putScope,
  upsertSchedule,
  upsertTarget
} from '@/lib/api'

const OpsHub = () => {
  const [scope, setScope] = useState<ScopePolicy | null>(null)
  const [targets, setTargets] = useState<Target[]>([])
  const [schedules, setSchedules] = useState<Schedule[]>([])
  const [notify, setNotify] = useState<NotifyConfig | null>(null)
  const [creds, setCreds] = useState<Credential[]>([])
  const [budget, setBudget] = useState<BudgetStats | null>(null)
  const [batchIPs, setBatchIPs] = useState('10.0.0.1\n10.0.0.2')
  const [tgtAddr, setTgtAddr] = useState('')
  const [schName, setSchName] = useState('')
  const [schTarget, setSchTarget] = useState('')
  const [schInterval, setSchInterval] = useState('24h')
  const [credName, setCredName] = useState('')
  const [credSecret, setCredSecret] = useState('')
  const [credUser, setCredUser] = useState('')

  const reload = async () => {
    const [sc, tg, sch, nt, cr, bd] = await Promise.all([
      getScope().catch(() => null),
      listTargets().catch(() => ({ targets: [] })),
      listSchedules().catch(() => ({ schedules: [] })),
      getNotify().catch(() => null),
      listCredentials().catch(() => ({ credentials: [] })),
      getBudget().catch(() => null)
    ])
    if (sc) setScope(sc)
    setTargets(tg.targets ?? [])
    setSchedules(sch.schedules ?? [])
    if (nt) setNotify(nt)
    setCreds(cr.credentials ?? [])
    if (bd) setBudget(bd)
  }

  useEffect(() => {
    void reload()
  }, [])

  return (
    <div className='flex flex-col gap-8'>
      <div className='flex flex-wrap items-end justify-between gap-4'>
        <div>
          <h1 className='font-mono text-xl font-semibold tracking-widest'>ENGAGEMENTS</h1>
          <p className='micro-label mt-1'>
            BATCH · TARGETS · SCOPE/ROE · SCHEDULES · WEBHOOKS · CREDENTIALS · BUDGET
          </p>
        </div>
        <GlobePanel className='hidden size-28 shrink-0 sm:block' state='running' activityLevel={0.4} gateCompact />
      </div>

      {budget && (
        <div className='grid grid-cols-2 gap-3 md:grid-cols-5'>
          {[
            ['STARTED', budget.runs_started],
            ['COMPLETED', budget.runs_completed],
            ['TOOLS', budget.tool_calls],
            ['LLM', budget.llm_calls],
            ['CRITICAL', budget.critical_findings]
          ].map(([l, v]) => (
            <Card key={String(l)} className='py-3'>
              <CardHeader className='px-3 py-0'>
                <CardTitle className='micro-label'>{l}</CardTitle>
              </CardHeader>
              <CardContent className='px-3 pt-1 font-mono text-2xl'>{v}</CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Batch */}
      <Card className='hud-corners'>
        <CardHeader>
          <CardTitle className='micro-label'>BATCH START</CardTitle>
          <CardDescription className='font-mono text-xs'>One IP per line · max 50 · description marked AUTHORIZED</CardDescription>
        </CardHeader>
        <CardContent className='space-y-3'>
          <Textarea value={batchIPs} onChange={e => setBatchIPs(e.target.value)} className='min-h-28 font-mono text-xs' />
          <Button
            className='font-mono text-xs tracking-widest uppercase'
            onClick={async () => {
              const ips = batchIPs.split(/[\n,]+/).map(s => s.trim()).filter(Boolean)
              try {
                const res = await batchStart({
                  ips,
                  description: 'Batch engagement AUTHORIZED',
                  agent_mode: 'full'
                })
                toast.success(`Started ${res.count} runs`)
              } catch (e) {
                toast.error(e instanceof Error ? e.message : 'batch failed')
              }
            }}
          >
            Launch batch
          </Button>
        </CardContent>
      </Card>

      {/* Scope */}
      <Card>
        <CardHeader>
          <CardTitle className='micro-label'>SCOPE / ROE</CardTitle>
        </CardHeader>
        <CardContent className='space-y-3 font-mono text-xs'>
          {scope && (
            <>
              <div className='flex items-center justify-between'>
                <Label>Enforce policy</Label>
                <Switch
                  checked={scope.enabled}
                  onCheckedChange={v => setScope({ ...scope, enabled: v })}
                />
              </div>
              <div className='flex items-center justify-between'>
                <Label>Require AUTHORIZED in description</Label>
                <Switch
                  checked={scope.require_auth_label}
                  onCheckedChange={v => setScope({ ...scope, require_auth_label: v })}
                />
              </div>
              <div className='flex items-center justify-between'>
                <Label>Auto-approve nmap on private IPs</Label>
                <Switch
                  checked={scope.auto_approve_nmap_private}
                  onCheckedChange={v => setScope({ ...scope, auto_approve_nmap_private: v })}
                />
              </div>
              <div>
                <Label>Allowed CIDRs (comma-separated)</Label>
                <Input
                  className='font-mono'
                  value={(scope.allowed_cidrs || []).join(', ')}
                  onChange={e =>
                    setScope({
                      ...scope,
                      allowed_cidrs: e.target.value.split(',').map(s => s.trim()).filter(Boolean)
                    })
                  }
                />
              </div>
              <div>
                <Label>Denied CIDRs</Label>
                <Input
                  className='font-mono'
                  value={(scope.denied_cidrs || []).join(', ')}
                  onChange={e =>
                    setScope({
                      ...scope,
                      denied_cidrs: e.target.value.split(',').map(s => s.trim()).filter(Boolean)
                    })
                  }
                />
              </div>
              <div>
                <Label>Max concurrent</Label>
                <Input
                  type='number'
                  className='font-mono'
                  value={scope.max_concurrent}
                  onChange={e => setScope({ ...scope, max_concurrent: Number(e.target.value) || 0 })}
                />
              </div>
              <Button
                size='sm'
                className='font-mono text-[10px] tracking-widest uppercase'
                onClick={async () => {
                  try {
                    const saved = await putScope(scope)
                    setScope(saved)
                    toast.success('Scope saved')
                  } catch (e) {
                    toast.error(e instanceof Error ? e.message : 'save failed')
                  }
                }}
              >
                Save scope
              </Button>
            </>
          )}
        </CardContent>
      </Card>

      <div className='grid gap-6 lg:grid-cols-2'>
        {/* Targets */}
        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>TARGETS</CardTitle>
          </CardHeader>
          <CardContent className='space-y-3'>
            <div className='flex gap-2'>
              <Input
                placeholder='IP / host / URL'
                className='font-mono text-xs'
                value={tgtAddr}
                onChange={e => setTgtAddr(e.target.value)}
              />
              <Button
                size='sm'
                onClick={async () => {
                  if (!tgtAddr.trim()) return
                  const isURL = tgtAddr.includes('://')
                  await upsertTarget(isURL ? { url: tgtAddr.trim(), address: '' } : { address: tgtAddr.trim() })
                  setTgtAddr('')
                  await reload()
                  toast.success('Target saved')
                }}
              >
                Add
              </Button>
            </div>
            <div className='max-h-64 space-y-1 overflow-auto'>
              {targets.map(t => (
                <div key={t.id} className='flex items-center gap-2 border-b border-border/40 py-1 font-mono text-[11px]'>
                  <span className='flex-1 truncate'>{t.label || t.address || t.url}</span>
                  {t.last_status && <Badge variant='outline'>{t.last_status}</Badge>}
                  <Link href={`/runs/new?ip=${encodeURIComponent(t.address || t.url || '')}`} className='text-primary'>
                    run
                  </Link>
                  <button type='button' className='text-destructive' onClick={() => deleteTarget(t.id).then(reload)}>
                    ×
                  </button>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Schedules */}
        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>SCHEDULES</CardTitle>
            <CardDescription className='font-mono text-[10px]'>Intervals: 1h, 6h, 24h, 7d</CardDescription>
          </CardHeader>
          <CardContent className='space-y-3'>
            <Input placeholder='Name' className='font-mono text-xs' value={schName} onChange={e => setSchName(e.target.value)} />
            <Input placeholder='Target IP' className='font-mono text-xs' value={schTarget} onChange={e => setSchTarget(e.target.value)} />
            <Input placeholder='Interval' className='font-mono text-xs' value={schInterval} onChange={e => setSchInterval(e.target.value)} />
            <Button
              size='sm'
              onClick={async () => {
                await upsertSchedule({
                  name: schName,
                  target: schTarget,
                  interval: schInterval,
                  enabled: true,
                  playbook_id: 'full-validation'
                })
                setSchName('')
                setSchTarget('')
                await reload()
                toast.success('Schedule saved')
              }}
            >
              Add schedule
            </Button>
            <div className='max-h-48 space-y-1 overflow-auto'>
              {schedules.map(s => (
                <div key={s.id} className='flex gap-2 border-b border-border/40 py-1 font-mono text-[11px]'>
                  <span className='flex-1'>
                    {s.name} → {s.target} every {s.interval}
                    {s.enabled ? '' : ' (off)'}
                  </span>
                  <button type='button' className='text-destructive' onClick={() => deleteSchedule(s.id).then(reload)}>
                    ×
                  </button>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      <div className='grid gap-6 lg:grid-cols-2'>
        {/* Notify */}
        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>WEBHOOK NOTIFICATIONS</CardTitle>
          </CardHeader>
          <CardContent className='space-y-3 font-mono text-xs'>
            {notify && (
              <>
                <Input
                  placeholder='https://hooks.slack.com/…'
                  className='font-mono'
                  value={notify.webhook_url}
                  onChange={e => setNotify({ ...notify, webhook_url: e.target.value })}
                />
                {(
                  [
                    ['on_complete', 'On complete'],
                    ['on_hitl', 'On HITL'],
                    ['on_critical_finding', 'On critical'],
                    ['on_error', 'On error']
                  ] as const
                ).map(([k, label]) => (
                  <div key={k} className='flex items-center justify-between'>
                    <Label>{label}</Label>
                    <Switch
                      checked={Boolean(notify[k])}
                      onCheckedChange={v => setNotify({ ...notify, [k]: v })}
                    />
                  </div>
                ))}
                <Button
                  size='sm'
                  onClick={async () => {
                    setNotify(await putNotify(notify))
                    toast.success('Notify saved')
                  }}
                >
                  Save webhooks
                </Button>
              </>
            )}
          </CardContent>
        </Card>

        {/* Credentials */}
        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>CREDENTIALS (ENCRYPTED)</CardTitle>
            <CardDescription className='font-mono text-[10px]'>AES-GCM · TALON_CRED_KEY</CardDescription>
          </CardHeader>
          <CardContent className='space-y-2'>
            <Input placeholder='Name' className='font-mono text-xs' value={credName} onChange={e => setCredName(e.target.value)} />
            <Input placeholder='Username' className='font-mono text-xs' value={credUser} onChange={e => setCredUser(e.target.value)} />
            <Input
              type='password'
              placeholder='Secret'
              className='font-mono text-xs'
              value={credSecret}
              onChange={e => setCredSecret(e.target.value)}
            />
            <Button
              size='sm'
              onClick={async () => {
                await addCredential({ name: credName, username: credUser, secret: credSecret })
                setCredName('')
                setCredUser('')
                setCredSecret('')
                await reload()
                toast.success('Credential stored')
              }}
            >
              Add credential
            </Button>
            <div className='max-h-40 space-y-1 overflow-auto'>
              {creds.map(c => (
                <div key={c.id} className='flex gap-2 font-mono text-[11px]'>
                  <span className='flex-1'>
                    {c.name} {c.username && `(${c.username})`}
                  </span>
                  <button type='button' className='text-destructive' onClick={() => deleteCredential(c.id).then(reload)}>
                    ×
                  </button>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default OpsHub
