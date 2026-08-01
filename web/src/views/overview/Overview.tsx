'use client'

// React Imports
import { useEffect, useMemo, useState } from 'react'

// Next Imports
import Link from 'next/link'

// Third-party Imports
import { Area, AreaChart, CartesianGrid, Cell, Pie, PieChart, XAxis, YAxis } from 'recharts'

// Type Imports
import type { RunsSummaryResponse, RunSummary } from '@/lib/api'

// Component Imports
import Elapsed from '@/components/shared/Elapsed'
import LiveDot from '@/components/shared/LiveDot'
import Logo from '@/components/shared/Logo'
import PageHeader from '@/components/shared/PageHeader'
import StatusBadge from '@/components/shared/StatusBadge'
import { TalonGlobe } from '@/components/shared/three'
import UtcClock from '@/components/shared/UtcClock'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ChartContainer, ChartTooltip, ChartTooltipContent } from '@/components/ui/chart'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

// Util Imports
import { getGlobalFindings, getSkills, listRuns, runsSummary } from '@/lib/api'
import { relativeTime } from '@/lib/format'

const ACTIVE_STATUSES = new Set(['running', 'awaiting_approval', 'initializing'])

const VERDICT_COLORS = {
  compromised: 'var(--chart-1)',
  clean: 'var(--muted-foreground)',
  error: 'var(--chart-2)'
}

