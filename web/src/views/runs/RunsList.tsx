'use client'

// React Imports
import { useEffect, useMemo, useState } from 'react'

// Next Imports
import Link from 'next/link'
import { useRouter } from 'next/navigation'

// Third-party Imports
import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  useReactTable
} from '@tanstack/react-table'
import type { ColumnDef, PaginationState, SortingState } from '@tanstack/react-table'
import { ArrowUpDownIcon } from 'lucide-react'

// Type Imports
import type { RunStatus, RunSummary } from '@/lib/api'

// Component Imports
import PageHeader from '@/components/shared/PageHeader'
import StatusBadge from '@/components/shared/StatusBadge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

// Util Imports
import { listRuns } from '@/lib/api'
import { relativeTime } from '@/lib/format'

const STATUS_OPTIONS: { value: RunStatus | 'all'; label: string }[] = [
  { value: 'all', label: 'ALL STATUSES' },
  { value: 'running', label: 'RUNNING' },
  { value: 'awaiting_approval', label: 'AWAITING APPROVAL' },
  { value: 'completed', label: 'COMPLETED' },
  { value: 'error', label: 'ERROR' },
  { value: 'initializing', label: 'INITIALIZING' }
]

const Verdict = ({ run }: { run: RunSummary }) => {
  if (run.judge_verdict === true) {
    return <span className='text-primary font-mono text-xs tracking-widest'>✓ COMPROMISED</span>
  }

  if (run.status === 'completed') {
    return <span className='text-muted-foreground font-mono text-xs tracking-widest'>✗ CLEAN</span>
  }

  return <span className='text-muted-foreground font-mono text-xs'>—</span>
}

const columns: ColumnDef<RunSummary>[] = [
  {
    accessorKey: 'status',
    header: () => <span className='micro-label'>STATUS</span>,
    cell: ({ row }) => <StatusBadge status={row.original.status} />,
    filterFn: (row, _id, value) => value === 'all' || row.original.status === value
  },
  {
    accessorKey: 'target',
    header: ({ column }) => (
      <button
        className='micro-label hover:text-foreground flex items-center gap-1'
        onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
      >
        TARGET <ArrowUpDownIcon className='size-3' />
      </button>
    ),
    cell: ({ row }) => <span className='font-mono text-sm'>{row.original.target}</span>
  },
  {
    accessorKey: 'cve_id',
    header: () => <span className='micro-label'>CVE</span>,
    cell: ({ row }) => <span className='text-muted-foreground font-mono text-xs'>{row.original.cve_id ?? '—'}</span>
  },
  {
    accessorKey: 'service_name',
    header: () => <span className='micro-label'>SERVICE</span>,
    cell: ({ row }) => (
      <span className='text-muted-foreground font-mono text-xs'>{row.original.service_name ?? '—'}</span>
    )
  },
  {
    accessorKey: 'tool_calls',
    header: () => <span className='micro-label text-right'>TOOLS</span>,
    cell: ({ row }) => <span className='font-mono text-sm'>{row.original.tool_calls}</span>
  },
  {
    id: 'verdict',
    header: () => <span className='micro-label'>VERDICT</span>,
    cell: ({ row }) => <Verdict run={row.original} />
  },
  {
    accessorKey: 'started_at',
    header: ({ column }) => (
      <button
        className='micro-label hover:text-foreground ml-auto flex items-center gap-1'
        onClick={() => column.toggleSorting(column.getIsSorted() === 'asc')}
      >
        STARTED <ArrowUpDownIcon className='size-3' />
      </button>
    ),
    cell: ({ row }) => (
      <span className='text-muted-foreground font-mono text-xs'>{relativeTime(row.original.started_at)}</span>
    )
  }
]

const PAGE_SIZE = 20

