# Multi-stage build for Node.js API
FROM node:20-alpine AS base
RUN apk add --no-cache libc6-compat python3 make g++
WORKDIR /workspace

# Dependencies stage
FROM base AS deps
COPY package*.json ./
COPY packages/database/package*.json ./packages/database/
COPY packages/common/package*.json ./packages/common/
COPY apps/api-monitor/package*.json ./apps/api-monitor/
RUN npm ci

# Builder stage
FROM base AS builder
COPY --from=deps /workspace/node_modules ./node_modules
COPY . .
RUN npm run build --workspace=@sales-monitor/database
RUN npm run build --workspace=@sales-monitor/common
RUN npm run build --workspace=@sales-monitor/api-monitor

# Runner stage
FROM base AS runner
ENV NODE_ENV=production
COPY --from=builder /workspace/node_modules ./node_modules
COPY --from=builder /workspace/packages/database/dist ./packages/database/dist
COPY --from=builder /workspace/packages/database/package.json ./packages/database/
COPY --from=builder /workspace/packages/common/dist ./packages/common/dist
COPY --from=builder /workspace/packages/common/package.json ./packages/common/
COPY --from=builder /workspace/apps/api-monitor/dist ./apps/api-monitor/dist
COPY --from=builder /workspace/apps/api-monitor/package.json ./apps/api-monitor/

WORKDIR /workspace/apps/api-monitor

EXPOSE 3000

CMD ["node", "dist/apps/api-monitor/src/main.js"]
