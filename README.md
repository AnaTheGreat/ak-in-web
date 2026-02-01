//################//
// File Structure //
//################//

ak-in-web/
│
├── README.md
├── Makefile
├── .gitignore +
├── .env.example +
├── docker-compose.yml
│
├── frontend/
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── package.json
│   ├── vite.config.ts
│   │
│   ├── public/
│   │   ├── covers/                # Book & movie cover images
│   │   └── resume.pdf
│   │
│   └── src/
│       ├── telemetry/             # Datadog RUM, error tracking
│       ├── api/                   # Typed API clients
│       │   ├── admin.ts
│       │   ├── library.ts
│       │   └── cinematheque.ts
│       │
│       ├── auth/                  # Admin auth logic
│       │   ├── AdminModal.tsx
│       │   └── useAdminAuth.ts
│       │
│       ├── components/
│       │   ├── layout/
│       │   │   ├── Header.tsx
│       │   │   ├── FolderTree.tsx
│       │   │   └── ScrollToTop.tsx
│       │   │
│       │   ├── cards/
│       │   │   ├── BookCard.tsx
│       │   │   └── MovieCard.tsx
│       │   │
│       │   └── admin/
│       │       ├── BookForm.tsx
│       │       ├── MovieForm.tsx
│       │       └── TagManager.tsx
│       │
│       ├── pages/
│       │   ├── Home.tsx
│       │   ├── Portfolio.tsx
│       │   ├── Education.tsx
│       │   ├── Experience.tsx
│       │   ├── Library.tsx
│       │   └── Cinematheque.tsx
│       │
│       ├── hooks/
│       ├── styles/
│       ├── App.tsx
│       └── main.tsx
│
├── backend/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   │
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── config/
│   │   ├── telemetry/
│   │   │   ├── setup.go
│   │   │   └── middleware/
│   │   │
│   │   ├── auth/
│   │   │   └── admin.go
│   │   │
│   │   ├── library/
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   └── repository.go
│   │   │
│   │   ├── cinematheque/
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   └── repository.go
│   │   │
│   │   └── router.go
│   │
│   ├── migrations/
│   │   ├── 001_init.sql
│   │   ├── 002_library.sql
│   │   └── 003_cinematheque.sql
│   │
│   └── config/
│       └── config.yaml.example
│
├── infrastructure/
│   ├── terraform/
│   │   ├── main.tf
│   │   ├── networking.tf
│   │   ├── ecs.tf
│   │   ├── rds.tf
│   │   ├── datadog.tf
│   │   └── outputs.tf
│   │
│   └── scripts/
│       ├── deploy.sh
│       └── setup-datadog.sh
│
├── .github/
│   └── workflows/
│       ├── ci.yml
│       ├── infrastructure.yml
│       └── deploy.yml
│
└── docs/
    ├── ARCHITECTURE.md
    ├── DEVELOPMENT.md
    └── DEPLOYMENT.md
//#############################################################################################//

- .gitignore enforces repo hygiene and prevents state or secrets from leaking, while .env.example defines a strict environment contract so the app is reproducible across local, CI, and AWS.

- I start with a dependency-free /health endpoint so infra and CI can validate the service before any business logic exists. Public read APIs come next, then auth, then admin write paths, and finally readiness checks once the system is production-bound.

- Next...

//#############################################################################################//

Phase 1 — Backend API Skeleton

Created first, because frontend depends on it

Why: establishes domain boundaries

Files created:

backend/cmd/api/main.go

internal/router.go

internal/config/

internal/telemetry/setup.go

At this stage: /health endpoint only.

Phase 2 — Domain Separation (Key Design Signal)

You explicitly model Library and Cinematheque as non-interacting bounded contexts.

internal/library/
internal/cinematheque/


Each has:

handler.go → HTTP layer

service.go → business rules (ratings 0–10, tag validation)

repository.go → SQL only

💼 This screams “senior design instincts”.

Phase 3 — Auth (Admin Only)

Created after domain APIs

internal/auth/admin.go

Why:

Single admin user

Stateless JWT or session cookie

No user table → intentional simplicity

Phase 4 — Database & Migrations

Created once API contracts stabilize

001_init.sql → users/admin

002_library.sql

003_cinematheque.sql

Explicit separation guarantees no accidental coupling.

Phase 5 — Frontend Shell

Created after API contracts exist

App.tsx

FolderTree.tsx

Header.tsx

Page routing

Why:

Folder-tree UI drives navigation state

URL mirrors directory structure

