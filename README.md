# Golang URL Shortener

A production-ready, cloud-native URL shortener built with **Hexagonal Architecture** and **AWS Serverless Stack**. This project demonstrates enterprise-grade development practices, DevOps automation, and scalable system design.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![AWS](https://img.shields.io/badge/AWS-Serverless-FF9900?style=flat&logo=amazon-aws)](https://aws.amazon.com)
[![Architecture](https://img.shields.io/badge/Architecture-Hexagonal-blue?style=flat)](https://alistair.cockburn.us/hexagonal-architecture/)

---

## Table of Contents

- [System Requirements](#system-requirements)
- [Key Features](#key-features)
- [Tech Stack](#tech-stack)
- [System Architecture](#system-architecture)
- [Quick Start](#quick-start)
- [Deployment](#deployment)
- [Testing](#testing)
- [Cost Analysis](#cost-analysis)
- [Architecture Breakdown](#architecture-breakdown)
- [Improvements](#improvements)
- [License](#license)

---

## System Requirements

### Functional Requirements

- **Shortening URLs** - Generate unique shortened URLs from long URLs
- **URL Redirection** - Fast redirection from short URL to original URL
- **Analytics** - Track click counts and usage statistics per shortened URL
- **API Access** - RESTful API endpoints for CRUD operations
- **User Notifications** - Slack integration for URL creation/deletion events

### Non-Functional Requirements

- **Scalability** - Auto-scaling to handle variable traffic loads
- **Performance** - High-speed URL redirection and generation with caching
- **Reliability** - High availability with fault tolerance
- **Security** - Protected against unauthorized access and abuse with input validation
- **Maintainability** - Clean code following Hexagonal Architecture principles
- **Monitoring** - Complete logging and performance tracking with CloudWatch

---

## Key Features

- **URL Generation** - Create shortened URLs efficiently with collision detection
- **Redirection** - High-performance redirect to original URLs with Redis caching
- **Statistics** - Real-time usage statistics and analytics with platform detection
- **Notifications** - Event-driven notifications via Slack integration
- **Deletion** - Safe removal of URLs with automatic cache invalidation
- **Caching** - Multi-layer caching strategy with ElastiCache (Redis) support
- **Security** - Input validation, malicious URL detection, and least-privilege IAM roles  

---

## Tech Stack

### Backend & Cloud Services

| Technology | Purpose | Why? |
|------------|---------|------|
| **Go (Golang)** | Primary Language | High performance, excellent concurrency |
| **AWS Lambda** | Serverless Compute | Auto-scaling, pay-per-use |
| **AWS DynamoDB** | NoSQL Database | Single-digit ms latency, fully managed |
| **ElastiCache (Redis)** | Caching Layer | Sub-millisecond response times |
| **AWS CloudFront** | CDN | Global content delivery, low latency |
| **AWS API Gateway** | API Management | REST API with built-in security |
| **AWS SQS** | Message Queue | Decoupled, scalable messaging |

### DevOps & Infrastructure

| Technology | Purpose |
|------------|---------|
| **AWS CloudFormation** | Infrastructure as Code |
| **AWS SAM CLI** | Local testing & deployment |
| **GitHub Actions** | CI/CD automation |

---

## System Architecture

### High-Level Design

![system design](./assets/system_design.png)

**Architecture Flow:**

1. **API Gateway** - Entry point with authentication and rate limiting
2. **Lambda Functions** - Serverless compute for each operation (generate, redirect, stats, delete, notify)
3. **ElastiCache (Redis)** - Fast caching layer for frequently accessed URLs
4. **DynamoDB** - Persistent storage with automatic scaling
5. **CloudFront CDN** - Global edge caching for minimal latency
6. **SQS** - Asynchronous processing for analytics and notifications

### Class Diagram

![class diagram](./assets/class_diagram.png)

---

## Quick Start

### Prerequisites

- Go 1.21 or higher installed
- AWS CLI configured with credentials
- AWS SAM CLI installed
- Basic knowledge of AWS services

### Installation

```bash
# Clone the repository
git clone https://github.com/dheeraj-vp/GoLang-URL-Shortener.git
cd golang-url-shortener

# Build the project
make build

# Run tests
make unit-test

# Deploy to AWS
make deploy
```

---

## Deployment

### Deploying to AWS Lambda

Deploy your serverless functions with a single command:

```bash
make deploy
```

This command uses **AWS SAM** to:
- Package Lambda functions
- Create/update CloudFormation stack
- Deploy all AWS resources
- Configure API Gateway endpoints

#### Deployment Options

**Without ElastiCache (Development):**
```bash
make build
sam deploy --parameter-overrides EnableElastiCache=false
```

**With ElastiCache (Production):**
```bash
make build
sam deploy --parameter-overrides \
  EnableElastiCache=true \
  ElastiCacheNodeType=cache.t3.micro \
  Environment=prod
```

### CI/CD Pipeline

The project includes **GitHub Actions** workflow (`.github/workflows/deploy.yml`) for automated:
- Code compilation and testing
- Linting with golangci-lint
- Security scanning with Trivy
- AWS deployment on push to main branch
- Automated rollback on failures

Configure the following GitHub secrets:
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `SLACK_TOKEN` (optional)
- `SLACK_CHANNEL_ID` (optional)

---

## Testing

### Run Unit Tests

```bash
make unit-test
```

Tests all core business logic in isolation, following Hexagonal Architecture principles. Includes comprehensive cache testing with mock implementations.

### Run Benchmark Tests

```bash
make benchmark-test
```

Performance benchmarks to ensure optimal response times.

### Clean Build Artifacts

```bash
make clean
```

Removes all compiled Lambda function binaries.

### Delete Deployed Stack

```bash
make delete
```

Removes all AWS resources created by CloudFormation stack.

---

## Cost Analysis

### Estimated Cost for 1 Million Requests

| AWS Service | Cost | Notes |
|-------------|------|-------|
| **AWS Lambda** | $0.20 | First 1M requests/month FREE, then $0.20 per 1M |
| **API Gateway** | $3.50 | First 1M requests/month FREE, then $3.50 per 1M |
| **DynamoDB** | $1.25 | On-demand pricing: 2 writes per generate, 1 per redirect, 2 reads per stats |
| **CloudFront** | $0.75 - $2.20 | Based on 1M HTTPS requests |
| **ElastiCache (Redis)** | $25-30/month | Depends on instance type (cache.t3.micro) |
| **SQS** | $0.40 | First 1M requests/month FREE, then $0.40 per 1M |

**Cost Optimization Strategies:**
- Redis caching reduces DynamoDB reads by approximately 85%
- CloudFront CDN minimizes Lambda cold starts
- Pay-per-use serverless model ensures no idle costs
- Cache-aside pattern significantly reduces database load

---

## Architecture Breakdown

### Hexagonal Architecture

![Hexagonal Architecture](./assets/hexagonal.png)

**What is Hexagonal Architecture?**

Also known as **Ports and Adapters Pattern**, it separates core business logic from external systems, making the application:
- **Testable** - Core logic tested independently with mock implementations
- **Maintainable** - Clear separation of concerns between layers
- **Flexible** - Easy to swap databases, APIs, or frameworks without changing business logic

### Project Structure

```
golang-url-shortener/
│
├── internal/
│   ├── adapters/              # Infrastructure Layer (Adapters)
│   │   ├── cache/            # Redis cache implementation
│   │   ├── repository/       # DynamoDB data access
│   │   ├── handlers/         # HTTP request handlers
│   │   └── functions/        # Lambda function entry points
│   │       ├── delete/       # Delete URL function
│   │       ├── generate/     # Generate short URL
│   │       ├── notification/ # Send notifications
│   │       ├── redirect/     # Redirect to original URL
│   │       └── stats/        # Get URL statistics
│   │
│   ├── core/                  # Domain Layer (Business Logic)
│   │   ├── domain/           # Domain models (link.go, stats.go)
│   │   ├── ports/            # Interface definitions (contracts)
│   │   └── services/         # Business logic implementation
│   │
│   ├── config/                # Configuration and constants
│   │   ├── config.go         # Environment configuration
│   │   └── constants.go      # Application constants
│   │
│   └── tests/
│       ├── benchmark/         # Performance tests
│       ├── unit/             # Unit tests
│       └── mock/             # Mock implementations
│
├── .github/
│   └── workflows/
│       └── deploy.yml        # CI/CD pipeline
│
├── assets/                    # Architecture diagrams
│   ├── system_design.png
│   ├── class_diagram.png
│   └── hexagonal.png
│
├── template.yaml              # AWS SAM/CloudFormation template
├── samconfig.toml            # SAM deployment configuration
├── Makefile                  # Build automation
├── .env.example              # Environment variables template
├── IMPROVEMENTS.md           # Detailed improvements documentation
├── MIGRATION_GUIDE.md        # Upgrade instructions
└── README.md                 # This file
```

### Why Hexagonal in Serverless?

**Benefits:**
1. **Decoupling** - Core logic independent of AWS services
2. **Testability** - Mock external dependencies easily with interface-based design
3. **Flexibility** - Replace Redis with Memcached without touching business logic
4. **Scalability** - Each Lambda function scales independently
5. **Cost-Effective** - Pay only for compute time used

**Example Flow:**

```
User Request
    ↓
API Gateway (Adapter)
    ↓
Lambda Handler (Adapter)
    ↓
Service Layer (Core/Port)
    ↓
Repository (Adapter)
    ↓
DynamoDB/Redis (External)
```

---

## Key Highlights

### DevOps Practices

- **Infrastructure as Code** - Complete CloudFormation templates with conditional resources
- **CI/CD Pipeline** - Automated GitHub Actions workflow with testing and security scanning
- **Serverless Architecture** - Auto-scaling, zero server management
- **Monitoring** - CloudWatch logs, metrics, and alarms for error rates and throttling
- **Cost Optimization** - Multi-layer caching strategy with ElastiCache support
- **Security** - Least-privilege IAM roles, input validation, and secure parameter handling

### Software Engineering Practices

- **Clean Architecture** - Hexagonal/Ports & Adapters pattern with clear layer separation
- **Go Best Practices** - Idiomatic Go code with proper error handling and context management
- **Microservices** - Event-driven, loosely coupled design with SQS integration
- **Testing** - Comprehensive unit tests with mock implementations and benchmarks
- **API Design** - RESTful endpoints with proper HTTP status codes and error handling
- **Code Quality** - Constants for magic numbers, structured logging, and comprehensive documentation  

---

## Improvements

This project has been enhanced with enterprise-grade improvements. See [IMPROVEMENTS.md](./IMPROVEMENTS.md) for detailed documentation.

**Key Enhancements:**
- Cache-aside pattern implementation with Redis/ElastiCache
- Collision detection for short URL generation
- Input validation and malicious URL detection
- Least-privilege IAM roles per Lambda function
- Context timeouts and proper error handling
- Platform detection from request headers
- Pagination support for DynamoDB operations
- Comprehensive test coverage with mock implementations
- CloudWatch alarms for monitoring
- VPC configuration for ElastiCache integration

For migration instructions, see [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md).

## Learning Outcomes

This project demonstrates:

**Cloud & DevOps:**
- AWS serverless stack (Lambda, API Gateway, DynamoDB, CloudFront, ElastiCache)
- Infrastructure as Code with CloudFormation and SAM
- CI/CD automation with GitHub Actions
- Cost-effective architecture design with caching strategies
- VPC configuration and security groups
- CloudWatch monitoring and alerting

**Software Engineering:**
- Hexagonal Architecture implementation
- Go programming best practices and idiomatic code
- RESTful API development with proper status codes
- Distributed systems concepts (caching, message queues)
- Test-driven development with comprehensive test suites
- Error handling and graceful degradation
- Context management and timeout handling

## License

This project is licensed under the MIT License - see the LICENSE file for details.
