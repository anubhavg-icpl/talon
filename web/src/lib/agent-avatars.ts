/**
 * Cropped filmstrip portraits → agent mode avatars.
 * Source: public/showcase/talon-agent-filmstrip.webp
 */

const BY_ID: Record<string, string> = {
  full: '/agents/full.webp',
  recon: '/agents/recon.webp',
  exploit: '/agents/exploit.webp',
  web: '/agents/web.webp',
  network: '/agents/network.webp',
  post: '/agents/post.webp'
}

const BY_CODENAME: Record<string, string> = {
  COMMANDER: '/agents/full.webp',
  GHOST: '/agents/recon.webp',
  STRIKER: '/agents/exploit.webp',
  'STRIKER-WEB': '/agents/web.webp',
  PHANTOM: '/agents/network.webp',
  CIPHER: '/agents/post.webp'
}

export function agentAvatarSrc(idOrCodename?: string | null): string {
  if (!idOrCodename) return BY_ID.full
  const raw = idOrCodename.trim()
  const lower = raw.toLowerCase()
  if (BY_ID[lower]) return BY_ID[lower]
  const upper = raw.toUpperCase()
  if (BY_CODENAME[upper]) return BY_CODENAME[upper]
  // codename-ish without separators
  if (upper.includes('WEB')) return BY_ID.web
  if (upper.includes('GHOST') || lower === 'explore') return BY_ID.recon
  if (upper.includes('PHANTOM') || lower === 'internal') return BY_ID.network
  if (upper.includes('CIPHER') || lower.includes('post')) return BY_ID.post
  if (upper.includes('STRIKER') || lower === 'exploit') return BY_ID.exploit
  return BY_ID.full
}

export const AGENT_AVATAR_IDS = Object.keys(BY_ID)
