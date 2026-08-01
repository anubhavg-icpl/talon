'use client'

// React Imports
import { useEffect, useState } from 'react'

// Third-party Imports
import { toast } from 'sonner'

// Type Imports
import type { ConfigEntry, MCPServerInfo, ServiceHealth } from '@/lib/api'

// Component Imports
import LiveDot from '@/components/shared/LiveDot'
import { Button, buttonVariants } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Switch } from '@/components/ui/switch'

// Next Imports
import Link from 'next/link'

// Util Imports
import { getAgents, getConfig, getMCPServers, getSkills, putConfig, serviceHealth } from '@/lib/api'
import { cn } from '@/lib/utils'

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

const ServiceCard = ({ svc }: { svc: ServiceHealth }) => (
  <Card className='hud-corners gap-2 py-4'>
    <CardHeader className='px-4'>
      <CardTitle className='flex items-center justify-between font-mono text-sm tracking-widest'>
        {svc.name}
        <StatusPill status={svc.status} />
      </CardTitle>
      <CardDescription className='font-mono text-xs break-words'>{svc.detail}</CardDescription>
    </CardHeader>
    <CardContent className='flex items-center justify-between px-4'>
      <p className='text-muted-foreground font-mono text-xs break-all'>{svc.endpoint}</p>
      {svc.status === 'online' && <p className='micro-label shrink-0 pl-2'>{svc.latency_ms}ms</p>}
    </CardContent>
  </Card>
)

const SERVICE_SLOTS = 6

/* ---------------------------------- CONFIG --------------------------------- */

const SECRET_MASK = '••••••••'

const LLM_MODEL_KEYS = new Set(['AGENT_MODEL_ID', 'JUDGE_MODEL_ID', 'CODE_MODEL_ID'])

type ConfigGroup = 'llm' | 'attacker' | 'features' | 'other'

const groupOf = (key: string): ConfigGroup => {
  if (key.startsWith('FEATURE_')) return 'features'
  if (key === 'LHOST' || key === 'LPORT') return 'attacker'
  if (key === 'LLM_PROVIDER' || key.startsWith('OPENAI_') || key.startsWith('OLLAMA_') || LLM_MODEL_KEYS.has(key))
    return 'llm'

  return 'other'
}

const GROUPS: { id: ConfigGroup; title: string }[] = [
  { id: 'llm', title: 'LLM PROVIDER' },
  { id: 'attacker', title: 'ATTACKER CONTEXT' },
  { id: 'features', title: 'FEATURES' },
  { id: 'other', title: 'OTHER' }
]

const LLM_PROVIDERS = ['bedrock', 'openai', 'ollama']

const parseBool = (value: string) => value === 'true' || value === '1'

