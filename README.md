# iRegistro - Electronic Registry System

iRegistro is a modern, full-stack electronic registry system designed for Italian schools. It manages academic records, attendance, communications, and administrative workflows.

## Features

- **Multi-Role Access**: Dedicated dashboards for Admin, Director, Secretary, Teacher, Parent, and Student.
- **Academic Management**: Classes, Subjects, Marks, and Absences.
- **Document Workflow**: Digital signature and approval process for official documents.
- **Communications**: Secure messaging and colloquium booking.

## Technology Stack

- **Backend**: Go (Golang) 1.24, PostgreSQL 15, Redis. Correct Clean Architecture.
- **Frontend**: Vue 3, Pinia, TypeScript, TailwindCSS.
- **Infrastructure**: Docker, Kubernetes, GitHub Actions CI/CD.

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.24+
- Node.js 18+

### Setup

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/your-org/iRegistro.git
    cd iRegistro
    ```
2.  **Start Infrastructure**:
    ```bash
    make up
    ```
    This spins up Postgres, Redis, and other dependencies.
3.  **Run Backend**:
    ```bash
    go run cmd/api/main.go
    ```
4.  **Run Frontend**:
    ```bash
    cd frontend
    npm install
    npm run dev
    ```

## Development & Testing

We follow TDD principles. Run the test suite:

```bash
make test        # Run all backend tests
make test-e2e   # Run E2E tests
npm run test:unit # Run frontend unit tests
```

See [docs/TESTING.md](docs/TESTING.md) for more details.

## CI/CD Workflows

We utilize GitHub Actions for our continuous integration and deployment pipeline:

| Workflow | Trigger | Description |
| :--- | :--- | :--- |
| **Lint & Test** | Push/PR | Runs Go/JS linting, Unit Tests, and Coverage checks. Fails if coverage < 85%. |
| **Security Scans** | Daily/Push | Performs Gosec (Go SAST), OWASP ZAP (DAST), and Dependency Checks. |
| **Build & Deploy** | Push to Main | Builds Docker images, pushes to GHCR, and deploys to Staging K8s. |
| **E2E Tests** | Post-Deploy | Runs Playwright E2E tests against the Staging environment. |

## Contribution

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.
