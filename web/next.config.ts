import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  basePath: process.env.BASEPATH ?? '',
  reactStrictMode: true,
  output: 'standalone',
  allowedDevOrigins: ['192.168.*.*'],
  pageExtensions: ['js', 'jsx', 'ts', 'tsx'],
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
