// Next Imports
import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

/**
 * Session gate: everything except the login page, the login API call itself,
 * and static assets requires the talon_session cookie (set by talon-core
 * through the /api/talon proxy). Without it → /login.
 */
export function proxy(req: NextRequest) {
  const token = req.cookies.get('talon_session')?.value

  if (!token) {
    const url = req.nextUrl.clone()

    url.pathname = '/login'
    url.search = ''

    return NextResponse.redirect(url)
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/((?!login|api/talon/auth/login|_next/static|_next/image|favicon.webp|.*\\.webp|.*\\.svg).*)']
}
