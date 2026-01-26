# --- Base Stage ---
FROM golang:1.25-alpine AS base
WORKDIR /app
# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Install migrate for database migrations
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4

# --- Development Stage ---
FROM base AS dev
# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest
COPY . .
# Air will handle the building and running of the binary via hotreloading
CMD ["air"]

# --- Test Stage ---
FROM base AS test
COPY . .
# Tail dev null keeps the container running but not doing anything. This keeps it available to run tests on
# without using much CPU
CMD ["tail", "-f", "/dev/null"]

