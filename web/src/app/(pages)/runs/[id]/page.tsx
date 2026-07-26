// Component Imports
import RunDetail from '@/views/runs/RunDetail'

const RunDetailPage = async ({ params }: { params: Promise<{ id: string }> }) => {
  const { id } = await params

  return <RunDetail runId={id} />
}

export default RunDetailPage
