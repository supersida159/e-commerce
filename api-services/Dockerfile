# FROM alpine

# RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
# RUN update-ca-certificates 

# WORKDIR /app/
# ADD ./app /app/


# ENTRYPOINT [ "./app" ]

# # CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app
# # docker build -t e-commerce-test .
 
# # docker run -d --name e-commerce-test-container -e DATABASE_URI="root:Tungpro123@@tcp(e-commerce2:3306)/e-commerce?parseTime=true" -e REDIS_URI="my-redis:6379" --network=fd-network e-commerce-test

# buid all services in a docker for loadtest

# This Dockerfile should be placed in each service directory
FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY config.yaml /app/config.yaml

# The actual command will vary per service
CMD ["./main"]