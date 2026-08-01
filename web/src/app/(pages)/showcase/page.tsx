import { redirect } from 'next/navigation'

/** Showcase tab removed — assets live on Overview / Agents / Skills / Playbooks / Login. */
export default function ShowcasePage() {
  redirect('/overview')
}
