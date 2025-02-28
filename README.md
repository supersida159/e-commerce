# E-commerce Microservices Platform

[![Go Version](https://img.shields.io/badge/Go-1.22.7+-00ADD8?style=flat&logo=go)](https://golang.org/doc/go1.22)
[![Gin Framework](https://img.shields.io/badge/Gin-Framework-00ADD8?style=flat&logo=go)](https://github.com/gin-gonic/gin)
[![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-Streaming-231F20?style=flat&logo=apache-kafka)](https://kafka.apache.org/)
[![Gorm](https://img.shields.io/badge/Gorm-v1.25.6-red)](https://gorm.io/)
[![Redis](https://img.shields.io/badge/Redis-v9.7.0-DC382D?style=flat&logo=redis)](https://redis.io/)
[![AWS SDK](https://img.shields.io/badge/AWS_SDK-v1.50.7-FF9900?style=flat&logo=amazon-aws)](https://aws.amazon.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A scalable and resilient e-commerce platform built with Go (Gin framework) using a microservices architecture and event-driven communication with Apache Kafka.

> **Note**: This project is made for practice and learning purposes. There are many simplifications and shortcomings compared to a production-ready implementation that would need improvement in a real-world scenario. See the [Areas for Improvement](#15-areas-for-improvement) section for details.

## 1. Architecture Overview

This project implements a microservices architecture with the following components:

- **API Central Service**: Acts as the gateway and orchestrator for all client requests
- **Order Service**: Handles order creation and management
- **Cart Service**: Manages user shopping carts
- **Inventory Service**: Tracks product inventory and availability
- **Authentication Service**: Handles user authentication, including Google OAuth
- **Image Upload Service**: Manages product image uploads to AWS S3

The system uses the **Saga Pattern** for distributed transactions, ensuring data consistency across services. Apache Kafka serves as the messaging backbone for event-driven communication between services.

### 1.1 Saga Pattern and Kafka Workflow

The system implements a choreography-based Saga pattern with Kafka as the event backbone:

![Architecture Diagram](https://github.com/supersida159/e-commerce/blob/Develop/workflow/e-commerce%20saga%20pattern-2024-10-25-023937.png)

### 1.2 Order Creation Flow:

1. Central System sends Create Order Request to the Orchestrator
2. Orchestrator initiates a transaction by:
   - 2a. Sending Create Cart Command to CartService.topic
   - 2b. Sending Update Cart Command to UpdateCart.topic
   - 2c. Sending Update Inventory Command to Inventory.topic
3. Each service updates its state and reports status to Saga.topic
   - 3a, 3b, 3c. Services send Update Status messages
4. Orchestrator monitors status via Saga.topic
5. If failure occurs, Orchestrator sends Rollback Commands to Rollback.topic
   - 6a, 6b, 6c. Services receive Rollback commands
   - 7a, 7b, 7c. Services update rollback status to Saga.topic
8. Central System receives Final Status from Orchestrator

This pattern ensures that all services either complete their transactions successfully or compensating transactions restore the system to a consistent state.

## 2. Key Features

- **Microservices Architecture**: Modular and independently deployable services
- **Event-Driven Design**: Apache Kafka for reliable message passing
- **Saga Pattern**: Distributed transaction management with compensation logic
- **Real-time Updates**: WebSocket support for live order status updates
- **HTTP Polling**: Alternative to WebSockets for order status updates
- **Authentication**: JWT-based authentication with Google OAuth support
- **Monitoring**: Grafana dashboards for system metrics
- **Load Testing**: K6 integration for performance testing
- **Containerization**: Docker for consistent development and deployment
- **Structured Logging**: Using Logrus for enhanced logging capabilities
- **Environment Management**: Flexible configuration via environment variables

## 3. Technology Stack

### 3.1 Core Technologies
- **Backend**: Go 1.22.7 (toolchain 1.23.6)
- **Web Framework**: Gin 1.9.1
- **Database**: MySQL 8.0 with GORM 1.25.6
- **Cache**: Redis 9.7.0
- **Message Broker**: Apache Kafka (IBM/sarama 1.43.3)
- **Service Discovery**: ZooKeeper

### 3.2 Authentication & Security
- **JWT**: dgrijalva/jwt-go 3.2.0
- **Google Auth**: cloud.google.com/go/auth 0.14.1
- **Validator**: go-playground/validator/v10 18.0

### 3.3 Real-time Communication
- **WebSockets**: gorilla/websocket 1.5.1
- **Socket.IO**: googollee/go-socket.io 1.7.0

### 3.4 Cloud Services
- **AWS SDK**: aws-sdk-go 1.50.7
- **S3 Storage**: For image uploads and static content

### 3.5 Development & Operations
- **Environment Management**: 
  - joho/godotenv 1.5.1
  - caarlos0/env/v6 6.10.1
- **Logging**: sirupsen/logrus 1.9.3
- **Monitoring**: Grafana, InfluxDB
- **Load Testing**: K6
- **Containerization**: Docker, Docker Compose

## 4. Services

### 4.1 API Central Service
Gateway service that routes requests to appropriate microservices and orchestrates Saga transactions.

### 4.2 Order Service
Handles the creation and management of customer orders.

### 4.3 Cart Service
Manages shopping cart operations and state.

### 4.4 Inventory Service
Tracks product inventory and handles reservation/release during order processing.

### 4.5 Authentication Service
Manages user authentication, including JWT issuance and Google OAuth integration.

### 4.6 Image Upload Service
Handles product image uploads to AWS S3.

## 5. Getting Started

### 5.1 Prerequisites

- Go 1.22.7+
- Docker and Docker Compose
- AWS account (for S3 image upload functionality)
- Google Developer account (for OAuth integration)

### 5.2 Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/ecommerce-microservices.git
   cd ecommerce-microservices
   ```

2. Configure environment variables:
   - Copy the sample configuration files for each service
   - Update with your specific settings
   ```bash
   cp ./api-services/pkg/config/config.sample.yaml ./api-services/pkg/config/config.yaml
   # Repeat for other services
   ```

3. Start the services using Docker Compose:
   ```bash
   docker-compose up -d
   ```

4. Access the application:
   - API Gateway: http://localhost:8090
   - Kafka UI: http://localhost:8080
   - Grafana: http://localhost:3000

### 5.3 Quick Start for Development

For local development without Docker:

```bash
# Set up environment variables
export GO111MODULE=on
export DATABASE_URI="root:password@tcp(localhost:3306)/e-commerce?parseTime=true"
export KAFKA_BROKERS="localhost:9092"

# Run the API service
cd api-services
go run main.go
```

## 6. Configuration

Each service requires proper configuration. A `config.sample.yaml` file with dummy values is provided for reference. **IMPORTANT: Never commit actual credentials to version control.**

### 6.1 Configuration Structure

```yaml
# General environment configuration
ENVIRONMENT: "Development"
HTTP_PORT: 8090
GRPC_PORT: 8889
AUTH_SECRET: "your-auth-secret"
SECRET_KEY: "your-secret-key"

# Database Configuration
DATABASE_URI: "root:password@tcp(mysql:3306)/e-commerce?parseTime=true"
DATABASE_TIMEOUT: 5 # Timeout in seconds for database connections

# Redis Configuration
REDIS_URI: "redis:6379"
REDIS_PASSWORD: ""
REDIS_DB: 0

# Kafka Configuration
KAFKA:
  KAFKA_BROKERS: ["kafka:29092"]
  KAFKA_RETRY: 3
  KAFKA_CONSUMER_OFFSET_RESET: "earliest"
  KAFKA_PRODUCER_REQUIRED_ACKS: 1
  KAFKA_ENABLE_TLS: false
  KAFKA_VERSION: "2.8.0"
  KAFKA_TIMEOUT: 3000

# AWS S3 Configuration
AWS_S3:
  AWS_REGION: "your-region"
  AWS_BUCKET_NAME: "your-bucket-name"
  AWS_ACCESS_KEY_ID: "your-access-key-id"
  AWS_SECRET_ACCESS_KEY: "your-secret-access-key"
  AWS_ENDPOINT: "https://your-bucket.s3.your-region.amazonaws.com"

# OAuth2 Configuration
OAUTH:
  OAUTH_CLIENT_ID: "your-oauth-client-id"
  OAUTH_CLIENT_SECRET: "your-oauth-client-secret"

# Logging Configuration
LOG:
  LOG_LEVEL: "info"
  LOG_FORMAT: "json"
  LOG_OUTPUT: "stdout"
```

### 6.2 Security Best Practices:

1. Use environment-specific configuration files
2. Store sensitive information in environment variables
3. Use Docker secrets or a secrets manager in production
4. Never commit actual credentials to version control

## 7. Development

### 7.1 Branch Structure

The project follows a structured branching strategy:

## Branch Code and Description
| Branch Code | Branch Description |
|--------------|--------------------|
| **01**       | Authentication Services |
| 01.01    | Authen using JWT |
| 01.02    | Login/signin using Google |
| **02**       | **API Central Services** |
| 02.01   | Build API, CRUD DB |
| 02.02    | Create Kafka for services, create topics |
| 02.03    | Create Producer/Consumer to work with Kafka services |
| 02.04    | Create subscriber with local pub/sub and goroutines (async jobs) |
| 02.05    | Create HTTP polling and WebSocket for order updates after creation |
| **03**       | **Order Services** |
| 03.01    | Create handler and compensation logic |
| 03.02    | Create Producer/Consumer to work with Kafka services |
| 03.03    | Create subscriber with local pub/sub and goroutines (async jobs) |
| **04**       | **Update Cart Services** |
| 04.01    | Create handler and compensation logic |
| 04.02    | Create Producer/Consumer to work with Kafka services |
| 04.03    | Create subscriber with local pub/sub and goroutines (async jobs) |
| **05**       | **Update Inventory Services** |
| 05.01    | Create handler and compensation logic |
| 05.02    | Create Producer/Consumer to work with Kafka services |
| 05.03    | Create subscriber with local pub/sub and goroutines (async jobs) |
| **06**       | **Upload Image Services** |
| 06.01    | Configure AWS S3 |
| 06.02    | Create Image Upload functionality with S3 |
| **07**       | **DevOps Services** |
| 07.01    | Create Docker Compose configurations |
| 07.02    | Set up K6 for load testing |
| 07.03    | Set up Grafana for monitoring |



### 7.2 Real-time Order Updates

The system provides two methods for clients to receive real-time order updates:

1. **WebSocket**:
   - Client establishes WebSocket connection with saga ID
   - Server pushes updates until order completion or failure
   - Implemented using gorilla/websocket and go-socket.io

2. **HTTP Polling**:
   - Client polls status endpoint with saga ID
   - Server responds with current order status
   - Optimized with Redis caching for performance

## 8. Kafka Topics

The system uses the following Kafka topics for communication:

1. **CartService.topic**: Handles create cart commands
2. **UpdateCart.topic**: Handles update cart commands
3. **Inventory.topic**: Handles inventory update commands
4. **Rollback.topic**: Handles compensation commands when transactions fail
5. **Saga.topic**: Central topic for monitoring transaction status across services

## 9. Infrastructure

### 9.1 Docker Compose Configuration

The project uses Docker Compose to orchestrate all required services. The configuration includes:

#### Infrastructure Services
- MySQL (Database)
- Redis (Caching)
- ZooKeeper (Kafka coordination)
- Kafka (Message broker)
- Kafka UI (Management interface)

#### Monitoring Stack
- InfluxDB (Time-series database)
- Grafana (Metrics visualization)
- K6 (Load testing)

#### Application Services
- API Services
- Create Order Service
- Update Cart Service
- Update Inventory Service

```yaml
# Docker Compose configuration excerpt
services:
  # Infrastructure Services
  mysql:
    image: mysql:8.0
    container_name: mysql
    ports:
      - "3308:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "your-password"
      MYSQL_DATABASE: "e-commerce"
    volumes:
      - mysql_data:/var/lib/mysql
    networks:
      - microservices-network
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
  
  # More services defined...
```

### 9.2 Load Testing with K6 and Grafana

The project includes a complete performance testing setup:

1. K6 for generating load
2. InfluxDB for storing test results
3. Grafana for visualizing performance metrics

Example K6 test results:
![LoadTest Sample Result Image](https://github.com/supersida159/e-commerce/blob/Develop/k6/loadtest%20sample%20result.png)

- HTTP Request Duration: Average response time and percentiles
- Virtual Users: Number of concurrent users during test
- Error Rate: Percentage of failed requests
- Requests Per Second: System throughput

## 10. API Documentation

### 10.1 Authentication
- `POST /api/auth/login`: Authenticate with username/password
- `POST /api/auth/google`: Authenticate with Google OAuth
- `POST /api/auth/refresh`: Refresh authentication token

### 10.2 Products
- `GET /api/products`: List all products
- `GET /api/products/:id`: Get product details
- `POST /api/products`: Create product (admin only)
- `PUT /api/products/:id`: Update product (admin only)
- `DELETE /api/products/:id`: Delete product (admin only)

### 10.3 Cart
- `GET /api/cart`: Get current user's cart
- `POST /api/cart/items`: Add item to cart
- `PUT /api/cart/items/:id`: Update cart item
- `DELETE /api/cart/items/:id`: Remove item from cart

### 10.4 Orders
- `POST /api/orders`: Create new order
- `GET /api/orders`: List user's orders
- `GET /api/orders/:id`: Get order details
- `GET /api/orders/:sagaId/status`: Poll order status by saga ID

### 10.5 Images
- `POST /api/images/upload`: Upload product image

### 10.6 WebSocket
- `WS /ws/orders/:sagaId`: Subscribe to real-time order updates

## 11. Security Considerations

### 11.1 Credentials Protection
- **NEVER** commit credentials to version control
- Use environment variables or secrets management
- Encrypt sensitive data in databases

### 11.2 API Security
- All endpoints are protected with JWT authentication
- OAuth integration for secure third-party authentication
- Rate limiting to prevent abuse

### 11.3 Data Protection
- Data validation at service boundaries
- SQL injection prevention
- XSS protection

## 12. Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Commit your changes: `git commit -m 'Add some feature'`
4. Push to the branch: `git push origin feature/your-feature-name`
5. Submit a pull request

## 13. License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 14. Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Apache Kafka](https://kafka.apache.org/)
- [Docker](https://www.docker.com/)
- [Grafana](https://grafana.com/)
- [IBM Sarama](https://github.com/IBM/sarama)
- [GORM](https://gorm.io/)
- [Redis](https://redis.io/)
- [AWS SDK](https://github.com/aws/aws-sdk-go)
- [Socket.IO](https://github.com/googollee/go-socket.io)
- [Logrus](https://github.com/sirupsen/logrus)

## 15. Areas for Improvement

As this project is developed for learning and practice purposes, there are several areas that would need improvement for a production-ready implementation:

### 15.1 Architecture Improvements
1. **Database Per Service**: Each microservice should have its own dedicated database to avoid bottleneck effects and ensure proper service isolation
2. **API Gateway**: Implement a dedicated API Gateway service (like Kong or Traefik) instead of relying on the central service for routing
3. **Service Mesh**: Add a service mesh (e.g., Istio, Linkerd) for advanced networking, security, and observability features
4. **Circuit Breaker**: Implement circuit breaker patterns to prevent cascading failures

### 15.2 Infrastructure Improvements
1. **Kubernetes Deployment**: Move from Docker Compose to Kubernetes for better scaling and management
2. **Multi-region Deployment**: Support deploying services across multiple regions for better availability
3. **Auto-scaling**: Implement auto-scaling based on traffic patterns
4. **Chaos Engineering**: Add chaos testing to ensure system resilience

### 15.3 Security Improvements
1. **Enhanced Authentication**: Add multi-factor authentication options
2. **Fine-grained Authorization**: Implement more detailed RBAC (Role-Based Access Control)
3. **Secrets Management**: Use a dedicated secrets management solution like HashiCorp Vault or AWS Secrets Manager
4. **Security Scanning**: Add automated security scanning for dependencies and containers

### 15.4 Operational Improvements
1. **Distributed Tracing**: Implement OpenTelemetry or Jaeger for request tracing across services
2. **Advanced Monitoring**: Add detailed service health checks and business metrics monitoring
3. **Centralized Logging**: Implement ELK stack or similar for centralized log management
4. **Disaster Recovery**: Develop comprehensive backup and recovery procedures

### 15.5 Performance Improvements
1. **CDN Integration**: Add CDN support for static content delivery
2. **Caching Strategy**: Implement a more sophisticated caching strategy across services
3. **Read Replicas**: Add database read replicas for read-heavy services
4. **Connection Pooling**: Optimize database and Redis connection pooling

These improvements would be necessary steps to transform this learning project into a production-ready e-commerce platform.
