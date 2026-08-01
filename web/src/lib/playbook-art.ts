/**
 * Generated engagement art for builtin playbooks.
 * Files live under public/playbooks/{id}.webp
 */

const BY_ID: Record<string, string> = {
  'full-validation': '/playbooks/full-validation.webp',
  'web-owasp': '/playbooks/web-owasp.webp',
  'internal-network': '/playbooks/internal-network.webp',
  'recon-only': '/playbooks/recon-only.webp',
  'post-proof': '/playbooks/post-proof.webp',
  'cve-lab': '/playbooks/cve-lab.webp'
}

const FALLBACK = '/playbooks/full-validation.webp'

export function playbookArtSrc(id?: string | null): string {
  if (!id) return FALLBACK
  return BY_ID[id] ?? FALLBACK
}

export const PLAYBOOK_ART_IDS = Object.keys(BY_ID)