const RunsList = () => {
  const router = useRouter()
  const [runs, setRuns] = useState<RunSummary[] | null>(null)
  const [total, setTotal] = useState(0)
  const [error, setError] = useState<string | null>(null)
  const [globalFilter, setGlobalFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState<RunStatus | 'all'>('all')
  const [sorting, setSorting] = useState<SortingState>([{ id: 'started_at', desc: true }])
  const [pagination, setPagination] = useState<PaginationState>({ pageIndex: 0, pageSize: PAGE_SIZE })

  useEffect(() => {
    let mounted = true

    const load = () =>
      listRuns(pagination.pageSize, pagination.pageIndex * pagination.pageSize)
        .then(res => {
          if (!mounted) return
          setRuns(res.runs ?? [])
          setTotal(res.total ?? 0)
          setError(null)
        })
        .catch(err => mounted && setError(err instanceof Error ? err.message : String(err)))

    load()
    const id = setInterval(load, 5000)

    return () => {
      mounted = false
      clearInterval(id)
    }
  }, [pagination.pageIndex, pagination.pageSize])

  const data = useMemo(() => runs ?? [], [runs])
  const pageCount = Math.max(1, Math.ceil(total / PAGE_SIZE))

  // Server-side pagination — sorting/filtering below applies to the current page only.
  // eslint-disable-next-line react-hooks/incompatible-library -- TanStack Table returns non-memoizable fns; safe for this table
  const table = useReactTable({
    data,
    columns,
    state: { globalFilter, sorting, pagination, columnFilters: [{ id: 'status', value: statusFilter }] },
    onGlobalFilterChange: setGlobalFilter,
    onSortingChange: setSorting,
    onPaginationChange: setPagination,
    globalFilterFn: 'includesString',
    manualPagination: true,
    pageCount,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel()
  })

  return (
    <div className='flex flex-col gap-6'>
      <PageHeader
        title='OPERATIONS'
        subtitle='ALL PENTEST RUNS'
        action={
          <Button
            className='font-mono text-xs font-semibold tracking-widest uppercase'
            render={<Link href='/runs/new' />}
            nativeButton={false}
          >
            [ + NEW OPERATION ]
          </Button>
        }
      />

      {error && (
        <div className='border-destructive/40 bg-destructive/10 text-destructive rounded-md border px-4 py-3 font-mono text-xs tracking-widest uppercase'>
          CORE UNREACHABLE — {error}
        </div>
      )}

      <Card>
        <CardHeader className='flex flex-wrap items-center justify-between gap-3'>
          <CardTitle className='micro-label'>RUN REGISTRY</CardTitle>
          <div className='flex flex-wrap items-center gap-2'>
            <Input
              value={globalFilter}
              onChange={e => setGlobalFilter(e.target.value)}
              placeholder='filter page…'
              className='w-full font-mono text-xs sm:w-64'
            />
            <Select value={statusFilter} onValueChange={value => setStatusFilter(value as RunStatus | 'all')}>
              <SelectTrigger className='w-full font-mono text-xs tracking-widest uppercase sm:w-48'>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {STATUS_OPTIONS.map(opt => (
                  <SelectItem key={opt.value} value={opt.value} className='font-mono text-xs tracking-widest uppercase'>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </CardHeader>
        <CardContent>
          {runs === null ? (
            <div className='flex flex-col gap-2'>
              <Skeleton className='h-10 w-full' />
              <Skeleton className='h-10 w-full' />
              <Skeleton className='h-10 w-full' />
            </div>
          ) : total === 0 ? (
            <div className='flex h-48 flex-col items-center justify-center gap-3'>
              <p className='micro-label'>NO OPERATIONS YET</p>
              <Link href='/runs/new' className='text-primary font-mono text-xs tracking-widest uppercase hover:underline'>
                [ + LAUNCH YOUR FIRST RUN ]
              </Link>
            </div>
          ) : (
            <>
              <div className='overflow-x-auto'>
                <Table>
                <TableHeader>
                  {table.getHeaderGroups().map(headerGroup => (
                    <TableRow key={headerGroup.id}>
                      {headerGroup.headers.map(header => (
                        <TableHead key={header.id}>
                          {header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}
                        </TableHead>
                      ))}
                    </TableRow>
                  ))}
                </TableHeader>
                <TableBody>
                  {table.getRowModel().rows.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={columns.length} className='micro-label h-32 text-center'>
                        NO RUNS MATCH FILTER
                      </TableCell>
                    </TableRow>
                  ) : (
                    table.getRowModel().rows.map(row => (
                      <TableRow
                        key={row.id}
                        className='hover:bg-muted/50 cursor-pointer'
                        onClick={() => router.push(`/runs/${row.original.run_id}`)}
                      >
                        {row.getVisibleCells().map(cell => (
                          <TableCell key={cell.id}>{flexRender(cell.column.columnDef.cell, cell.getContext())}</TableCell>
                        ))}
                      </TableRow>
                    ))
                  )}
                </TableBody>
                </Table>
              </div>
              <div className='mt-4 flex items-center justify-between'>
                <p className='micro-label'>
                  PAGE {pagination.pageIndex + 1} / {pageCount} OF {total}
                </p>
                <div className='flex gap-2'>
                  <Button
                    variant='outline'
                    size='sm'
                    className='font-mono text-xs tracking-widest uppercase'
                    disabled={!table.getCanPreviousPage()}
                    onClick={() => table.previousPage()}
                  >
                    ← PREV
                  </Button>
                  <Button
                    variant='outline'
                    size='sm'
                    className='font-mono text-xs tracking-widest uppercase'
                    disabled={!table.getCanNextPage()}
                    onClick={() => table.nextPage()}
                  >
                    NEXT →
                  </Button>
                </div>
              </div>
            </>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

export default RunsList
