'use client'

// React Imports
import { useEffect, useState } from 'react'

// Next Imports
import { useRouter, useSearchParams } from 'next/navigation'

// Third-party Imports
import { zodResolver } from '@hookform/resolvers/zod'
import { Crosshair, Rocket } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { toast } from 'sonner'
import { z } from 'zod'

// Component Imports
import GlobeWallpaper from '@/components/shared/GlobeWallpaper'
import PageHeader from '@/components/shared/PageHeader'
import { TalonGlobe } from '@/components/shared/three'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'

// Util Imports
import type { AgentInfo, Playbook } from '@/lib/api'
import { agentAvatarSrc } from '@/lib/agent-avatars'
import { getAgents, getPlaybooks, startRun } from '@/lib/api'
import { playbookArtSrc } from '@/lib/playbook-art'
import { cn } from '@/lib/utils'

const schema = z.object({
  ip: z
    .string()
    .min(1, 'Target IP is required')
    .regex(
      /^(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}$/,
      'Must be a valid IPv4 address'
    ),
  cve_id: z
    .string()
    .regex(/^CVE-\d{4}-\d+$/, 'Format: CVE-2024-1234')
    .optional()
    .or(z.literal('')),
  service_name: z.string().optional(),
  description: z.string().optional(),
  lhost: z.string().optional(),
  lport: z.union([z.literal(''), z.coerce.number().int().min(1).max(65535)]).optional(),
  agent_mode: z.string().optional()
})

type FormValues = z.infer<typeof schema>

const FALLBACK_AGENTS: AgentInfo[] = [
  { id: 'full', name: 'Full Pipeline', codename: 'COMMANDER', focus: 'general', description: '', delegates: [] },
  { id: 'recon', name: 'Recon', codename: 'GHOST', focus: 'recon', description: '', delegates: [] },
  { id: 'web', name: 'Web Application', codename: 'STRIKER-WEB', focus: 'web', description: '', delegates: [] },
  { id: 'network', name: 'Internal Network', codename: 'PHANTOM', focus: 'network', description: '', delegates: [] },
  { id: 'exploit', name: 'Exploit', codename: 'STRIKER', focus: 'exploit', description: '', delegates: [] },
  { id: 'post', name: 'Post-Exploit', codename: 'CIPHER', focus: 'post', description: '', delegates: [] }
]

const Field = ({
  label,
  error,
  children,
  className
}: {
  label: string
  error?: string
  children: React.ReactNode
  className?: string
}) => (
  <div className={cn('flex min-w-0 flex-col gap-1.5', className)}>
    <Label className='micro-label text-muted-foreground'>{label}</Label>
    {children}
    <p className={cn('min-h-4 font-mono text-[11px]', error ? 'text-destructive' : 'text-transparent')}>
      {error || '·'}
    </p>
  </div>
)

