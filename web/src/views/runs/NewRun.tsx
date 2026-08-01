'use client'

// React Imports
import { useEffect, useState } from 'react'

// Next Imports
import { useRouter, useSearchParams } from 'next/navigation'

// Third-party Imports
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { toast } from 'sonner'
import { z } from 'zod'

// Component Imports
import { TalonGlobe } from '@/components/shared/three'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'

// Util Imports
import type { AgentInfo, Playbook } from '@/lib/api'
import { getAgents, getPlaybooks, startRun } from '@/lib/api'

const schema = z.object({
  ip: z
    .string()
    .min(1, 'Target IP is required')
    .regex(/^(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}$/, 'Must be a valid IPv4 address'),
  cve_id: z
    .string()
    .regex(/^CVE-\d{4}-\d+$/, 'Format: CVE-2024-1234')
    .optional()
    .or(z.literal('')),
  service_name: z.string().optional(),
  description: z.string().optional(),
  lhost: z.string().optional(),
  lport: z
    .union([z.literal(''), z.coerce.number().int().min(1).max(65535)])
    .optional(),
  agent_mode: z.string().optional()
})

type FormValues = z.infer<typeof schema>

const Field = ({
  label,
  error,
  children
}: {
  label: string
  error?: string
  children: React.ReactNode
}) => (
  <div className='flex flex-col gap-1.5'>
    <Label className='micro-label'>{label}</Label>
    {children}
    {error && <p className='text-destructive font-mono text-[11px]'>{error}</p>}
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
    const ipQ = searchParams.get('ip')
    if (ipQ) {
      // prefill via form default is hard; set through setValue if available — use reset on ip field via DOM
    }
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
    if (pb) {
      setAgentMode(pb.agent_mode || 'full')
    }
  }, [playbookId, playbooks])

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors, isSubmitting }
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { ip: '', cve_id: '', service_name: '', description: '', lhost: '', lport: '', agent_mode: 'full' }
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

  return (
    <div className='mx-auto flex max-w-5xl flex-col gap-6'>
      <div>
        <h1 className='font-mono text-xl font-semibold tracking-widest'>NEW OPERATION</h1>
        <p className='micro-label mt-1'>PROVISION A PENTEST RUN · THREE.JS TARGET HUD</p>
      </div>

      <div className='grid gap-6 lg:grid-cols-[1fr_240px]'>
      <Card className='hud-corners scanlines relative overflow-hidden'>
        <CardHeader>
          <CardTitle className='text-primary font-mono text-sm tracking-widest'>$ talon run start</CardTitle>
          <CardDescription className='font-mono text-xs'>
            Target is queued for autonomous reconnaissance → exploitation → judge verification.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className='flex flex-col gap-4' noValidate>
            <Field label='TARGET IP *' error={errors.ip?.message}>
              <Input {...register('ip')} placeholder='10.10.10.5' className='font-mono' autoFocus />
            </Field>

            {playbooks.length > 0 && (
              <Field label='PLAYBOOK (OPTIONAL)'>
                <Select
                  value={playbookId || '__none__'}
                  onValueChange={v => setPlaybookId(!v || v === '__none__' ? '' : v)}
                >
                  <SelectTrigger className='font-mono'>
                    <SelectValue placeholder='None' />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value='__none__' className='font-mono text-xs'>
                      None — custom
                    </SelectItem>
                    {playbooks.map(pb => (
                      <SelectItem key={pb.id} value={pb.id} className='font-mono text-xs'>
                        {pb.codename} — {pb.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </Field>
            )}

            <Field label='AGENT MODE (SPECIALIST)'>
              <Select value={agentMode} onValueChange={v => setAgentMode(v ?? 'full')}>
                <SelectTrigger className='font-mono'>
                  <SelectValue placeholder='full' />
                </SelectTrigger>
                <SelectContent>
                  {(agents.length
                    ? agents
                    : [
                        { id: 'full', name: 'Full Pipeline', codename: 'COMMANDER' },
                        { id: 'recon', name: 'Recon', codename: 'GHOST' },
                        { id: 'web', name: 'Web Application', codename: 'STRIKER-WEB' },
                        { id: 'network', name: 'Internal Network', codename: 'PHANTOM' },
                        { id: 'exploit', name: 'Exploit', codename: 'STRIKER' },
                        { id: 'post', name: 'Post-Exploit', codename: 'CIPHER' }
                      ]
                  ).map(a => (
                    <SelectItem key={a.id} value={a.id} className='font-mono text-xs'>
                      {a.codename} — {a.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </Field>

            <div className='grid gap-4 sm:grid-cols-2'>
              <Field label='CVE ID' error={errors.cve_id?.message}>
                <Input {...register('cve_id')} placeholder='CVE-2021-41773' className='font-mono' />
              </Field>
              <Field label='SERVICE NAME' error={errors.service_name?.message}>
                <Input {...register('service_name')} placeholder='apache' className='font-mono' />
              </Field>
            </div>

            <Field label='DESCRIPTION' error={errors.description?.message}>
              <Textarea
                {...register('description')}
                placeholder='Optional operator notes / engagement context…'
                className='min-h-24 font-mono'
              />
            </Field>

            <div className='grid gap-4 sm:grid-cols-2'>
              <Field label='LHOST' error={errors.lhost?.message}>
                <Input {...register('lhost')} placeholder='reverse-shell listener host' className='font-mono' />
              </Field>
              <Field label='LPORT' error={errors.lport?.message}>
                <Input {...register('lport')} placeholder='4444' inputMode='numeric' className='font-mono' />
              </Field>
            </div>

            <Button
              type='submit'
              disabled={isSubmitting}
              className='mt-2 font-mono text-xs font-semibold tracking-widest uppercase'
            >
              {isSubmitting ? 'LAUNCHING…' : '[ EXECUTE OPERATION ]'}
            </Button>
          </form>
        </CardContent>
      </Card>

      <Card className='hud-corners hidden overflow-hidden lg:block'>
        <CardHeader className='pb-2'>
          <CardTitle className='micro-label'>TARGET HUD</CardTitle>
          <CardDescription className='font-mono text-[10px]'>Three.js · compact globe</CardDescription>
        </CardHeader>
        <CardContent className='flex justify-center p-2 pb-4'>
          <TalonGlobe className='aspect-square w-full max-w-[220px]' variant='compact' state='thinking' activityLevel={0.5} />
        </CardContent>
      </Card>
      </div>
    </div>
  )
}

export default NewRun
