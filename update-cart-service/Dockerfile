# update-cart-service/Dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

# Install git and SSL ca certificates
RUN apk update && \
    apk add --no-cache \
    git \
    ca-certificates \
    gcc \
    musl-dev \
    && update-ca-certificates

# Set necessary Git configs for private repos if needed
RUN git config --global http.sslVerify false

# Create and set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Set GOPROXY for better download reliability
ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o main .

# Final stage
FROM alpine:latest

# Add necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -g '' appuser

# Create and set working directory
WORKDIR /app

# Copy binary and config from builder
COPY --from=builder /app/main .
COPY --from=builder /app/pkg/config/config.sample.yaml /app/pkg/config/config.sample.yaml

# Set ownership of the application files
RUN chown -R appuser:appuser /app

# Use non-root user
USER appuser
EXPOSE 8090
CMD ["./main"]