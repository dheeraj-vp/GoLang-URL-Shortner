# 🚀 Golang URL Shortener

> A production-ready, cloud-native URL shortener built with **Hexagonal Architecture** and **AWS Serverless Stack**. Demonstrating enterprise-grade development practices, DevOps automation, and scalable system design.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![AWS](https://img.shields.io/badge/AWS-Serverless-FF9900?style=flat&logo=amazon-aws)](https://aws.amazon.com)
[![Architecture](https://img.shields.io/badge/Architecture-Hexagonal-blue?style=flat)](https://alistair.cockburn.us/hexagonal-architecture/)

---

## 📋 Table of Contents

- [System Requirements](#-system-requirements)
- [Key Features](#-key-features)
- [Tech Stack](#-tech-stack)
- [System Architecture](#-system-architecture)
- [Quick Start](#-quick-start)
- [Deployment](#-deployment)
- [Testing](#-testing)
- [Cost Analysis](#-cost-analysis)
- [Architecture Breakdown](#-architecture-breakdown)
- [License](#-license)

---

## 📊 System Requirements

### ✅ Functional Requirements

- **Shortening URLs** - Generate unique shortened URLs from long URLs
- **URL Redirection** - Fast redirection from short URL to original URL
- **Analytics** - Track click counts and usage statistics per shortened URL
- **API Access** - RESTful API endpoints for CRUD operations
- **User Notifications** - Slack integration for URL creation/deletion events

### ⚡ Non-Functional Requirements

- **Scalability** - Auto-scaling to handle variable traffic loads
- **Performance** - High-speed URL redirection and generation
- **Reliability** - High availability with fault tolerance
- **Security** - Protected against unauthorized access and abuse
- **Maintainability** - Clean code following Hexagonal Architecture principles
- **Monitoring** - Complete logging and performance tracking

---

## ✨ Key Features

🔗 **URL Generation** - Create shortened URLs efficiently  
🎯 **Redirection** - Lightning-fast redirect to original URLs  
📈 **Stats** - Real-time usage statistics and analytics  
🔔 **Notification** - Event-driven notifications via Slack  
🗑️ **Deletion** - Safe removal of URLs with data cleanup  

---

## 💻 Tech Stack

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

## 🏗 System Architecture

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

## 🚀 Quick Start

### Prerequisites

```bash
✅ Go 1.21 or higher installed
✅ AWS CLI configured with credentials
✅ AWS SAM CLI installed
✅ Basic knowledge of AWS services
```

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

## 🌐 Deployment

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

### CI/CD Pipeline

The project includes **GitHub Actions** workflow for automated:
- Code compilation and testing
- AWS deployment on push to main branch
- Automated rollback on failures

---

## 🧪 Testing

### Run Unit Tests

```bash
make unit-test
```

Tests all core business logic in isolation, following Hexagonal Architecture principles.

### Run Benchmark Tests

```bash
make benchmark-test
```

Performance benchmarks to ensure optimal response times.

### Clean Build Artifacts

```bash
make clean
```

### Delete Deployed Stack

```bash
make delete
```

Removes all AWS resources created by CloudFormation.

---

## 💰 Cost Analysis

### Estimated Cost for 1 Million Requests

| AWS Service | Cost | Notes |
|-------------|------|-------|
| **AWS Lambda** | $0.20 | First 1M requests/month FREE, then $0.20 per 1M |
| **API Gateway** | $3.50 | First 1M requests/month FREE, then $3.50 per 1M |
| **DynamoDB** | $1.25 | On-demand pricing: 2 writes per generate, 1 per redirect, 2 reads per stats |
| **CloudFront** | $0.75 - $2.20 | Based on 1M HTTPS requests |
| **ElastiCache (Redis)** | Variable | Depends on instance type and runtime hours |
| **SQS** | $0.40 | First 1M requests/month FREE, then $0.40 per 1M |

**💡 Cost Optimization:**
- Redis caching reduces DynamoDB reads by ~85%
- CloudFront CDN minimizes Lambda cold starts
- Pay-per-use serverless model ensures no idle costs

---

## 🏛 Architecture Breakdown

### Hexagonal Architecture

![Hexagonal Architecture](./assets/hexagonal.png)

**What is Hexagonal Architecture?**

Also known as **Ports and Adapters Pattern**, it separates core business logic from external systems, making the application:
- ✅ **Testable** - Core logic tested independently
- ✅ **Maintainable** - Clear separation of concerns
- ✅ **Flexible** - Easy to swap databases, APIs, or frameworks

### Project Structure

```
golang-url-shortener/
│
├── internal/
│   ├── adapters/              # 🔌 Infrastructure Layer (Adapters)
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
│   └── core/                  # 💎 Domain Layer (Business Logic)
│       ├── domain/           # Domain models (link.go, stats.go)
│       ├── ports/            # Interface definitions (contracts)
│       └── services/         # Business logic implementation
│
├── tests/
│   ├── benchmark/            # Performance tests
│   └── unit/                 # Unit tests
│
├── template.yaml             # AWS SAM/CloudFormation template
├── samconfig.toml           # SAM deployment configuration
├── Makefile                 # Build automation
└── README.md                # Documentation
```

### Why Hexagonal in Serverless?

**Benefits:**
1. **🔄 Decoupling** - Core logic independent of AWS services
2. **🧪 Testability** - Mock external dependencies easily
3. **🔧 Flexibility** - Replace Redis with Memcached without touching business logic
4. **📈 Scalability** - Each Lambda function scales independently
5. **💰 Cost-Effective** - Pay only for compute time used

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

## 🎯 Key Highlights

### DevOps Practices 

✅ **Infrastructure as Code** - Complete CloudFormation templates  
✅ **CI/CD Pipeline** - GitHub Actions workflow  
✅ **Serverless Architecture** - Auto-scaling, zero server management  
✅ **Monitoring** - CloudWatch logs and metrics integration  
✅ **Cost Optimization** - Multi-layer caching strategy  

### SDE Practices 

✅ **Clean Architecture** - Hexagonal/Ports & Adapters pattern  
✅ **Go Best Practices** - Idiomatic Go code  
✅ **Microservices** - Event-driven, loosely coupled design  
✅ **Testing** - Unit tests and benchmarks  
✅ **API Design** - RESTful endpoints with proper error handling  

---

## 📚 Learning Outcomes

This project demonstrates:

**☁️ Cloud & DevOps:**
- AWS serverless stack (Lambda, API Gateway, DynamoDB, CloudFront)
- Infrastructure as Code with CloudFormation
- CI/CD automation with GitHub Actions
- Cost-effective architecture design

**💻 Software Engineering:**
- Hexagonal Architecture implementation
- Go programming best practices
- RESTful API development
- Distributed systems concepts
- Test-driven development
