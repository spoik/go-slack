# --- Base Stage ---
FROM node:25-alpine AS base
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# --- Development Stage ---
FROM base AS dev
CMD ["npm", "run", "dev"]

# --- Test Stage ---
FROM base AS test
# Tail dev null keeps the container running but not doing anything. This keeps it available to run tests on
# without using much CPU
CMD ["tail", "-f", "/dev/null"]
