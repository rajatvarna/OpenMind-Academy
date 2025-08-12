# --- 1. Base Image ---
# Use a specific Node.js version for consistency
FROM node:18-alpine AS base

# --- 2. Installer Stage ---
# Install dependencies
FROM base AS installer
WORKDIR /app
COPY package.json ./
RUN npm install

# --- 3. Builder Stage ---
# Build the Next.js application
FROM base AS builder
WORKDIR /app
COPY --from=installer /app/node_modules ./node_modules
COPY . .
RUN npm run build

# --- 4. Runner Stage ---
# Create the final, minimal image
FROM base AS runner
WORKDIR /app

# Set environment variables for production
ENV NODE_ENV=production

# Copy the standalone output from the builder stage.
# This includes the minimal server, .next/static, and public assets.
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# The Next.js standalone server runs on port 3000 by default
EXPOSE 3000

# The command to start the Next.js server
CMD ["node", "server.js"]
