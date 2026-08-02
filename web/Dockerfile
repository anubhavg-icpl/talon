# syntax=docker/dockerfile:1

FROM node:22-alpine AS base
# Pin pnpm — "latest" (11.x) changed build-script approval behavior and
# fails on this lockfile; 10.33.2 is the version the lockfile was built with.
# Prefer npm install over corepack prepare — corepack often times out offline
# / behind flaky GitHub registry access on lab networks.
RUN npm install -g pnpm@10.33.2

# ---- deps ----
FROM base AS deps
WORKDIR /app
COPY package.json pnpm-lock.yaml .npmrc pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

# ---- build ----
FROM base AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN pnpm build

# ---- runner ----
FROM node:22-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production \
    TALON_CORE_URL=http://localhost:8000 \
    PORT=3000 \
    HOSTNAME=0.0.0.0

COPY --from=build /app/.next/standalone ./
COPY --from=build /app/.next/static ./.next/static
COPY --from=build /app/public ./public

EXPOSE 3000
CMD ["node", "server.js"]