Phase 6 — Admin UX

Created after read-only views work

AdminModal.tsx

BookForm.tsx

MovieForm.tsx

TagManager.tsx

Admin capabilities are feature-flagged by auth, not separate app.

Phase 7 — Observability

Added after app works locally

Frontend:

Datadog RUM

Error tracking

Backend:

OpenTelemetry middleware

Request traces

SQL latency

This timing matters. You don’t instrument broken systems.

Phase 8 — Docker & Local Parity

Dockerfile (multi-stage)

docker-compose.yml

Why:

Local ≈ prod

CI uses same artifacts

Phase 9 — Infrastructure as Code

Only after app is containerized

ECS Fargate

ALB

RDS

Datadog agent

Phase 10 — CI/CD

Last step

ci.yml → test + build

infrastructure.yml → terraform plan/apply

deploy.yml → ECS rollout + smoke test

//#############################################################################################//
// plan //
//#############################################################################################//

## 📋 **PHASE 0 — DEFINE THE CONTRACT (1 HOUR)**
**Objective:** Lock architectural decisions to prevent rework

### ✅ **0.1 — Architecture Decision Document**
- [ ] Single-page TODO app with CRUD operations
- [ ] REST API backend with Go
- [ ] React/TypeScript frontend
- [ ] PostgreSQL database
- [ ] AWS ECS Fargate deployment
- [ ] Datadog for observability

### ✅ **0.2 — README Foundation**
```markdown
# Golden Path (Locked)
- Backend: Go, Chi router, structured logging
- Frontend: React + TypeScript, Vite
- Database: PostgreSQL (SQLite for initial dev)
- Deployment: AWS ECS Fargate + ALB
- Observability: Datadog (APM, RUM, Logs, Metrics)
- CI/CD: GitHub Actions

# Non-Goals (Explicit)
- Multi-region deployment
- Kubernetes orchestration
- Advanced authentication (JWT only)
- Mobile apps
- Real-time collaboration
```

### ✅ **0.3 — API Contract**
- [ ] Define OpenAPI spec (Swagger)
- [ ] Endpoints: `GET/POST/PUT/DELETE /todos`
- [ ] Todo schema: `{id, title, completed, created_at}`
- [ ] Error response format

---

## 📋 **PHASE 1 — REPO SKELETON & HYGIENE (2 HOURS)**
**Objective:** Create safe rails before speed

### ✅ **1.1 — Repository Structure**
```bash
todo-fullstack-aws/
├── README.md
├── .gitignore
├── .env.example
├── Makefile
├── docker-compose.yml
├── frontend/
├── backend/
└── infrastructure/
```

### ✅ **1.2 — Development Hygiene**
- [ ] `.gitignore` for all languages/tools
- [ ] `.env.example` with all required variables
- [ ] `Makefile` with: `make up`, `make test`, `make clean`
- [ ] Basic `docker-compose.yml` skeleton
- [ ] Pre-commit hooks (optional but recommended)
- [ ] GitHub repo created and initialized

---

## 📋 **PHASE 2 — BACKEND FIRST (1 DAY)**
**Objective:** Establish the system of record

### ✅ **2.1 — Minimal Backend Skeleton**
- [ ] `go mod init` with proper module name
- [ ] `cmd/api/main.go` with HTTP server
- [ ] `/health` endpoint returning 200 OK
- [ ] Chi router configured
- [ ] Graceful shutdown handler

### ✅ **2.2 — TELEMETRY FIRST (Non-negotiable)**
- [ ] Datadog tracing middleware (`dd-trace-go`)
- [ ] Structured logging with JSON output
- [ ] Request ID middleware for correlation
- [ ] Metrics stub for request count/duration
- [ ] `/health` endpoint emits trace + log

### ✅ **2.3 — TODO API (In-memory)**
- [ ] `GET /todos` - list all todos
- [ ] `POST /todos` - create new todo
- [ ] `PUT /todos/:id` - update todo
- [ ] `DELETE /todos/:id` - delete todo
- [ ] Proper HTTP status codes (200, 201, 400, 404, 500)
- [ ] Input validation
- [ ] Repository interface (in-memory implementation)

### ✅ **2.4 — Database Layer**
- [ ] Choose PostgreSQL driver (`lib/pq` or `pgx`)
- [ ] Database connection pool configuration
- [ ] Migration system (go-migrate or simple SQL files)
- [ ] `migrations/001_create_todos.up.sql`
- [ ] `migrations/001_create_todos.down.sql`
- [ ] Repository implementation with real DB
- [ ] Integration tests for repository

