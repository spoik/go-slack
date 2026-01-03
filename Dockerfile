# --- Base Stage ---
FROM golang:1.25-alpine AS base
WORKDIR /app
# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

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