const SourceBadge = ({ source }: { source: ConfigEntry['source'] }) => {
  const tone =
    source === 'database'
      ? 'text-primary border-primary/40'
      : source === 'env'
        ? 'text-cyan-400 border-cyan-400/40'
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
    <Card>
      <CardHeader className='flex flex-wrap items-center justify-between gap-3'>
        <CardTitle className='micro-label'>CONFIGURATION</CardTitle>
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
                  <p className='micro-label'>{group.title}</p>
                  <div className='flex flex-col gap-3'>
                    {groupEntries.map(e => (
                      <div key={e.key} className='grid items-center gap-2 sm:grid-cols-[minmax(0,18rem)_1fr]'>
                        <div className='flex min-w-0 items-center gap-2'>
                          <span className='truncate font-mono text-xs' title={e.key}>
                            {e.label || e.key}
                          </span>
                          <SourceBadge source={e.source} />
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

/* -------------------------------- MCP SERVERS ------------------------------- */

const TOOL_PREVIEW_COUNT = 12

const MCPServerCard = ({ server }: { server: MCPServerInfo }) => {
  const [expanded, setExpanded] = useState(false)

  const tools = server.tools ?? []
  const visible = expanded ? tools : tools.slice(0, TOOL_PREVIEW_COUNT)

  return (
    <Card className='hud-corners gap-2 py-4'>
      <CardHeader className='px-4'>
        <CardTitle className='flex items-center justify-between font-mono text-sm tracking-widest'>
          {server.name}
          <span className='micro-label border-primary/40 text-primary rounded-sm border px-1.5 py-0.5'>
            {tools.length} TOOLS
          </span>
        </CardTitle>
      </CardHeader>
      <CardContent className='flex flex-col gap-2 px-4'>
        <ul className='text-muted-foreground flex flex-col gap-1 font-mono text-xs break-all'>
          {visible.map(tool => (
            <li key={tool}>▸ {tool}</li>
          ))}
        </ul>
        {tools.length > TOOL_PREVIEW_COUNT && (
          <button
            className='text-primary micro-label mt-1 self-start hover:underline'
            onClick={() => setExpanded(prev => !prev)}
          >
            {expanded ? '[ SHOW LESS ]' : `[ SHOW ALL ${tools.length} ]`}
          </button>
        )}
      </CardContent>
    </Card>
  )
}

const MCPPanel = () => {
  const [servers, setServers] = useState<MCPServerInfo[] | null>(null)
  const [a2a, setA2a] = useState<{ model?: string; notes?: string[] } | null>(null)
  const [skillStats, setSkillStats] = useState<Record<string, number> | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let mounted = true

    getMCPServers()
      .then(res => {
        if (!mounted) return
        setServers(res.servers ?? [])
        setA2a(res.agent_to_agent ?? null)
        setSkillStats(res.skill_stats ?? null)
      })
      .catch(err => mounted && setError(err instanceof Error ? err.message : String(err)))

    return () => {
      mounted = false
    }
  }, [])

  return (
    <div className='flex flex-col gap-4'>
      <div>
        <h2 className='font-mono text-sm font-semibold tracking-widest'>MCP + AGENT-TO-AGENT</h2>
        <p className='micro-label mt-1'>
          STDIO MCP (ARSENAL / STRIKE) · IN-PROCESS SKILLS & FINDINGS · ORCHESTRATOR DELEGATES
        </p>
      </div>
      {a2a && (
        <Card className='border-primary/30'>
          <CardHeader className='pb-2'>
            <CardTitle className='micro-label'>AGENT COMMUNICATION</CardTitle>
            <CardDescription className='font-mono text-xs'>{a2a.model}</CardDescription>
          </CardHeader>
          <CardContent className='space-y-1 font-mono text-[11px]'>
            {(a2a.notes ?? []).map((n, i) => (
              <p key={i} className='text-muted-foreground'>
                ▸ {n}
              </p>
            ))}
            {skillStats && (
              <p className='text-primary mt-2'>
                CyberStrike skills loaded: {skillStats.total ?? 0} (disk {skillStats.src_disk ?? 0} · builtin{' '}
                {skillStats.src_builtin ?? 0}) — agents call skill_search / skill_get
              </p>
            )}
          </CardContent>
        </Card>
      )}
      {error ? (
        <p className='text-destructive font-mono text-xs tracking-widest uppercase'>MCP SERVERS UNAVAILABLE — {error}</p>
      ) : servers === null ? (
        <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-3'>
          <Skeleton className='h-28 w-full' />
          <Skeleton className='h-28 w-full' />
        </div>
      ) : servers.length === 0 ? (
        <p className='micro-label py-6 text-center'>NO MCP SERVERS REGISTERED</p>
      ) : (
        <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-3'>
          {servers.map(server => (
            <MCPServerCard key={server.name} server={server} />
          ))}
        </div>
      )}
    </div>
  )
}

/* ------------------------------ INTELLIGENCE ------------------------------- */

const IntelligencePanel = () => {
  const [skillsCount, setSkillsCount] = useState<number | null>(null)
  const [agentsCount, setAgentsCount] = useState<number | null>(null)

  useEffect(() => {
    getSkills({ brief: true, limit: 1 })
      .then(r => setSkillsCount(r.total ?? r.count ?? 0))
      .catch(() => setSkillsCount(0))
    getAgents()
      .then(r => setAgentsCount(r.count ?? 0))
      .catch(() => setAgentsCount(0))
  }, [])

  return (
    <div className='flex flex-col gap-4'>
      <div>
        <h2 className='font-mono text-sm font-semibold tracking-widest'>INTELLIGENCE LAYER</h2>
        <p className='micro-label mt-1'>SKILLS · AGENTS · FINDINGS — FULL API ↔ UI WIRING</p>
      </div>
      <div className='grid gap-4 sm:grid-cols-3'>
        <Card className='hud-corners'>
          <CardHeader>
            <CardTitle className='micro-label'>SKILLS</CardTitle>
            <CardDescription className='font-mono text-xs'>
              {skillsCount === null ? '…' : `${skillsCount} loaded`} — methodology pack injected into agents
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link
              href='/skills'
              className={cn(
                buttonVariants({ variant: 'outline', size: 'sm' }),
                'font-mono text-[10px] tracking-widest uppercase'
              )}
            >
              Open Skills
            </Link>
          </CardContent>
        </Card>
        <Card className='hud-corners'>
          <CardHeader>
            <CardTitle className='micro-label'>AGENTS</CardTitle>
            <CardDescription className='font-mono text-xs'>
              {agentsCount === null ? '…' : `${agentsCount} modes`} — full / recon / web / network / exploit / post
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Link
              href='/agents'
              className={cn(
                buttonVariants({ variant: 'outline', size: 'sm' }),
                'font-mono text-[10px] tracking-widest uppercase'
              )}
            >
              Open Agents
            </Link>
          </CardContent>
        </Card>
        <Card className='hud-corners'>
          <CardHeader>
            <CardTitle className='micro-label'>FINDINGS</CardTitle>
            <CardDescription className='font-mono text-xs'>Global 3-gate findings registry across runs</CardDescription>
          </CardHeader>
          <CardContent>
            <Link
              href='/findings'
              className={cn(
                buttonVariants({ variant: 'outline', size: 'sm' }),
                'font-mono text-[10px] tracking-widest uppercase'
              )}
            >
              Open Findings
            </Link>
          </CardContent>
        </Card>
      </div>
    </div>
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

  return (
    <div className='flex flex-col gap-6'>
      <div className='flex flex-wrap items-end justify-between gap-2'>
        <div>
          <h1 className='font-mono text-xl font-semibold tracking-widest'>SYSTEM</h1>
          <p className='micro-label mt-1'>LIVE SERVICE HEALTH — PROBED SERVER-SIDE EVERY 10S</p>
        </div>
        {lastProbe && (
          <p className='micro-label'>
            LAST PROBE {lastProbe.toLocaleTimeString('en-GB', { hour12: false })}
          </p>
        )}
      </div>

      {probeError && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs tracking-widest uppercase'>
          CORE UNREACHABLE — {probeError}
        </div>
      )}

      <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-3'>
        {services === null
          ? Array.from({ length: SERVICE_SLOTS }).map((_, i) => <Skeleton key={i} className='h-28 w-full' />)
          : services.map(svc => <ServiceCard key={svc.name} svc={svc} />)}
      </div>

      <ConfigPanel />

      <IntelligencePanel />

      <MCPPanel />

      <div className='grid gap-4 lg:grid-cols-2'>
        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>DATA PERSISTENCE</CardTitle>
          </CardHeader>
          <CardContent className='flex flex-col gap-2 font-mono text-xs'>
            <p className='text-muted-foreground'>
              Run state, tool logs and reports are persisted by talon-core under <span className='text-foreground'>TALON_DATA_DIR</span>.
            </p>
            <p className='text-muted-foreground'>Back up that directory to retain operation history.</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>VERSION / BUILD</CardTitle>
          </CardHeader>
          <CardContent className='flex flex-col gap-2 font-mono text-xs'>
            <p>
              <span className='text-muted-foreground'>console:</span> talon-console (Next.js 16.2.12 / standalone)
            </p>
            <p>
              <span className='text-muted-foreground'>core:</span>{' '}
              {coreOnline ? 'reachable — version not exposed by API' : 'unknown (core offline)'}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>API PROXY</CardTitle>
          </CardHeader>
          <CardContent className='flex flex-col gap-2 font-mono text-xs'>
            <p className='text-muted-foreground'>
              All browser traffic is proxied through <span className='text-foreground'>/api/talon/*</span> →{' '}
              <span className='text-foreground'>TALON_CORE_URL</span> (default http://localhost:8000).
            </p>
            <p className='text-muted-foreground'>Run streams prefer WebSocket, degrade to SSE then polling.</p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default Settings