---

## 📋 **PHASE 3 — FRONTEND (1 DAY)**
**Objective:** Consume the stable API contract

### ✅ **3.1 — Frontend Skeleton**
- [ ] `npm create vite@latest` with React + TypeScript
- [ ] Configure absolute imports (`@/components/`)
- [ ] Environment configuration (`.env`, `.env.example`)
- [ ] Base layout component
- [ ] Routing setup (React Router if multi-page)

### ✅ **3.2 — Frontend Telemetry**
- [ ] Datadog RUM initialization
- [ ] Error boundary with error tracking
- [ ] API timing hooks/interceptors
- [ ] User action tracking (optional)
- [ ] Environment-aware configuration (dev/prod)

### ✅ **3.3 — TODO UI Implementation**
- [ ] `TodoList` component (fetches + displays todos)
- [ ] `TodoItem` component (with complete/delete)
- [ ] `CreateTodo` component (form with validation)
- [ ] API client (`todoApi.ts`) with axios/fetch
- [ ] Loading/error states
- [ ] Basic CSS/SCSS styling

### ✅ **3.4 — API Integration**
- [ ] TypeScript interfaces from OpenAPI spec
- [ ] React Query or similar for state management
- [ ] Optimistic updates for better UX
- [ ] Error handling with user-friendly messages

---

## 📋 **PHASE 4 — LOCAL ORCHESTRATION (½ DAY)**
**Objective:** Production-like local development

### ✅ **4.1 — Complete docker-compose.yml**
```yaml
services:
  postgres:
    image: postgres:15
    environment: [...]
    volumes: [...]
    
  backend:
    build: ./backend
    ports: ["8080:8080"]
    depends_on: [postgres]
    environment: [...]
    
  frontend:
    build: ./frontend
    ports: ["5173:5173"]
    depends_on: [backend]
    environment: [...]
```

### ✅ **4.2 — One-Command Development**
- [ ] `make up` → `docker-compose up --build`
- [ ] `make test` → runs all tests
- [ ] `make down` → `docker-compose down`
- [ ] `make logs` → view all service logs
- [ ] `make db-migrate` → run migrations

### ✅ **4.3 — Local Validation**
- [ ] Fresh clone can run app in < 5 minutes
- [ ] Database migrations run automatically
- [ ] All services start in correct order
- [ ] Health checks pass
- [ ] Frontend can communicate with backend

---

## 📋 **PHASE 5 — PRODUCTION CONTAINERS (½ DAY)**
**Objective:** Immutable, reproducible artifacts

### ✅ **5.1 — Multi-stage Dockerfiles**
**Backend Dockerfile:**
- [ ] Stage 1: Builder (Go compile)
- [ ] Stage 2: Distroless/minimal base image
- [ ] Non-root user
- [ ] Health check endpoint
- [ ] Proper signal handling
- [ ] Optimized layer caching

**Frontend Dockerfile:**
- [ ] Stage 1: Node builder
- [ ] Stage 2: Nginx for static files
- [ ] Optimized build (compression, cache headers)
- [ ] Security headers
- [ ] Nginx config with routing

### ✅ **5.2 — Local Production Testing**
- [ ] `docker build -t backend:prod ./backend`
- [ ] `docker build -t frontend:prod ./frontend`
- [ ] Run containers locally with production configs
- [ ] Verify health endpoints work
- [ ] Test with different environment variables

---

## 📋 **PHASE 6 — INFRASTRUCTURE (2 DAYS)**
**Objective:** Cloud deployment without surprises

### ✅ **6.1 — Terraform Foundation**
- [ ] Remote state configuration (S3 + DynamoDB)
- [ ] Provider versions locked
- [ ] Variable definitions (`variables.tf`)
- [ ] Output definitions (`outputs.tf`)
- [ ] Backend initialization works

### ✅ **6.2 — Networking Layer**
- [ ] VPC with public/private subnets
- [ ] Security groups (minimum necessary ports)
- [ ] Internet Gateway / NAT Gateway
- [ ] Route tables
- [ ] Test: Can resources communicate?

### ✅ **6.3 — Database (RDS)**
- [ ] PostgreSQL RDS instance
- [ ] Parameter group for tuning
- [ ] Automated backups configured
- [ ] Read replica (optional for v1)
- [ ] Security: encryption at rest, no public access