const Overview = () => {
  const [runs, setRuns] = useState<RunSummary[] | null>(null)
  const [summary, setSummary] = useState<RunsSummaryResponse | null>(null)
  const [findingsTotal, setFindingsTotal] = useState(0)
  const [skillsTotal, setSkillsTotal] = useState(0)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let mounted = true

    const load = () =>
      Promise.all([
        listRuns(10),
        runsSummary(),
        getGlobalFindings({ limit: 200 }).catch(() => ({ findings: [], count: 0 })),
        getSkills({ brief: true, limit: 1 }).catch(() => ({ skills: [], count: 0, total: 0 }))
      ])
        .then(([runsRes, summaryRes, findingsRes, skillsRes]) => {
          if (!mounted) return
          setRuns(runsRes.runs ?? [])
          setSummary(summaryRes)
          setFindingsTotal(findingsRes.count ?? findingsRes.findings?.length ?? 0)
          // API: total = catalog size; count = page size only
          setSkillsTotal(skillsRes.total ?? skillsRes.count ?? 0)
          setError(null)
        })
        .catch(err => mounted && setError(err instanceof Error ? err.message : String(err)))

    load()
    const id = setInterval(load, 5000)

    return () => {
      mounted = false
      clearInterval(id)
    }
  }, [])

  const stats = useMemo(
    () => ({
      total: summary?.total ?? 0,
      active: summary?.active ?? 0,
      compromised: summary?.compromised ?? 0,
      awaiting: summary?.awaiting_approval ?? 0,
      findings: findingsTotal,
      skills: skillsTotal
    }),
    [summary, findingsTotal, skillsTotal]
  )

  const activity = useMemo(() => {
    const buckets = new Map<string, number>()

    for (const run of runs ?? []) {
      const day = run.started_at.slice(0, 10)

      if (day) buckets.set(day, (buckets.get(day) ?? 0) + 1)
    }

    return [...buckets.entries()]
      .sort(([a], [b]) => a.localeCompare(b))
      .map(([day, count]) => ({ day: day.slice(5), runs: count }))
  }, [runs])

  const verdicts = useMemo(() => {
    const all = runs ?? []
    const compromised = all.filter(r => r.judge_verdict === true).length
    const errored = all.filter(r => r.status === 'error').length
    const clean = all.filter(r => r.status === 'completed' && r.judge_verdict !== true).length

    return [
      { name: 'COMPROMISED', value: compromised, fill: VERDICT_COLORS.compromised },
      { name: 'CLEAN', value: clean, fill: VERDICT_COLORS.clean },
      { name: 'ERROR', value: errored, fill: VERDICT_COLORS.error }
    ].filter(v => v.value > 0)
  }, [runs])

  const activeOps = useMemo(
    () => (runs ?? []).filter(r => ACTIVE_STATUSES.has(r.status)),
    [runs]
  )

  const recent = useMemo(
    () =>
      [...(runs ?? [])]
        .sort((a, b) => b.started_at.localeCompare(a.started_at))
        .slice(0, 8),
    [runs]
  )

  const statCards = [
    { label: 'TOTAL RUNS', value: stats.total, tone: 'text-foreground' },
    { label: 'ACTIVE OPS', value: stats.active, tone: 'text-primary' },
    { label: 'COMPROMISED', value: stats.compromised, tone: 'text-primary text-glow' },
    { label: 'AWAITING APPROVAL', value: stats.awaiting, tone: stats.awaiting > 0 ? 'text-warning' : 'text-foreground' },
    { label: 'FINDINGS', value: stats.findings, tone: 'text-primary' },
    { label: 'SKILLS LOADED', value: stats.skills, tone: 'text-foreground' }
  ]

  // Globe pulse: active ops drive activity; running state spins faster.
  const globeState = stats.active > 0 ? 'running' : 'idle'
  const globeLevel = Math.min(1, (stats.active || 0) * 0.35 + (stats.awaiting > 0 ? 0.25 : 0))

  return (
    <div className='flex flex-col gap-6'>
      <PageHeader
        title='OVERVIEW'
        description='Fleet status · live engagements · findings roll-up'
        actions={
          <Link
            href='/runs/new'
            className='bg-primary text-primary-foreground hover:bg-primary/90 rounded-sm px-3 py-2 font-mono text-[10px] tracking-widest uppercase'
          >
            New operation
          </Link>
        }
      />

      {/* Hero + single WebGL globe (prod: no matrix rain stack) */}
      <div className='hud-corners relative overflow-hidden rounded-sm border border-primary/15'>
        <div className='scanlines pointer-events-none absolute inset-0 opacity-60' />
        <div className='relative grid items-center gap-6 px-6 py-8 md:grid-cols-[1fr_minmax(200px,280px)] lg:grid-cols-[1fr_300px]'>
          <div className='flex flex-col gap-2'>
            <p className='micro-label'>AI PENTEST ORCHESTRATION</p>
            <Logo className='text-2xl sm:text-3xl [&_svg]:size-7' />
            <UtcClock className='text-muted-foreground font-mono text-xs tracking-widest' />
            <p className='text-muted-foreground mt-2 max-w-md font-mono text-[11px] leading-relaxed'>
              Operator globe tracks active runs · CyberStrike skills + multi-agent pipeline · drag to orbit
            </p>
            {stats.active > 0 && (
              <p className='text-primary micro-label mt-1 flex items-center gap-2'>
                <LiveDot tone='cyan' /> {stats.active} ACTIVE · C2 LIVE
              </p>
            )}
            <Link
              href='/showcase'
              className='text-primary micro-label mt-3 inline-flex w-fit tracking-widest underline-offset-4 hover:underline'
            >
              OPEN PRODUCT SHOWCASE →
            </Link>
          </div>
          <div className='relative mx-auto aspect-square w-full max-w-[280px] lg:max-w-[300px] overflow-hidden rounded-sm border border-primary/20 bg-black/40'>
            <TalonGlobe
              className='h-full w-full'
              variant='hero'
              interactive
              state={globeState}
              activityLevel={globeLevel}
              onClick={() => {
                window.location.assign('/runs/new')
              }}
            />
          </div>
        </div>
      </div>

      {error && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs tracking-widest uppercase'>
          CORE UNREACHABLE — {error}
        </div>
      )}

      {/* Stat cards */}
      <div className='grid grid-cols-2 gap-4 lg:grid-cols-3 xl:grid-cols-6'>
        {statCards.map(card => (
          <Card key={card.label} className='hud-corners gap-2 py-4'>
            <CardHeader className='px-4'>
              <CardTitle className='micro-label'>{card.label}</CardTitle>
            </CardHeader>
            <CardContent className='px-4'>
              {summary === null ? (
                <Skeleton className='h-9 w-16' />
              ) : (
                <p className={`font-mono text-4xl font-semibold ${card.tone}`}>{card.value}</p>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Charts */}
      <div className='grid gap-4 lg:grid-cols-3'>
        <Card className='lg:col-span-2'>
          <CardHeader>
            <CardTitle className='micro-label'>RUN ACTIVITY</CardTitle>
          </CardHeader>
          <CardContent>
            {runs === null ? (
              <Skeleton className='h-56 w-full' />
            ) : activity.length === 0 ? (
              <p className='micro-label flex h-56 items-center justify-center'>NO ACTIVITY RECORDED</p>
            ) : (
              <ChartContainer
                className='h-56 w-full'
                config={{ runs: { label: 'Runs', color: 'var(--chart-1)' } }}
              >
                <AreaChart data={activity} margin={{ left: -20, right: 8, top: 8 }}>
                  <CartesianGrid vertical={false} strokeDasharray='3 3' />
                  <XAxis dataKey='day' tickLine={false} axisLine={false} tickMargin={8} />
                  <YAxis allowDecimals={false} tickLine={false} axisLine={false} width={32} />
                  <ChartTooltip content={<ChartTooltipContent />} />
                  <Area
                    type='monotone'
                    dataKey='runs'
                    stroke='var(--chart-1)'
                    fill='var(--chart-1)'
                    fillOpacity={0.12}
                    strokeWidth={2}
                  />
                </AreaChart>
              </ChartContainer>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className='micro-label'>VERDICTS</CardTitle>
          </CardHeader>
          <CardContent>
            {runs === null ? (
              <Skeleton className='h-56 w-full' />
            ) : verdicts.length === 0 ? (
              <p className='micro-label flex h-56 items-center justify-center'>NO VERDICTS YET</p>
            ) : (
              <ChartContainer className='h-56 w-full' config={{ verdicts: { label: 'Verdicts' } }}>
                <PieChart>
                  <ChartTooltip content={<ChartTooltipContent hideLabel />} />
                  <Pie data={verdicts} dataKey='value' nameKey='name' innerRadius='55%' outerRadius='80%' strokeWidth={2}>
                    {verdicts.map(v => (
                      <Cell key={v.name} fill={v.fill} />
                    ))}
                  </Pie>
                </PieChart>
              </ChartContainer>
            )}
            {verdicts.length > 0 && (
              <div className='mt-2 flex justify-center gap-4'>
                {verdicts.map(v => (
                  <span key={v.name} className='micro-label flex items-center gap-1.5'>
                    <span className='inline-block size-2 rounded-full' style={{ background: v.fill }} />
                    {v.name} {v.value}
                  </span>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Active operations + recent runs */}
      <div className='grid gap-4 lg:grid-cols-3'>
        <Card className='lg:col-span-1'>
          <CardHeader>
            <CardTitle className='micro-label flex items-center gap-2'>
              <LiveDot tone='green' pulse={activeOps.length > 0} /> ACTIVE OPERATIONS
            </CardTitle>
          </CardHeader>
          <CardContent className='flex flex-col gap-3'>
            {runs === null ? (
              <>
                <Skeleton className='h-10 w-full' />
                <Skeleton className='h-10 w-full' />
              </>
            ) : activeOps.length === 0 ? (
              <p className='micro-label py-6 text-center'>NO ACTIVE OPERATIONS</p>
            ) : (
              activeOps.map(run => (
                <Link
                  key={run.run_id}
                  href={`/runs/${run.run_id}`}
                  className='hover:bg-muted/50 flex items-center justify-between gap-3 rounded-md border px-3 py-2 transition-colors'
                >
                  <span className='flex min-w-0 items-center gap-2'>
                    <LiveDot tone={run.status === 'awaiting_approval' ? 'amber' : 'green'} />
                    <span className='truncate font-mono text-sm'>{run.target}</span>
                  </span>
                  <Elapsed since={run.started_at} className='text-muted-foreground shrink-0 font-mono text-xs' />
                </Link>
              ))
            )}
          </CardContent>
        </Card>

        <Card className='lg:col-span-2'>
          <CardHeader>
            <CardTitle className='micro-label'>RECENT RUNS</CardTitle>
          </CardHeader>
          <CardContent>
            {runs === null ? (
              <Skeleton className='h-48 w-full' />
            ) : recent.length === 0 ? (
              <div className='flex h-48 flex-col items-center justify-center gap-3'>
                <p className='micro-label'>NO OPERATIONS YET</p>
                <Link href='/runs/new' className='text-primary font-mono text-xs tracking-widest uppercase hover:underline'>
                  [ + LAUNCH YOUR FIRST RUN ]
                </Link>
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className='micro-label'>STATUS</TableHead>
                    <TableHead className='micro-label'>TARGET</TableHead>
                    <TableHead className='micro-label hidden sm:table-cell'>CVE</TableHead>
                    <TableHead className='micro-label text-right'>STARTED</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {recent.map(run => (
                    <TableRow key={run.run_id} className='cursor-pointer'>
                      <TableCell>
                        <Link href={`/runs/${run.run_id}`} className='block'>
                          <StatusBadge status={run.status} />
                        </Link>
                      </TableCell>
                      <TableCell className='font-mono text-sm'>
                        <Link href={`/runs/${run.run_id}`} className='block'>
                          {run.target}
                        </Link>
                      </TableCell>
                      <TableCell className='text-muted-foreground hidden font-mono text-xs sm:table-cell'>
                        {run.cve_id ?? '—'}
                      </TableCell>
                      <TableCell className='text-muted-foreground text-right font-mono text-xs'>
                        {relativeTime(run.started_at)}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default Overview
