# Stage 1: Build the application (installs all dependencies for building)
FROM node:24-alpine AS builder
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile
COPY . .
RUN yarn build

# Stage 2: Prepare production dependencies (installs only production dependencies)
FROM node:24-alpine AS production_deps
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install --production --frozen-lockfile --ignore-scripts

# Stage 3: Final production image
FROM node:24-alpine
WORKDIR /app
ENV NODE_ENV=production

# Copy build artifacts from builder stage
COPY --from=builder /app/build ./build

# Copy production node_modules from production_deps stage
COPY --from=production_deps /app/node_modules ./node_modules

# Copy package.json and yarn.lock (optional, but good for completeness/debugging)
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/yarn.lock ./yarn.lock

EXPOSE 3000
CMD ["node", "build/index.js"]