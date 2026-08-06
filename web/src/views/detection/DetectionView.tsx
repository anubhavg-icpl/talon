'use client'

import { useCallback, useEffect, useState } from 'react'
import { toast } from 'sonner'
import { getDetectionSkillsByType, getSkill, type DetectionSkill, type DetectionSkillType } from '@/lib/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Skeleton } from '@/components/ui/skeleton'

const TYPE_DESCRIPTIONS: Record<DetectionSkillType, string> = {
  triage: 'Decide escalate or dismiss on security alerts using structured checklists with majority-rule verdicts.',
  investigation: 'Confirm or rule out threats with multi-signal analysis (beacon confirmation, attribution, process scoping, operator activity).',
  tuning: 'Propose detection changes (exclude, include, modify, fork) after case closure to reduce noise and catch real threats.',
}

export default function DetectionView() {
  const [activeType, setActiveType] = useState<DetectionSkillType>('triage')
  const [skills, setSkills] = useState<DetectionSkill[] | null>(null)
  const [search, setSearch] = useState('')
  const [selected, setSelected] = useState<DetectionSkill | null>(null)
  const [skillBody, setSkillBody] = useState<string | null>(null)
  const [loadingBody, setLoadingBody] = useState(false)

  const loadSkills = useCallback(async (type: DetectionSkillType) => {
    setSkills(null)
    setSelected(null)
    setSkillBody(null)
    try {
      const res = await getDetectionSkillsByType(type)
      setSkills(res.skills ?? [])
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to load skills')
      setSkills([])
    }
  }, [])

  const viewSkillBody = useCallback(async (skillId: string) => {
    setLoadingBody(true)
    setSkillBody(null)
    try {
      const skill = await getSkill(skillId)
      setSkillBody(skill.body || 'No content available')
    } catch {
      setSkillBody('Failed to load skill content.')
    } finally {
      setLoadingBody(false)
    }
  }, [])

  useEffect(() => {
    void loadSkills(activeType)
  }, [activeType, loadSkills])

  const filtered = skills?.filter(s =>
    !search || s.name.toLowerCase().includes(search.toLowerCase())
  ) ?? []

  return (
    <div className='grid-bg min-h-screen p-6'>
      <div className='mx-auto max-w-6xl space-y-6'>
        <div>
          <h1 className='font-mono text-2xl font-bold tracking-widest uppercase text-cyan-400'>
            Alert Triage
          </h1>
          <p className='text-muted-foreground mt-1 font-mono text-xs'>
            SOC alert triage, investigation, and detection tuning pipeline. 50 built-in playbooks for structured alert analysis.
          </p>
        </div>

        <Card>
          <CardContent className='pt-6'>
            <p className='text-sm text-muted-foreground'>{TYPE_DESCRIPTIONS[activeType]}</p>
          </CardContent>
        </Card>

        <Tabs value={activeType} onValueChange={(v) => setActiveType(v as DetectionSkillType)}>
          <TabsList variant='line' className='w-full justify-start font-mono text-xs tracking-widest uppercase'>
            <TabsTrigger value='triage'>Triage (38)</TabsTrigger>
            <TabsTrigger value='investigation'>Investigation (8)</TabsTrigger>
            <TabsTrigger value='tuning'>Tuning (4)</TabsTrigger>
          </TabsList>

          {(['triage', 'investigation', 'tuning'] as DetectionSkillType[]).map(type => (
            <TabsContent key={type} value={type} className='pt-4'>
              <div className='flex gap-6'>
                {/* Skill list */}
                <div className='w-1/2 space-y-2'>
                  <Input
                    placeholder='Search skills...'
                    value={search}
                    onChange={e => setSearch(e.target.value)}
                    className='font-mono text-xs'
                  />
                  {!skills ? (
                    <Skeleton className='h-64 w-full' />
                  ) : filtered.length === 0 ? (
                    <p className='micro-label py-8 text-center'>NO SKILLS FOUND</p>
                  ) : (
                    <div className='max-h-[32rem] space-y-1 overflow-y-auto'>
                      {filtered.map(skill => (
                        <button
                          key={skill.id}
                          onClick={() => setSelected(skill)}
                          className={`w-full rounded border px-3 py-2 text-left font-mono text-xs transition-colors ${
                            selected?.id === skill.id
                              ? 'border-primary/40 bg-primary/10 text-primary'
                              : 'border-border/50 hover:border-border hover:bg-muted/50'
                          }`}
                        >
                          <div className='flex items-center justify-between'>
                            <span className='font-semibold'>{skill.name}</span>
                            <Badge variant='outline' className='text-[10px]'>{skill.stage}</Badge>
                          </div>
                          <div className='text-muted-foreground mt-1 text-[10px]'>{skill.id}</div>
                        </button>
                      ))}
                    </div>
                  )}
                </div>

                {/* Detail panel */}
                <div className='w-1/2'>
                  {!selected ? (
                    <Card>
                      <CardContent className='flex items-center justify-center py-16'>
                        <p className='micro-label text-center'>SELECT A SKILL TO VIEW DETAILS</p>
                      </CardContent>
                    </Card>
                  ) : (
                    <Card>
                      <CardHeader>
                        <CardTitle className='font-mono text-sm tracking-wider uppercase'>
                          {selected.name}
                        </CardTitle>
                      </CardHeader>
                      <CardContent className='space-y-3'>
                        <div className='grid grid-cols-2 gap-3 font-mono text-xs'>
                          <div>
                            <div className='micro-label text-[10px]'>ID</div>
                            <div className='text-muted-foreground'>{selected.id}</div>
                          </div>
                          <div>
                            <div className='micro-label text-[10px]'>CATEGORY</div>
                            <div className='text-muted-foreground'>{selected.category}</div>
                          </div>
                          <div>
                            <div className='micro-label text-[10px]'>STAGE</div>
                            <div className='text-muted-foreground'>{selected.stage}</div>
                          </div>
                          <div>
                            <div className='micro-label text-[10px]'>PATH</div>
                            <div className='text-muted-foreground truncate'>{selected.path}</div>
                          </div>
                        </div>
                        <div className='pt-2 flex gap-2'>
                          <Button
                            variant='outline'
                            size='sm'
                            className='font-mono text-xs'
                            onClick={() => {
                              toast.info(`Use skill_get with id "${selected.id}" in a run to load the full methodology`)
                            }}
                          >
                            Use in Run →
                          </Button>
                          <Button
                            variant='outline'
                            size='sm'
                            className='font-mono text-xs'
                            onClick={() => viewSkillBody(selected.id)}
                            disabled={loadingBody}
                          >
                            {loadingBody ? 'Loading...' : 'View Content →'}
                          </Button>
                        </div>
                        {skillBody && (
                          <div className='mt-3 max-h-48 overflow-auto rounded border border-border/50 bg-black/40 p-3'>
                            <pre className='whitespace-pre-wrap font-mono text-[11px] text-foreground/70'>
                              {skillBody}
                            </pre>
                          </div>
                        )}
                      </CardContent>
                    </Card>
                  )}
                </div>
              </div>
            </TabsContent>
          ))}
        </Tabs>
      </div>
    </div>
  )
}
