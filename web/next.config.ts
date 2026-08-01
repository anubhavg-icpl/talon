import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  basePath: process.env.BASEPATH ?? '',
  reactStrictMode: true,
  output: 'standalone',
  poweredByHeader: false,
  compress: true,
  allowedDevOrigins: ['192.168.*.*'],
  pageExtensions: ['js', 'jsx', 'ts', 'tsx'],
  images: {
    formats: ['image/webp', 'image/avif'],
    dangerouslyAllowSVG: false
  },
  headers: async () => [
    {
      source: '/(.*)',
      headers: [
        { key: 'X-Content-Type-Options', value: 'nosniff' },
        { key: 'X-Frame-Options', value: 'SAMEORIGIN' },
        { key: 'Referrer-Policy', value: 'strict-origin-when-cross-origin' },
        { key: 'Permissions-Policy', value: 'camera=(), microphone=(), geolocation=()' }
      ]
    }
  ],
  redirects: async () => {
    return [
      {
        source: '/',
        destination: '/overview',
        permanent: true
      }
    ]
  }
}

export default nextConfig