const NewRun = () => {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [agents, setAgents] = useState<AgentInfo[]>([])
  const [playbooks, setPlaybooks] = useState<Playbook[]>([])
  const [agentMode, setAgentMode] = useState('full')
  const [playbookId, setPlaybookId] = useState('')

  useEffect(() => {
    const mode = searchParams.get('mode')
    if (mode) setAgentMode(mode)
    const pb = searchParams.get('playbook')
    if (pb) setPlaybookId(pb)
    getAgents()
      .then(res => setAgents(res.agents ?? []))
      .catch(() => {})
    getPlaybooks()
      .then(res => setPlaybooks(res.playbooks ?? []))
      .catch(() => {})
  }, [searchParams])

  useEffect(() => {
    if (!playbookId) return
    const pb = playbooks.find(p => p.id === playbookId)
    if (pb) setAgentMode(pb.agent_mode || 'full')
  }, [playbookId, playbooks])

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors, isSubmitting }
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      ip: '',
      cve_id: '',
      service_name: '',
      description: '',
      lhost: '',
      lport: '',
      agent_mode: 'full'
    }
  })

  useEffect(() => {
    const ipQ = searchParams.get('ip')
    if (ipQ) setValue('ip', ipQ)
  }, [searchParams, setValue])

  const onSubmit = async (values: FormValues) => {
    try {
      const res = await startRun({
        ip: values.ip,
        ...(values.cve_id ? { cve_id: values.cve_id } : {}),
        ...(values.service_name ? { service_name: values.service_name } : {}),
        ...(values.description ? { description: values.description } : {}),
        ...(values.lhost ? { lhost: values.lhost } : {}),
        ...(values.lport ? { lport: Number(values.lport) } : {}),
        agent_mode: agentMode || 'full',
        ...(playbookId ? { playbook_id: playbookId } : {})
      })

      toast.success(`Operation launched — ${res.run_id}`)
      router.push(`/runs/${res.run_id}`)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to start operation')
    }
  }

  const agentList = agents.length ? agents : FALLBACK_AGENTS
  const selectedAgent = agentList.find(a => a.id === agentMode)
  const selectedPlaybook = playbooks.find(p => p.id === playbookId)

  return (
    <div className='flex w-full flex-col gap-6'>
      <PageHeader
        title={
          <span className='inline-flex items-center gap-2.5'>
            <Crosshair className='text-primary size-6 shrink-0' />
            NEW OPERATION
          </span>
        }
        subtitle='PROVISION A PENTEST RUN · TARGET HUD · AGENT PIPELINE'
      />

      <div className='grid w-full items-start gap-6 lg:grid-cols-[minmax(0,1fr)_minmax(260px,300px)]'>
        {/* Form column */}
        <Card className='hud-corners scanlines relative min-w-0 overflow-hidden'>
          <CardHeader className='space-y-1 border-b border-primary/10 pb-4'>
            <CardTitle className='text-primary font-mono text-sm tracking-widest'>$ talon run start</CardTitle>
            <CardDescription className='font-mono text-xs leading-relaxed'>
              Target is queued for autonomous reconnaissance → exploitation → judge verification.
            </CardDescription>
          </CardHeader>
          <CardContent className='pt-5'>
            <form onSubmit={handleSubmit(onSubmit)} className='flex flex-col gap-1' noValidate>
              <Field label='TARGET IP *' error={errors.ip?.message}>
                <Input {...register('ip')} placeholder='10.10.10.5' className='h-10 font-mono' autoFocus />
              </Field>

              <div className='grid gap-x-4 sm:grid-cols-2'>
                {playbooks.length > 0 && (
                  <Field label='PLAYBOOK (OPTIONAL)'>
                    <Select
                      value={playbookId || '__none__'}
                      onValueChange={v => setPlaybookId(!v || v === '__none__' ? '' : v)}
                    >
                      <SelectTrigger className='h-10 w-full font-mono'>
                        <SelectValue placeholder='None' />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value='__none__' className='font-mono text-xs'>
                          None — custom
                        </SelectItem>
                        {playbooks.map(pb => (
                          <SelectItem key={pb.id} value={pb.id} className='font-mono text-xs'>
                            <span className='flex min-w-0 items-center gap-2'>
                              {/* eslint-disable-next-line @next/next/no-img-element */}
                              <img
                                src={playbookArtSrc(pb.id)}
                                alt=''
                                className='size-5 shrink-0 rounded-sm object-cover ring-1 ring-primary/30'
                              />
                              <span className='truncate'>
                                {pb.codename} — {pb.name}
                              </span>
                            </span>
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </Field>
                )}

                <Field label='AGENT MODE (SPECIALIST)' className={playbooks.length === 0 ? 'sm:col-span-2' : undefined}>
                  <Select value={agentMode} onValueChange={v => setAgentMode(v ?? 'full')}>
                    <SelectTrigger className='h-10 w-full font-mono'>
                      <SelectValue placeholder='full' />
                    </SelectTrigger>
                    <SelectContent>
                      {agentList.map(a => (
                        <SelectItem key={a.id} value={a.id} className='font-mono text-xs'>
                          <span className='flex min-w-0 items-center gap-2'>
                            {/* eslint-disable-next-line @next/next/no-img-element */}
                            <img
                              src={agentAvatarSrc(a.id)}
                              alt=''
                              className='size-5 shrink-0 rounded-full object-cover ring-1 ring-primary/30'
                            />
                            <span className='truncate'>
                              {a.codename} — {a.name}
                            </span>
                          </span>
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </Field>
              </div>

              <div className='grid gap-x-4 sm:grid-cols-2'>
                <Field label='CVE ID' error={errors.cve_id?.message}>
                  <Input {...register('cve_id')} placeholder='CVE-2021-41773' className='h-10 font-mono' />
                </Field>
                <Field label='SERVICE NAME' error={errors.service_name?.message}>
                  <Input {...register('service_name')} placeholder='apache' className='h-10 font-mono' />
                </Field>
              </div>

              <Field label='DESCRIPTION' error={errors.description?.message}>
                <Textarea
                  {...register('description')}
                  placeholder='Optional operator notes / engagement context…'
                  className='min-h-24 resize-y font-mono'
                />
              </Field>

              <div className='grid gap-x-4 sm:grid-cols-2'>
                <Field label='LHOST' error={errors.lhost?.message}>
                  <Input
                    {...register('lhost')}
                    placeholder='reverse-shell listener host'
                    className='h-10 font-mono'
                  />
                </Field>
                <Field label='LPORT' error={errors.lport?.message}>
                  <Input
                    {...register('lport')}
                    placeholder='4444'
                    inputMode='numeric'
                    className='h-10 font-mono'
                  />
                </Field>
              </div>

              <div className='mt-2 flex flex-col gap-3 border-t border-primary/10 pt-4 sm:flex-row sm:items-center sm:justify-between'>
                <p className='micro-label text-muted-foreground order-2 sm:order-1'>
                  {selectedAgent ? (
                    <>
                      MODE · <span className='text-primary'>{selectedAgent.codename}</span>
                      {selectedPlaybook ? (
                        <>
                          {' '}
                          · PB · <span className='text-primary'>{selectedPlaybook.codename}</span>
                        </>
                      ) : null}
                    </>
                  ) : (
                    'AUTHORIZED TARGETS ONLY'
                  )}
                </p>
                <Button
                  type='submit'
                  disabled={isSubmitting}
                  className='glow-red order-1 h-10 w-full font-mono text-xs font-semibold tracking-widest uppercase sm:order-2 sm:w-auto sm:min-w-52'
                >
                  <Rocket className='mr-2 size-3.5' />
                  {isSubmitting ? 'LAUNCHING…' : '[ EXECUTE OPERATION ]'}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>

        {/* HUD column — sticky, aligned with form top */}
        <aside className='flex min-w-0 flex-col gap-4 lg:sticky lg:top-20'>
          <Card className='hud-corners overflow-hidden'>
            <CardHeader className='space-y-1 border-b border-primary/10 pb-3'>
              <CardTitle className='micro-label'>TARGET HUD</CardTitle>
              <CardDescription className='font-mono text-[10px] leading-relaxed'>
                Dark planet · optional Live 3D
              </CardDescription>
            </CardHeader>
            <CardContent className='p-3'>
              <div className='mx-auto aspect-square w-full max-w-[280px]'>
                <GlobeWallpaper
                  compact
                  label='TARGET HUD'
                  className='size-full'
                  live3d={
                    <TalonGlobe
                      className='h-full w-full'
                      variant='compact'
                      state='thinking'
                      activityLevel={0.45}
                      interactive
                    />
                  }
                />
              </div>
            </CardContent>
          </Card>

          <Card className='hud-corners border-primary/15 bg-card/60'>
            <CardContent className='space-y-3 p-4 font-mono text-[10px] tracking-wide'>
              <div className='flex items-center justify-between gap-2'>
                <span className='text-muted-foreground uppercase'>Pipeline</span>
                <span className='text-primary'>recon → exploit → judge</span>
              </div>
              <div className='flex items-center justify-between gap-2'>
                <span className='text-muted-foreground uppercase'>Agent</span>
                <span className='truncate text-right uppercase'>{selectedAgent?.codename ?? 'COMMANDER'}</span>
              </div>
              <div className='flex items-center justify-between gap-2'>
                <span className='text-muted-foreground uppercase'>Playbook</span>
                <span className='truncate text-right uppercase'>{selectedPlaybook?.codename ?? 'CUSTOM'}</span>
              </div>
              <div className='border-t border-primary/10 pt-3 text-muted-foreground leading-relaxed'>
                Scope and HITL gates apply after launch. Monitor live progress on the run detail stream.
              </div>
            </CardContent>
          </Card>
        </aside>
      </div>
    </div>
  )
}

export default NewRun
