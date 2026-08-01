'use client'

import { useEffect, useState, type ReactNode } from 'react'

import Link from 'next/link'

import { version as nextVersion } from 'next/package.json'
import { ChevronDown, CircleHelp, ExternalLink } from 'lucide-react'
import { toast } from 'sonner'

import type { ConfigEntry, MCPServerInfo, ServiceHealth } from '@/lib/api'
import LiveDot from '@/components/shared/LiveDot'
import PageHeader from '@/components/shared/PageHeader'
import { Button, buttonVariants } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Switch } from '@/components/ui/switch'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { getAgents, getConfig, getMCPServers, getSkills, putConfig, serviceHealth } from '@/lib/api'
import { cn } from '@/lib/utils'

/* --------------------------------- helpers -------------------------------- */

const HelpTip = ({ label, children }: { label: string; children: ReactNode }) => (
  <Tooltip>
    <TooltipTrigger
      type='button'
      className='text-muted-foreground hover:text-primary inline-flex size-5 shrink-0 items-center justify-center rounded-sm border border-border/60 transition-colors'
      aria-label={label}
    >
      <CircleHelp className='size-3' />
    </TooltipTrigger>
    <TooltipContent side='top' className='max-w-xs text-left font-mono text-[11px] leading-relaxed'>
      {children}
    </TooltipContent>
  </Tooltip>
)

const StatusPill = ({ status }: { status: ServiceHealth['status'] | 'loading' }) => {
  if (status === 'loading') return <Skeleton className='size-2 rounded-full' />
  if (status === 'online')
    return (
      <span className='text-primary micro-label flex items-center gap-1.5'>
        <LiveDot tone='green' /> ONLINE
      </span>
    )
  if (status === 'offline')
    return (
      <span className='text-destructive micro-label flex items-center gap-1.5'>
        <LiveDot tone='red' /> OFFLINE
      </span>
    )

  return (
    <span className='micro-label flex items-center gap-1.5'>
      <LiveDot tone='muted' pulse={false} /> N/A
    </span>
  )
}

const ServiceCard = ({ svc }: { svc: ServiceHealth }) => {
  const optional = svc.name === 'ollama' || svc.name === 'onnx-slm'
  return (
    <Card className='hud-corners gap-2 py-3'>
      <CardHeader className='px-4 py-0'>
        <CardTitle className='flex items-center justify-between gap-2 font-mono text-sm tracking-widest'>
          <span className='flex min-w-0 items-center gap-1.5'>
            <span className='truncate'>{svc.name}</span>
            {optional && (
              <HelpTip label={`${svc.name} help`}>
                Optional dependency. Stack can run with LLM_PROVIDER=openai (or other) without this service.
              </HelpTip>
            )}
          </span>
          <StatusPill status={svc.status} />
        </CardTitle>
        <CardDescription className='line-clamp-2 font-mono text-[11px] break-words'>{svc.detail}</CardDescription>
      </CardHeader>
      <CardContent className='flex items-center justify-between gap-2 px-4 pt-1'>
        <p className='text-muted-foreground truncate font-mono text-[11px]' title={svc.endpoint}>
          {svc.endpoint}
        </p>
        {svc.status === 'online' && <p className='micro-label shrink-0'>{svc.latency_ms}ms</p>}
      </CardContent>
    </Card>
  )
}

const SERVICE_SLOTS = 6

/* ---------------------------------- CONFIG --------------------------------- */

const SECRET_MASK = '••••••••'
const LLM_MODEL_KEYS = new Set(['AGENT_MODEL_ID', 'JUDGE_MODEL_ID', 'CODE_MODEL_ID'])

type ConfigGroup = 'llm' | 'attacker' | 'features' | 'other'

const groupOf = (key: string): ConfigGroup => {
  if (key.startsWith('FEATURE_')) return 'features'
  if (key === 'LHOST' || key === 'LPORT') return 'attacker'
  if (
    key === 'LLM_PROVIDER' ||
    key.startsWith('OPENAI_') ||
    key.startsWith('OLLAMA_') ||
    key.startsWith('ONNX_') ||
    LLM_MODEL_KEYS.has(key)
  )
    return 'llm'

  return 'other'
}

