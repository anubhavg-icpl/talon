'use client'

import { useCallback, useEffect, useState } from 'react'

import { toast } from 'sonner'

import type { Collaborator, Engagement, EngagementRole, ShareLink } from '@/lib/api'
import {
  createEngagement,
  createEngagementShare,
  listCollaborators,
  listEngagementShares,
  listEngagements,
  removeCollaborator,
  revokeEngagementShare
} from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'

const roleVariant = (
  role: EngagementRole | string
): 'destructive' | 'default' | 'secondary' | 'outline' => {
  switch (role) {
    case 'owner':
      return 'destructive'
    case 'build':
      return 'default'
    case 'use':
      return 'secondary'
    default:
      return 'outline'
  }
}

const ROLES: EngagementRole[] = ['owner', 'build', 'use']

const EngagementsView = () => {
  const [items, setItems] = useState<Engagement[] | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [newName, setNewName] = useState('')
  const [creating, setCreating] = useState(false)
  const [expanded, setExpanded] = useState<Set<string>>(new Set())
  const [shares, setShares] = useState<Record<string, ShareLink[]>>({})
  const [collabs, setCollabs] = useState<Record<string, Collaborator[]>>({})
  const [shareRole, setShareRole] = useState<Record<string, EngagementRole>>({})

  const reload = useCallback(async () => {
    try {
      const res = await listEngagements()
      setItems(res ?? [])
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err))
      setItems([])
    }
  }, [])

  useEffect(() => {
    void reload()
  }, [reload])

  const toggleExpand = async (id: string) => {
    const next = new Set(expanded)
    if (next.has(id)) {
      next.delete(id)
      setExpanded(next)
      return
    }
    next.add(id)
    setExpanded(next)
    try {
      const [sh, co] = await Promise.all([
        listEngagementShares(id),
        listCollaborators(id)
      ])
      setShares(prev => ({ ...prev, [id]: sh ?? [] }))
      setCollabs(prev => ({ ...prev, [id]: co ?? [] }))
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'failed to load details')
    }
  }

  const onCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newName.trim()) return
    setCreating(true)
    try {
      await createEngagement({ name: newName.trim() })
      setNewName('')
      toast.success('Engagement created')
      await reload()
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'create failed')
    } finally {
      setCreating(false)
    }
  }

  const onCreateShare = async (id: string) => {
    const role = shareRole[id] ?? 'use'
    try {
      const link = await createEngagementShare(id, role)
      setShares(prev => ({ ...prev, [id]: [...(prev[id] ?? []), link] }))
      toast.success('Share link created')
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'share failed')
    }
  }

  const onRevokeShare = async (id: string, linkId: string) => {
    try {
      await revokeEngagementShare(id, linkId)
      setShares(prev => ({
        ...prev,
        [id]: (prev[id] ?? []).map(s => (s.id === linkId ? { ...s, revoked: true } : s))
      }))
      toast.success('Share link revoked')
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'revoke failed')
    }
  }

  const onRemoveCollab = async (id: string, userId: string) => {
    try {
      await removeCollaborator(id, userId)
      setCollabs(prev => ({ ...prev, [id]: (prev[id] ?? []).filter(c => c.user_id !== userId) }))
      toast.success('Collaborator removed')
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'remove failed')
    }
  }

  const copyToken = (token: string) => {
    void navigator.clipboard.writeText(token).then(() => toast.success('Token copied'))
  }

  return (
    <div className='flex flex-col gap-6'>
      <div>
        <h1 className='font-mono text-xl font-semibold tracking-widest'>ENGAGEMENTS</h1>
        <p className='micro-label mt-1'>SHARED PENTEST SCOPES — COLLABORATION, SHARE LINKS & TEAM ACCESS</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className='micro-label'>NEW ENGAGEMENT</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onCreate} className='flex flex-wrap gap-2'>
            <Input
              placeholder='Engagement name'
              className='font-mono text-xs'
              value={newName}
              onChange={e => setNewName(e.target.value)}
            />
            <Button type='submit' size='sm' disabled={creating} className='font-mono text-[10px] tracking-widest uppercase'>
              {creating ? 'Creating…' : 'Create'}
            </Button>
          </form>
        </CardContent>
      </Card>

      {error && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs'>
          {error}
        </div>
      )}

      {!items ? (
        <Skeleton className='h-40 w-full' />
      ) : items.length === 0 ? (
        <p className='micro-label py-12 text-center'>NO ENGAGEMENTS — CREATE ONE ABOVE</p>
      ) : (
        <div className='flex flex-col gap-3'>
          {items.map(eng => {
            const isOpen = expanded.has(eng.id)
            const engShares = shares[eng.id]
            const engCollabs = collabs[eng.id]
            return (
              <Card key={eng.id}>
                <CardHeader className='pb-2'>
                  <div className='flex flex-wrap items-center gap-2'>
                    <Badge variant='secondary' className='font-mono text-[9px] uppercase'>
                      {eng.run_ids?.length ?? 0} runs
                    </Badge>
                    <span className='text-muted-foreground font-mono text-[10px]'>
                      {new Date(eng.created_at).toLocaleDateString()}
                    </span>
                  </div>
                  <button type='button' onClick={() => toggleExpand(eng.id)} className='mt-1 flex items-center gap-2 text-left'>
                    <CardTitle className='font-mono text-sm tracking-widest'>{eng.name}</CardTitle>
                    <span className='text-muted-foreground font-mono text-[10px]'>{isOpen ? '▲' : '▼'}</span>
                  </button>
                </CardHeader>
                <CardContent className='space-y-4'>
                  <div className='font-mono text-[11px]'>
                    <span className='text-muted-foreground'>ID: </span>
                    <span>{eng.id}</span>
                  </div>

                  {isOpen && (
                    <>
                      {/* Share links */}
                      <div className='space-y-2'>
                        <div className='micro-label'>SHARE LINKS</div>
                        <div className='flex flex-wrap items-center gap-2'>
                          <Label className='font-mono text-[10px] uppercase'>Role</Label>
                          <select
                            value={shareRole[eng.id] ?? 'use'}
                            onChange={e => setShareRole(prev => ({ ...prev, [eng.id]: e.target.value as EngagementRole }))}
                            className='rounded border border-border bg-background px-2 py-1 font-mono text-[11px]'
                          >
                            {ROLES.map(r => (
                              <option key={r} value={r}>{r}</option>
                            ))}
                          </select>
                          <Button
                            size='sm'
                            variant='outline'
                            onClick={() => onCreateShare(eng.id)}
                            className='font-mono text-[10px] tracking-widest uppercase'
                          >
                            Generate link
                          </Button>
                        </div>
                        {!engShares ? (
                          <p className='text-muted-foreground font-mono text-[11px]'>Loading…</p>
                        ) : engShares.length === 0 ? (
                          <p className='text-muted-foreground font-mono text-[11px]'>No share links</p>
                        ) : (
                          <div className='space-y-1'>
                            {engShares.map(s => (
                              <div key={s.id} className='flex flex-wrap items-center gap-2 border-b border-border/40 py-1 font-mono text-[11px]'>
                                <Badge variant={roleVariant(s.role)} className='font-mono text-[9px] uppercase'>{s.role}</Badge>
                                <span className='text-muted-foreground truncate'>{s.label || s.token.slice(0, 16) + '…'}</span>
                                {s.revoked && <Badge variant='destructive' className='font-mono text-[9px]'>REVOKED</Badge>}
                                {!s.revoked && (
                                  <>
                                    <button type='button' onClick={() => copyToken(s.token)} className='text-primary underline'>
                                      copy
                                    </button>
                                    <button type='button' onClick={() => onRevokeShare(eng.id, s.id)} className='text-destructive ml-auto underline'>
                                      revoke
                                    </button>
                                  </>
                                )}
                              </div>
                            ))}
                          </div>
                        )}
                      </div>

                      {/* Collaborators */}
                      <div className='space-y-2'>
                        <div className='micro-label'>COLLABORATORS</div>
                        {!engCollabs ? (
                          <p className='text-muted-foreground font-mono text-[11px]'>Loading…</p>
                        ) : engCollabs.length === 0 ? (
                          <p className='text-muted-foreground font-mono text-[11px]'>No collaborators</p>
                        ) : (
                          <div className='space-y-1'>
                            {engCollabs.map(c => (
                              <div key={c.user_id} className='flex items-center gap-2 border-b border-border/40 py-1 font-mono text-[11px]'>
                                <Badge variant={roleVariant(c.role)} className='font-mono text-[9px] uppercase'>{c.role}</Badge>
                                <span>{c.username}</span>
                                <button
                                  type='button'
                                  onClick={() => onRemoveCollab(eng.id, c.user_id)}
                                  className='text-destructive ml-auto underline'
                                >
                                  remove
                                </button>
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                    </>
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

export default EngagementsView
