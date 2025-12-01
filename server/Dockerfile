# Build stage
# FROM golang:1.21-alpine AS builder
FROM golang:1.24.10-alpine AS builder
# docker pull golang:1.24.10-nanoserver

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN chmod +x kubectl
RUN mv kubectl /usr/local/bin/

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]