const GROUPS: { id: ConfigGroup; title: string; hint: string }[] = [
  {
    id: 'llm',
    title: 'LLM PROVIDER',
    hint: 'bedrock | openai | ollama | onnx. Models and base URLs only apply to the active provider.'
  },
  {
    id: 'attacker',
    title: 'ATTACKER CONTEXT',
    hint: 'LHOST/LPORT used for reverse shells and listener defaults on authorized lab targets.'
  },
  {
    id: 'features',
    title: 'FEATURES',
    hint: 'Runtime feature flags persisted when Postgres is available.'
  },
  { id: 'other', title: 'OTHER', hint: 'Remaining operator-tunable keys from the control plane.' }
]

const LLM_PROVIDERS = ['bedrock', 'openai', 'ollama', 'onnx']
const parseBool = (value: string) => value === 'true' || value === '1'

const SourceBadge = ({ source }: { source: ConfigEntry['source'] }) => {
  const tone =
    source === 'database'
      ? 'text-primary border-primary/40'
      : source === 'env'
        ? 'text-red-400 border-red-400/40'
        : 'text-muted-foreground border-border'

  return <span className={`micro-label rounded-sm border px-1.5 py-0.5 ${tone}`}>{source.toUpperCase()}</span>
}

const ConfigPanel = () => {
  const [entries, setEntries] = useState<ConfigEntry[] | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [draft, setDraft] = useState<Record<string, string | boolean>>({})
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    let mounted = true
    getConfig()
      .then(res => mounted && setEntries(res.config ?? []))
      .catch(err => mounted && setError(err instanceof Error ? err.message : String(err)))
    return () => {
      mounted = false
    }
  }, [])

  const writable = (entries ?? []).length > 0 && (entries ?? []).every(e => e.writable)

  const currentValue = (e: ConfigEntry): string | boolean => {
    if (e.key in draft) return draft[e.key]
    if (e.key.startsWith('FEATURE_')) return parseBool(e.value)
    if (e.secret) return ''
    return e.value
  }

  const setValue = (key: string, value: string | boolean) => setDraft(prev => ({ ...prev, [key]: value }))

  const collectDirty = (): Record<string, unknown> => {
    const dirty: Record<string, unknown> = {}
    for (const e of entries ?? []) {
      if (e.key.startsWith('FEATURE_')) {
        const value = Boolean(currentValue(e))
        if (value !== parseBool(e.value)) dirty[e.key] = value
      } else if (e.secret) {
        const value = String(currentValue(e))
        if (value !== '' && value !== SECRET_MASK) dirty[e.key] = value
      } else {
        const value = String(currentValue(e))
        if (value !== e.value) dirty[e.key] = value
      }
    }
    return dirty
  }

  const dirtyCount = Object.keys(collectDirty()).length

  const save = async () => {
    setSaving(true)
    try {
      const res = await putConfig(collectDirty())
      toast.success(res.note || `UPDATED ${res.updated.length} KEYS`)
      setDraft({})
      const fresh = await getConfig()
      setEntries(fresh.config ?? [])
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to save config')
    } finally {
      setSaving(false)
    }
  }

  const renderControl = (e: ConfigEntry) => {
    if (e.key === 'LLM_PROVIDER') {
      return (
        <Select
          value={String(currentValue(e))}
          onValueChange={value => value && setValue(e.key, value)}
          disabled={!writable}
        >
          <SelectTrigger className='w-full font-mono text-xs tracking-widest uppercase sm:w-64'>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {LLM_PROVIDERS.map(p => (
              <SelectItem key={p} value={p} className='font-mono text-xs tracking-widest uppercase'>
                {p}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      )
    }

    if (e.key.startsWith('FEATURE_')) {
      return (
        <Switch
          checked={Boolean(currentValue(e))}
          onCheckedChange={checked => setValue(e.key, checked)}
          disabled={!writable}
        />
      )
    }

    if (e.secret) {
      return (
        <Input
          type='password'
          value={String(currentValue(e))}
          onChange={ev => setValue(e.key, ev.target.value)}
          placeholder={`${SECRET_MASK} (unchanged)`}
          disabled={!writable}
          className='w-full font-mono text-xs sm:max-w-md'
        />
      )
    }

    return (
      <Input
        value={String(currentValue(e))}
        onChange={ev => setValue(e.key, ev.target.value)}
        disabled={!writable}
        className='w-full font-mono text-xs sm:max-w-md'
      />
    )
  }

  return (
    <Card className='hud-corners'>
      <CardHeader className='flex flex-wrap items-center justify-between gap-3'>
        <div className='flex items-center gap-2'>
          <CardTitle className='micro-label'>CONFIGURATION</CardTitle>
          <HelpTip label='Config help'>
            ENV keys come from process environment. DATABASE source can be edited when Postgres is available. Leave
            secrets blank to keep the current value.
          </HelpTip>
        </div>
        <Button
          size='sm'
          className='font-mono text-xs font-semibold tracking-widest uppercase'
          disabled={!writable || saving || dirtyCount === 0}
          onClick={save}
        >
          {saving ? '[ SAVING… ]' : `[ SAVE CONFIG${dirtyCount > 0 ? ` (${dirtyCount})` : ''} ]`}
        </Button>
      </CardHeader>
      <CardContent className='flex flex-col gap-6'>
        {error ? (
          <p className='text-destructive font-mono text-xs tracking-widest uppercase'>CONFIG UNAVAILABLE — {error}</p>
        ) : entries === null ? (
          <div className='flex flex-col gap-2'>
            <Skeleton className='h-9 w-full' />
            <Skeleton className='h-9 w-full' />
            <Skeleton className='h-9 w-full' />
          </div>
        ) : (
          <>
            {!writable && (
              <div className='border-warning/40 bg-warning/10 text-warning rounded-md border px-4 py-3 font-mono text-xs tracking-widest uppercase'>
                CONFIG EDITING REQUIRES POSTGRES
              </div>
            )}
            {GROUPS.map(group => {
              const groupEntries = entries.filter(e => groupOf(e.key) === group.id)
              if (groupEntries.length === 0) return null

              return (
                <div key={group.id} className='flex flex-col gap-3'>
                  <div className='flex items-center gap-2'>
                    <p className='micro-label'>{group.title}</p>
                    <HelpTip label={group.title}>{group.hint}</HelpTip>
                  </div>
                  <div className='flex flex-col gap-3'>
                    {groupEntries.map(e => (
                      <div key={e.key} className='grid items-center gap-2 sm:grid-cols-[minmax(0,18rem)_1fr]'>
                        <div className='flex min-w-0 items-center gap-2'>
                          <span className='truncate font-mono text-xs' title={e.key}>
                            {e.label || e.key}
                          </span>
                          <SourceBadge source={e.source} />
                          <HelpTip label={e.key}>
                            <span className='font-semibold'>{e.key}</span>
                            {e.hot ? ' · hot-reload' : ''}
                            {e.secret ? ' · secret' : ''}
                            {e.set ? ' · set' : ' · unset'}
                          </HelpTip>
                        </div>
                        <div className='flex items-center gap-2'>{renderControl(e)}</div>
                      </div>
                    ))}
                  </div>
                </div>
              )
            })}
          </>
        )}
      </CardContent>
    </Card>
  )
}

/* --------------------------- operator help (extra) -------------------------- */

const OperatorHelp = ({ coreOnline }: { coreOnline: boolean }) => {
  const [open, setOpen] = useState(false)
  const [skillsCount, setSkillsCount] = useState<number | null>(null)
  const [agentsCount, setAgentsCount] = useState<number | null>(null)
  const [servers, setServers] = useState<MCPServerInfo[] | null>(null)
  const [a2a, setA2a] = useState<{ model?: string; notes?: string[] } | null>(null)
  const [skillStats, setSkillStats] = useState<Record<string, number> | null>(null)
  const [mcpError, setMcpError] = useState<string | null>(null)
  const [expandedServer, setExpandedServer] = useState<string | null>(null)

  useEffect(() => {
    if (!open) return
    getSkills({ brief: true, limit: 1 })
      .then(r => setSkillsCount(r.total ?? r.count ?? 0))
      .catch(() => setSkillsCount(0))
    getAgents()
      .then(r => setAgentsCount(r.count ?? 0))
      .catch(() => setAgentsCount(0))
    getMCPServers()
      .then(res => {
        setServers(res.servers ?? [])
        setA2a(res.agent_to_agent ?? null)
        setSkillStats(res.skill_stats ?? null)
        setMcpError(null)
      })
      .catch(err => setMcpError(err instanceof Error ? err.message : String(err)))
  }, [open])

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <Card className='hud-corners border-primary/20'>
        <CollapsibleTrigger className='flex w-full items-center justify-between gap-3 px-4 py-3 text-left'>
          <div className='flex items-center gap-2'>
            <CircleHelp className='text-primary size-4' />
            <div>
              <p className='font-mono text-xs font-semibold tracking-widest uppercase'>Operator help</p>
              <p className='micro-label mt-0.5'>Architecture · MCP · skills · proxy — toggle for details</p>
            </div>
          </div>
          <ChevronDown className={cn('text-muted-foreground size-4 shrink-0 transition-transform', open && 'rotate-180')} />
        </CollapsibleTrigger>
        <CollapsibleContent>
          <CardContent className='flex flex-col gap-5 border-t border-border/60 pt-4'>
            {/* Quick links */}
            <div>
              <p className='micro-label mb-2'>INTELLIGENCE</p>
              <div className='grid gap-2 sm:grid-cols-3'>
                <Link
                  href='/skills'
                  className={cn(
                    buttonVariants({ variant: 'outline', size: 'sm' }),
                    'h-auto justify-between py-2 font-mono text-[10px] tracking-widest uppercase'
                  )}
                >
                  Skills {skillsCount != null ? `(${skillsCount})` : ''}
                  <ExternalLink className='size-3 opacity-60' />
                </Link>
                <Link
                  href='/agents'
                  className={cn(
                    buttonVariants({ variant: 'outline', size: 'sm' }),
                    'h-auto justify-between py-2 font-mono text-[10px] tracking-widest uppercase'
                  )}
                >
                  Agents {agentsCount != null ? `(${agentsCount})` : ''}
                  <ExternalLink className='size-3 opacity-60' />
                </Link>
                <Link
                  href='/findings'
                  className={cn(
                    buttonVariants({ variant: 'outline', size: 'sm' }),
                    'h-auto justify-between py-2 font-mono text-[10px] tracking-widest uppercase'
                  )}
                >
                  Findings
                  <ExternalLink className='size-3 opacity-60' />
                </Link>
              </div>
            </div>

            {/* A2A notes */}
            {a2a && (
              <div className='rounded-sm border border-primary/20 bg-primary/5 p-3'>
                <p className='micro-label mb-1'>AGENT COMMUNICATION</p>
                <p className='text-muted-foreground mb-2 font-mono text-[11px]'>{a2a.model}</p>
                <ul className='text-muted-foreground space-y-1 font-mono text-[11px]'>
                  {(a2a.notes ?? []).map((n, i) => (
                    <li key={i}>▸ {n}</li>
                  ))}
                </ul>
                {skillStats && (
                  <p className='text-primary mt-2 font-mono text-[11px]'>
                    CyberStrike: {skillStats.total ?? 0} skills (disk {skillStats.src_disk ?? 0} · builtin{' '}
                    {skillStats.src_builtin ?? 0})
                  </p>
                )}
              </div>
            )}

            {/* MCP compact */}
            <div>
              <p className='micro-label mb-2'>MCP SERVERS</p>
              {mcpError ? (
                <p className='text-destructive font-mono text-[11px]'>Unavailable — {mcpError}</p>
              ) : servers === null ? (
                <Skeleton className='h-16 w-full' />
              ) : servers.length === 0 ? (
                <p className='micro-label'>No MCP servers registered</p>
              ) : (
                <div className='flex flex-col gap-2'>
                  {servers.map(s => {
                    const tools = s.tools ?? []
                    const openSrv = expandedServer === s.name
                    return (
                      <div key={s.name} className='rounded-sm border border-border/60 px-3 py-2'>
                        <button
                          type='button'
                          className='flex w-full items-center justify-between gap-2 text-left'
                          onClick={() => setExpandedServer(openSrv ? null : s.name)}
                        >
                          <span className='font-mono text-xs font-semibold tracking-wide'>{s.name}</span>
                          <span className='micro-label text-primary'>{tools.length} tools {openSrv ? '▾' : '▸'}</span>
                        </button>
                        {openSrv && (
                          <ul className='text-muted-foreground mt-2 max-h-40 space-y-0.5 overflow-auto font-mono text-[10px]'>
                            {tools.map(t => (
                              <li key={t}>▸ {t}</li>
                            ))}
                          </ul>
                        )}
                      </div>
                    )
                  })}
                </div>
              )}
            </div>

            {/* Meta footnotes */}
            <div className='grid gap-3 md:grid-cols-3'>
              <div className='rounded-sm border border-border/50 p-3'>
                <p className='micro-label mb-1'>DATA</p>
                <p className='text-muted-foreground font-mono text-[11px] leading-relaxed'>
                  Run state, tool logs and reports live under <span className='text-foreground'>TALON_DATA_DIR</span>.
                  Back up that directory for history.
                </p>
              </div>
              <div className='rounded-sm border border-border/50 p-3'>
                <p className='micro-label mb-1'>VERSION</p>
                <p className='text-muted-foreground font-mono text-[11px] leading-relaxed'>
                  console: Next.js {nextVersion} standalone
                  <br />
                  core: {coreOnline ? 'reachable' : 'offline'}
                </p>
              </div>
              <div className='rounded-sm border border-border/50 p-3'>
                <p className='micro-label mb-1'>API PROXY</p>
                <p className='text-muted-foreground font-mono text-[11px] leading-relaxed'>
                  Browser → <span className='text-foreground'>/api/talon/*</span> → TALON_CORE_URL. Streams: WS → SSE →
                  poll.
                </p>
              </div>
            </div>
          </CardContent>
        </CollapsibleContent>
      </Card>
    </Collapsible>
  )
}

/* --------------------------------- SETTINGS -------------------------------- */

const Settings = () => {
  const [services, setServices] = useState<ServiceHealth[] | null>(null)
  const [probeError, setProbeError] = useState<string | null>(null)
  const [lastProbe, setLastProbe] = useState<Date | null>(null)

  useEffect(() => {
    let mounted = true
    const check = () =>
      serviceHealth()
        .then(res => {
          if (!mounted) return
          setServices(res.services ?? [])
          setProbeError(null)
          setLastProbe(new Date())
        })
        .catch(err => {
          if (!mounted) return
          setProbeError(err instanceof Error ? err.message : 'unreachable')
        })

    check()
    const id = setInterval(check, 10000)
    return () => {
      mounted = false
      clearInterval(id)
    }
  }, [])

  const coreOnline = services?.find(s => s.name === 'talon-core')?.status === 'online'
  const onlineCount = services?.filter(s => s.status === 'online').length ?? 0
  const totalCount = services?.length ?? 0

  return (
    <TooltipProvider delay={200}>
      <div className='flex flex-col gap-6'>
        <PageHeader
          title='SYSTEM'
          subtitle='Service health · configuration'
          action={
            <div className='flex flex-wrap items-center gap-3'>
              {services && (
                <span className='micro-label text-primary'>
                  {onlineCount}/{totalCount} ONLINE
                </span>
              )}
              {lastProbe && (
                <p className='micro-label'>PROBE {lastProbe.toLocaleTimeString('en-GB', { hour12: false })}</p>
              )}
              <HelpTip label='System page help'>
                Live probes hit core every 10s. Optional LLM backends (ollama / onnx-slm) may be offline if you use
                OpenAI. Extra architecture docs are under Operator help.
              </HelpTip>
            </div>
          }
        />

        {probeError && (
          <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs tracking-widest uppercase'>
            CORE UNREACHABLE — {probeError}
          </div>
        )}

        <div className='grid gap-3 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4'>
          {services === null
            ? Array.from({ length: SERVICE_SLOTS }).map((_, i) => <Skeleton key={i} className='h-24 w-full' />)
            : services.map(svc => <ServiceCard key={svc.name} svc={svc} />)}
        </div>

        <ConfigPanel />

        <OperatorHelp coreOnline={Boolean(coreOnline)} />
      </div>
    </TooltipProvider>
  )
}

export default Settings