### ✅ **6.4 — ECS Cluster & Services**
- [ ] ECS Fargate cluster
- [ ] Task definitions for backend/frontend
- [ ] IAM roles for tasks
- [ ] Service discovery (optional)
- [ ] Auto-scaling configuration
- [ ] Test: Tasks can pull from ECR, start, pass health checks

### ✅ **6.5 — Load Balancer & Routing**
- [ ] Application Load Balancer
- [ ] Target groups (backend:8080, frontend:80)
- [ ] Listener rules for routing
- [ ] SSL certificate (ACM) for HTTPS
- [ ] DNS record (Route53) if custom domain
- [ ] Test: ALB routes traffic correctly

### ✅ **6.6 — Observability in Cloud**
- [ ] Datadog agent as sidecar container
- [ ] CloudWatch log groups
- [ ] Container insights enabled
- [ ] Metrics collection verified
- [ ] Traces visible in Datadog
- [ ] Logs flowing to Datadog

---

## 📋 **PHASE 7 — CI/CD (1 DAY)**
**Objective:** Automate what already works manually

### ✅ **7.1 — CI Pipeline (.github/workflows/ci.yml)**
- [ ] Run on PR: lint, test, build
- [ ] Backend: unit tests, go vet, gofmt
- [ ] Frontend: type check, unit tests, build
- [ ] Docker: build images, security scan (trivy)
- [ ] Push images to ECR with Git SHA tag
- [ ] Integration tests against running containers

### ✅ **7.2 — Infrastructure Pipeline (.github/workflows/terraform.yml)**
- [ ] Terraform fmt/validate on PR
- [ ] Terraform plan output in PR comments
- [ ] Manual approval for apply
- [ ] Separate workflows for dev/staging/prod
- [ ] State locking during apply
- [ ] Notifications on success/failure

### ✅ **7.3 — Deployment Pipeline (.github/workflows/deploy.yml)**
- [ ] Update ECS task definitions with new images
- [ ] Wait for service stabilization
- [ ] Run smoke tests against new deployment
- [ ] Rollback on failure
- [ ] Notify Slack/email on deployment

---

## 📋 **PHASE 8 — DOCUMENTATION & HARDENING (½ DAY)**
**Objective:** Make it explainable and resilient

### ✅ **8.1 — Architecture Documentation**
- [ ] System architecture diagram (draw.io/Excalidraw)
- [ ] Data flow diagram
- [ ] Deployment architecture
- [ ] Monitoring architecture

### ✅ **8.2 — Runbooks & Procedures**
- [ ] Deployment runbook (step-by-step)
- [ ] Rollback procedure
- [ ] Database backup/restore procedure
- [ ] Incident response checklist
- [ ] Troubleshooting guide

### ✅ **8.3 — Security Hardening**
- [ ] Secrets in AWS Secrets Manager (not env vars)
- [ ] Least-privilege IAM roles verified
- [ ] Security headers in ALB/nginx
- [ ] DDoS protection basics (WAF optional)
- [ ] Vulnerability scan results documented

### ✅ **8.4 — Performance & Cost**
- [ ] Load testing results (even basic)
- [ ] Cost estimation per environment
- [ ] Cost optimization recommendations
- [ ] Performance baseline metrics

### ✅ **8.5 — Final README**
- [ ] Quick start (local development)
- [ ] Deployment guide
- [ ] API documentation
- [ ] Monitoring dashboard links
- [ ] Team contacts
- [ ] License

---

## 🚀 **LAUNCH CHECKLIST**

### **Pre-Launch (Day Before)**
- [ ] All tests passing
- [ ] Production database backed up
- [ ] Team notified of deployment window
- [ ] Rollback procedure tested
- [ ] Monitoring dashboards ready

### **Launch Day**
- [ ] Final code review completed
- [ ] CI/CD pipelines green
- [ ] Deploy to production
- [ ] Smoke tests pass
- [ ] Monitor metrics for 30 minutes
- [ ] Announce successful deployment

### **Post-Launch (First Week)**
- [ ] Review Datadog dashboards daily
- [ ] Monitor error rates
- [ ] Check cost alerts
- [ ] Gather team feedback
- [ ] Schedule retrospective

---

## 📊 **SUCCESS METRICS (Track These)**

| Metric | Target | How to Measure |
|--------|--------|----------------|
| Deployment Time | < 10 minutes | GitHub Actions logs |
| P95 API Latency | < 100ms | Datadog APM |
| Error Rate | < 1% | Datadog Error Tracking |
| Uptime | > 99.5% | Datadog Synthetic |
| Cost | < $100/month | AWS Cost Explorer |