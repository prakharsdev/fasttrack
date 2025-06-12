# FastTrack Data Ingestion System

## Project Overview

This project implements a complete data ingestion pipeline using **Golang**, **RabbitMQ**, **MySQL**, and **Docker Compose**. It simulates a real-world event streaming use case where payloads are published into a message broker, consumed in real-time, and stored into a relational database for further analytics.

The project is fully containerized to support reproducibility, CI/CD integration, and easy handover across teams. The structure strictly follows a modular approach following Golang standards.

---

## System Architecture Diagram

![Publisher Golang](https://github.com/user-attachments/assets/5c01ee4b-accb-492f-9e71-ba0ec8c33bf7)


---

## Tech Stack

| Layer            | Technology Used                                                                       |
| ---------------- | ------------------------------------------------------------------------------------- |
| Language         | Golang (1.24.x)                                                                       |
| Messaging Broker | RabbitMQ (3.x)                                                                        |
| Database         | MySQL 8.x                                                                             |
| Containerization | Docker + Docker Compose                                                               |
| Configuration    | .env + Config Loader                                                                  |
| Project Layout   | [golang-standards/project-layout](https://github.com/golang-standards/project-layout) |

---

## Project Structure
The repository follows a modular Go monorepo-style layout:
```
fasttrack/
├── build/
│   └── docker-compose.yml       # Docker Compose for full stack orchestration
│
├── cmd/app/
│   └── main.go                  # Application entrypoint
│
├── config/
│   └── config.go                # Environment loader (.env file reader)
│
├── internal/
│   ├── consumer/
│   │   └── consumer.go          # RabbitMQ consumer logic
│   │
│   ├── db/
│   │   └── mysql.go             # MySQL connection + table initialization
│   │
│   ├── logger/
│   │   └── logger.go            # Centralized logger wrapper
│   │
│   └── publisher/
│       └── publisher.go         # RabbitMQ publisher logic
│
├── test/
│   ├── consumer_test.go         # Unit tests for consumer logic
│   └── publisher_test.go        # Unit tests for publisher logic
│
├── .env                          # Environment variables (loaded via godotenv)
├── Dockerfile                    # Application container build
├── Makefile                      # Automation commands
├── go.mod / go.sum               # Go module dependencies
└── README.md                     # Documentation

```

---

## Setup Instructions

This section explains how to fully set up and run the project locally, including dependencies, environment preparation, and startup sequence.

### 1. Clone the repository

```bash
git clone https://github.com/prakharsdev/fasttrack.git
cd fasttrack
```

### 2. Prepare your environment

Ensure you have the following installed:

| Dependency | Recommended Version | Notes                                  |
| ---------- | ------------------- | -------------------------------------- |
| Docker     | v26.1.1             | Required for full stack startup        |
| Golang     | v1.24.4             | Required for local development         |
| Make       | v4.4.1              | Used for automation scripts            |
| RabbitMQ   | v3.13.x (optional)  | Only if running standalone tests       |
| Erlang     | OTP 26.x (optional) | Dependency for RabbitMQ native install |

Note: You don't need RabbitMQ or Erlang locally if you only run via Docker Compose (make up).
These are only needed for advanced local debugging without containers.

### 3. Setup environment variables

Create a `.env` file in the project root (if not already present):

```env
# MySQL configuration
MYSQL_USER=root
MYSQL_PASSWORD=root
MYSQL_DATABASE=fasttrack
MYSQL_HOST=mysql
MYSQL_PORT=3306

# RabbitMQ configuration
RABBITMQ_URI=amqp://guest:guest@rabbitmq:5672/
```

> Note: The service names (`mysql`, `rabbitmq`) match the container names declared in Docker Compose and are used for internal DNS resolution.

### 4. Start the system using Makefile

```bash
make up
```

This will:

* Build the application Docker image
* Spin up RabbitMQ and MySQL containers
* Start your app container with healthchecks and dependency ordering

You’ll see logs from all services in your terminal.

### 5. View application logs

```bash
make logs
```

This will tail live logs from all containers, including publisher and consumer activity.

### 6. Run tests (optional)

```bash
make test
```

Executes unit tests for publisher and consumer logic to validate the ingestion pipeline.


---
## Problem Statement

The objective behind this implementation was to simulate a highly reliable, production-like ingestion system capable of:

* Ingesting transactional payment events via RabbitMQ
* Consuming, validating, and persisting records into MySQL
* Handling transient infrastructure failures gracefully via connection retries
* Ensuring strict idempotency at the database layer
* Maintaining full modular separation of concerns for future scalability
* Providing automation-driven local setup for reproducibility and ease of onboarding

---

## Design Decisions

### 1️ Monorepo Structure with Go Standards

I followed the `golang-standards/project-layout` convention because it aligns well with how production-grade Go services are structured across most organizations. This structure allows me to:

* Clearly separate responsibilities between configuration, internal services, domain logic, and entrypoints.
* Onboard any future engineer or reviewer very quickly as the folder structure is self-explanatory.
* Isolate and test individual components (publisher, consumer, DB, logger) independently without tight coupling.
* Allow the system to organically scale into additional modules (e.g., APIs, scheduled jobs, monitoring) without re-architecting.

My personal experience with Go has taught me that upfront modularity pays long-term dividends, especially when the system grows.

---

### 2️ Fully Dockerized Stack with Dependency Healthchecks

For this ingestion pipeline, I wanted to simulate a production-like environment locally. Docker Compose allows me to:

* Bundle MySQL, RabbitMQ, and my application as isolated services.
* Eliminate host-to-container network issues by leveraging Docker's internal service name resolution (e.g., `mysql`, `rabbitmq`).
* Ensure correct startup order using Docker Compose `healthcheck` and `depends_on: condition: service_healthy` to avoid race conditions at startup.
* Allow any engineer (or CI/CD system) to spin up the full stack on any machine by running a single `make up`.

This design reduces "works-on-my-machine" type issues during onboarding, testing, or handover.

---

### 3️ RabbitMQ as Messaging Backbone

I chose RabbitMQ because it's widely adopted for durable, reliable message brokering in many backend ecosystems. In this pipeline:

* I kept the messaging model simple using classic work queues.
* Messages are only acknowledged (`Ack`) after successful database insertion to ensure at-least-once delivery.
* Duplicate handling is fully addressed at the database level to ensure safe ingestion, as explained in the dedicated Duplicate Handling section.
* Publisher and consumer logic are separated to support independent scaling, making the design horizontally scalable for production workloads.

This architecture mimics the ingestion patterns I've seen work reliably under load in real-world systems.

---

### 4️ MySQL for Durable Storage

While message brokers ensure reliable delivery, they are not meant for long-term state. MySQL provides:

* Durable persistence for historical payments data.
* Transactional guarantees to maintain referential and financial integrity.
* Simple but effective support for idempotent writes using unique keys.
* Native connection pooling via Go's `database/sql` package.

I implemented full retry loops for MySQL to handle cold starts where DB containers may take extra time to fully initialize.

---

### 5️ Retry and Backoff Strategy

One of my priorities was to ensure resilience during temporary service outages, which often happen in distributed deployments (especially in cloud-native environments). Therefore:

* Both RabbitMQ and MySQL connection logic implement retry loops with exponential backoff.
* Startup retries ensure that transient network failures don’t bring down the entire system immediately.
* After exhausting retries, the app fails gracefully — allowing orchestrators (like Kubernetes or ECS) to handle restarts.

This pattern is something I've seen work reliably when designing cloud-native ingestion systems.

---

### 6️ Centralized Logging Abstraction

I built a centralized `logger` package which wraps Go’s standard logger:

* It logs file names and line numbers for easier debugging.
* Keeps logs consistent across publisher, consumer, and database layers.
* The design allows me to easily replace the logger later with fully structured logging engines like `zap` or `zerolog` without major refactoring.

I always build systems expecting observability to be a first-class concern, even in early stages.

---

### 7️ Environment-Driven Configuration

I deliberately separated all sensitive and runtime configs into `.env` files (loaded via `godotenv`), which allows me to:

* Avoid hardcoding credentials directly in code.
* Simplify configuration overrides across environments (local, staging, production).
* Support future secret management integrations (AWS Secrets Manager, Vault, etc).

In production, I would externalize all configs entirely from containers using environment providers.

---

### 8️ Automation-First Development

One of my non-negotiable principles is full automation from day one. The `Makefile`:

* Orchestrates full build, run, test, deploy, and teardown steps.
* Minimizes human errors.
* Simplifies developer onboarding.
* Keeps my local dev environment very close to production behavior.

This approach drastically reduces friction for CI/CD pipelines as well.

---

##  Testing Strategy

* Isolated unit tests for publisher and consumer logic (under /test)

* Validate DB insert correctness, retry mechanisms, duplicate protection

* End-to-end integration tested through full pipeline execution

---


##  Minimal Local Test Plan

This section provides a complete step-by-step local validation plan that fully covers all test cases and business logic requirements for the Fast Track Data Ingestion assignment.

---

###  1. Start the System

Bring up the entire stack:

```bash
make up
```

Expected log highlights:

```text
Successfully connected to MySQL
RabbitMQ queue declared
Published message: {UserID:1 PaymentID:1 DepositAmount:10}
Published message: {UserID:1 PaymentID:2 DepositAmount:20}
Published message: {UserID:2 PaymentID:3 DepositAmount:20}
Inserted payment: {UserID:1 PaymentID:1 DepositAmount:10}
Inserted payment: {UserID:1 PaymentID:2 DepositAmount:20}
Inserted payment: {UserID:2 PaymentID:3 DepositAmount:20}
Duplicate message published for testing
Duplicate detected, moved to skipped_messages: {UserID:1 PaymentID:1 DepositAmount:10}
```

---

###  2. Validate MySQL Tables & Initial Data

Connect to MySQL:

```bash
docker exec -it mysql mysql -uroot -proot fasttrack
```

#### Verify Tables Exist

```sql
SHOW TABLES;
```

Expected output:

```text
+---------------------+
| Tables_in_fasttrack |
+---------------------+
| payment_events      |
| skipped_messages    |
+---------------------+
```

#### Verify `payment_events` Contents

```sql
SELECT * FROM payment_events;
```

Expected:

```text
+---------+------------+----------------+
| user_id | payment_id | deposit_amount |
+---------+------------+----------------+
|       1 |          1 |             10 |
|       1 |          2 |             20 |
|       2 |          3 |             20 |
+---------+------------+----------------+
```

#### Verify `skipped_messages` Contents

```sql
SELECT * FROM skipped_messages;
```

Initially:

```text
+---------+------------+----------------+
| user_id | payment_id | deposit_amount |
+---------+------------+----------------+
|       1 |          1 |             10 |
+---------+------------+----------------+
```

 *Note:* The first duplicate (payment\_id `1`) is automatically redirected here during startup.

---

###  3. Manual Duplicate Insertion via RabbitMQ UI

Login to RabbitMQ Management UI:
 `http://localhost:15672` (default: `guest` / `guest`)

Under `payments` queue ➔ **Publish Message** ➔ enter:

```json
{
  "user_id": 1,
  "payment_id": 1,
  "deposit_amount": 10
}
```

#### Expected Application Logs:

* If duplicate already exists in `skipped_messages`:

```bash
Failed to insert into skipped_messages: Error 1062 (23000): Duplicate entry '1' for key 'skipped_messages.PRIMARY'
```

* If skipped\_messages was empty before:

```bash
Duplicate detected, moved to skipped_messages: {UserID:1 PaymentID:1 DepositAmount:10}
```

---

###  4. Insert Unique Payment via RabbitMQ UI

Again publish via UI:

```json
{
  "user_id": 20,
  "payment_id": 12,
  "deposit_amount": 13
}
```

#### Expected behavior:

* Inserted successfully into `payment_events`.
* `skipped_messages` remains unchanged.

#### Expected log:

```bash
Inserted payment: {UserID:20 PaymentID:12 DepositAmount:13}
```

#### Verify MySQL:

```sql
SELECT * FROM payment_events;
```

Should now include:

```text
|      20 |         12 |             13 |
```

---

###  5. Insert New Duplicate (new payment\_id clash)

Now publish:

```json
{
  "user_id": 2,
  "payment_id": 2,
  "deposit_amount": 11
}
```

#### Expected behavior:

* `payment_events` rejects due to primary key conflict.
* Inserted into `skipped_messages` if not already present.

#### Expected log:

```bash
Duplicate detected, moved to skipped_messages: {UserID:2 PaymentID:2 DepositAmount:11}
```

#### MySQL Verification:

```sql
SELECT * FROM skipped_messages;
```

Expected:

```text
+---------+------------+----------------+
| user_id | payment_id | deposit_amount |
+---------+------------+----------------+
|       1 |          1 |             10 |
|       2 |          2 |             11 |
+---------+------------+----------------+
```

---

###  6. Graceful Shutdown Test

While app is running, press:

```bash
CTRL + C
```

Expected log:

```bash
Shutting down gracefully.
```

 All services shut down cleanly.

---

###  7. Health Check Validation

Verify app liveness endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:

```text
OK
```

---

###  Summary Test Coverage Table

| Test Scenario                    | Expected Behavior                                    |
| -------------------------------- | ---------------------------------------------------- |
| Initial seed ingestion           | Inserts 3 valid events into `payment_events`.        |
| Initial duplicate (1,1,10)       | Inserted into `skipped_messages` automatically.      |
| Manual duplicate via UI (1,1,10) | Skipped\_messages insertion conflict (PK violation). |
| New duplicate test (2,2,11)      | Inserted into `skipped_messages`.                    |
| New valid message (20,12,13)     | Successfully inserted into `payment_events`.         |
| Shutdown (CTRL+C)                | Graceful shutdown log shown.                         |
| Health check (`/health`)         | Returns `OK`.                                        |

---

 This **Minimal Local Test Plan** fully validates correctness of:

* Functional requirements 
* Duplicate handling 
* Fault tolerance 
* End-to-end flow 

---
##  Environment Variables

All configurations are loaded via `.env` or system variables.

| Variable        | Description             | Default                                   |
| --------------- | ----------------------- | ------------------------------------------|
| MYSQL_USER      | MySQL Username          | root                                      |
| MYSQL_PASSWORD  | MySQL Password          | root                                      |
| MYSQL_HOST      | MySQL Host              | mysql                                     |
| MYSQL_PORT      | MySQL Port              | 3306                                      |
| MYSQL_DATABASE  | MySQL DB Name           | fasttrack                                 |
| RABBITMQ_URI    | RabbitMQ Connection URI | amqp://guest:guest@rabbitmq:5672/         |

---
> Note: Docker Compose automatically uses container service names for host resolution (mysql, rabbitmq).

##  Makefile Commands

| Command            | Description                              |
| ------------------ | ---------------------------------------- |
| `make up`          | Build and start entire Docker stack      |
| `make start`       | Start existing containers (no rebuild)   |
| `make rebuild-app` | Rebuild only app container               |
| `make logs`        | View live logs                           |
| `make down`        | Stop all containers                      |
| `make clean`       | Full cleanup (volumes, networks, images) |
| `make build`       | Build Go binary locally                  |
| `make run-local`   | Run app locally (non-Docker)             |
| `make test`        | Run unit tests                           |

---

##  Deployment Flow (Local Dev)

```bash
# Start full stack
make up

# View app logs
make logs

# Tear down
make down

# Cleanup everything
make clean
```

---

##  Failure Handling

* Full retry logic if MySQL or RabbitMQ are temporarily unavailable during startup.

* Graceful application termination when retries are exhausted.

* Docker Compose orchestrates service startup order via healthchecks to minimize cold start failures.

* Idempotent database writes prevent data duplication during retry scenarios.

---

## Duplicate Handling: MySQL Error 1062 Explained

When inserting records into MySQL, any violation of the `PRIMARY KEY` constraint results in error code `1062`:

```text
Error 1062 (23000): Duplicate entry 'X' for key 'PRIMARY'
```

In my implementation:

* The `payment_events` table has `payment_id` as the primary key to enforce uniqueness.
* While consuming messages from RabbitMQ, I first try to insert the payload directly into `payment_events`.
* If MySQL raises `Error 1062` (indicating a duplicate `payment_id`), I explicitly catch this error using MySQL’s error code via Go's MySQL driver (`github.com/go-sql-driver/mysql`).
* Once I detect the duplicate, I redirect that message into the `skipped_messages` table for tracking.
* If the same duplicate arrives again (already present in `skipped_messages`), MySQL again raises `1062`, which I log safely without breaking the consumer.

By directly checking for MySQL error code `1062`, I ensure:

* Idempotent inserts.
* Clean and predictable handling of duplicates.
* Zero consumer crashes even when duplicate messages are reprocessed multiple times.

This approach gives the ingestion system resilience, correctness, and production-grade safety under real-world streaming conditions.

---

##  Future Improvements

| Planned Feature       | Description                                      |
| --------------------- | ------------------------------------------------ |
| Dead Letter Queues    | Add RabbitMQ DLQ for failed message processing   |
| Observability         | Metrics via Prometheus + Grafana                 |
| CI/CD Integration     | Automated tests + deployments via GitHub Actions |
| Graceful Shutdown     | OS signal interception for safe container exit   |
| Structured Logging    | JSON logs via zap/zerolog                        |
| Distributed Tracing   | OpenTelemetry instrumentation                    |
| Secrets Management    | Externalized config using Vault providers        |

---

## Author Reflections

This system reflects the real-world ingestion architectures I have implemented professionally, with a strong focus on correctness, operational stability, simplicity, and future scalability.

The goal is to:

* Build once, run anywhere.

* Onboard junior engineers quickly.

* Support feature growth without architectural rewrites.

* Fail gracefully but recover automatically.

---

##  Security Notes

* All credentials are loaded from `.env` files to separate code from secrets.
* In production, secrets should be handled via vault providers (AWS Secrets Manager, HashiCorp Vault, etc).

---


