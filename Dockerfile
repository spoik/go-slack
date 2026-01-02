FROM golang:1.25-alpine

WORKDIR /app

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Air will handle the building and running of the binary
CMD ["air"]
