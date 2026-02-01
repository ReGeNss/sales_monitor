# Multi-stage build for Node.js API
FROM node:20-alpine AS base
RUN apk add --no-cache libc6-compat python3 make g++
WORKDIR /workspace

# Dependencies stage
FROM base AS deps
COPY package*.json ./
COPY node/packages/database/package*.json ./node/packages/database/
COPY node/packages/common/package*.json ./node/packages/common/
COPY node/apps/api-monitor/package*.json ./node/apps/api-monitor/
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
COPY --from=builder /workspace/node/packages/database/dist ./node/packages/database/dist
COPY --from=builder /workspace/node/packages/database/package.json ./node/packages/database/
COPY --from=builder /workspace/node/packages/common/dist ./node/packages/common/dist
COPY --from=builder /workspace/node/packages/common/package.json ./node/packages/common/
COPY --from=builder /workspace/node/apps/api-monitor/dist ./node/apps/api-monitor/dist
COPY --from=builder /workspace/node/apps/api-monitor/package.json ./node/apps/api-monitor/

WORKDIR /workspace/node/apps/api-monitor

EXPOSE 3000

CMD ["node", "dist/apps/api-monitor/src/main.js"]
