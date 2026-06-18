# EucastanPay

## Overview

EucastanPay is a distributed fintech payment platform built using Go, PostgreSQL, Kafka, gRPC, and REST APIs. The system follows a microservices architecture where each service owns its own domain, database, and business logic.

The project was designed to explore real-world backend engineering concepts including:

- Microservices
- Event-Driven Architecture
- CQRS
- Distributed Transactions
- Idempotency
- Audit Logging
- Double-Entry Ledger Accounting
- Retry and Dead Letter Queues (DLQ)
- Service-to-Service Communication using gRPC
- JWT Authentication and Authorization

The primary goal is to build a reliable and scalable payment processing platform while learning production-grade backend architecture.

---

# Architecture

The system uses a hybrid communication model:

### Synchronous Communication

Used when immediate responses are required.

- REST APIs
- gRPC

### Asynchronous Communication

Used for cross-service workflows and event propagation.

- Apache Kafka

---

# Services

## User Service

Responsible for user management and onboarding.

### Responsibilities

- User registration
- User profile management
- KYC verification
- Authentication support

### Events Published

- user.registered
- user.register.failed
- user.kyc.verified

---

## Account Service

Responsible for bank account management and balance updates.

### Responsibilities

- Account creation
- Account lookup
- Balance management
- Debit operations
- Credit operations

### Events Published

- account.created
- debit.completed
- debit.failed
- credit.completed
- credit.failed
- reverse.debit

### Architectural Decision

Account balances are maintained here for fast reads while the Ledger Service acts as the source of truth for financial history.

---

## Transfer Service

Coordinates money movement between accounts.

### Responsibilities

- Transfer initiation
- Transaction orchestration
- Failure handling
- Compensation workflows

### Events Published

- transfer.initiated
- transfer.completed
- transfer.failed

### Architectural Decision

The Transfer Service acts as the workflow coordinator. It does not directly own balances but coordinates Account and Ledger operations through events.

---

## Ledger Service

Implements double-entry bookkeeping.

### Responsibilities

- Create debit ledger entries
- Create credit ledger entries
- Maintain immutable transaction records
- Calculate balances from ledger entries
- Support reconciliation

### Events Published

- ledger.created
- ledger.reconciliation.alert

### Architectural Decision

Ledger entries are immutable and never updated after creation.

Every transfer generates:

1. Debit entry
2. Credit entry

This ensures a complete financial audit trail.

---

## Audit Service

Provides compliance and observability capabilities.

### Responsibilities

- Store raw events
- Maintain searchable audit records
- Support investigation and reporting
- Track transaction history

### Components

### audit_logs

Stores immutable event payloads exactly as received.

### audit_read

Stores optimized queryable audit records.

### Architectural Decision

The service follows CQRS principles:

- Write Model → audit_logs
- Read Model → audit_read

This allows efficient searching without modifying historical event data.

---

# Event-Driven Architecture

Kafka is used as the system event bus.

Services communicate by publishing and consuming events rather than directly calling each other whenever possible.

Benefits:

- Loose coupling
- Better scalability
- Independent deployments
- Improved fault tolerance

---

# Idempotency

Distributed systems may deliver the same message more than once.
To prevent duplicate processing, services maintain a processed_events table.
Each event is processed exactly once based on its unique identifier.

---

# Outbox Pattern

Outbox pattern is used to prevent data loss due to failed published.
So, we save data before publishing it with outbox worker.

---

# Retry and Dead Letter Queue

Consumers automatically retry failed messages.
This prevents message loss while avoiding infinite retry loops.

---

# Security

Authentication is implemented using JWT tokens.
Protected endpoints require valid access tokens.
Additional security features include:

- Request logging
- Recovery middleware
- Authentication middleware
- RBAC middleware
- gRPC interceptors

---

# Database Strategy

Each service owns its own database schema.

Benefits:

- Service independence
- Reduced coupling
- Better scalability
- Easier deployments

Services never access another service's database directly.

Communication occurs through:

- gRPC
- REST APIs
- Kafka Events

---

# Future Enhancements

Planned improvements include:

- Notification Service
- Reconciliation Engine
- Fraud Detection Service
- OpenTelemetry Tracing
- Prometheus Metrics
- Grafana Dashboards
- Kubernetes Deployment

---

# Learning Objectives

This project serves as a practical implementation of:

- Domain-Driven Design (DDD)
- Clean Architecture
- Event-Driven Systems
- CQRS
- Microservices
- Distributed Systems
- Fintech Backend Engineering

The focus is not only on feature development but also on understanding the trade-offs and operational challenges involved in building reliable financial systems.
