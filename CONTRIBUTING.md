# Contributing to iRegistro

We welcome contributions to the iRegistro project! This document provides guidelines for contributing to the repository.

## Getting Started

1.  **Fork the repository** on GitHub.
2.  **Clone your fork**:
    ```bash
    git clone https://github.com/your-username/iRegistro.git
    cd iRegistro
    ```
3.  **Create a branch** for your feature or bugfix:
    ```bash
    git checkout -b feature/my-new-feature
    ```

## Development Workflow

### Backend (Go) Service
- Ensure you have Go 1.24+ installed.
- Run tests: `make test`
- Run linting: `golangci-lint run`

### Frontend (Vue 3 + Typescript)
- Navigate to `frontend/`.
- Install dependencies: `npm install --legacy-peer-deps`
- Run dev server: `npm run dev`
- Run tests: `npm run test:unit`

## CI/CD Pipeline

We use GitHub Actions for our pipeline.

- **Lint & Test**: Runs on Pull Requests. Checks Go/JS linting, unit tests, and coverage.
- **Security**: Runs nightly and on push. Checks for vulnerabilities (Gosec, Dependency Check).
- **Build & Deploy**: Runs on push to `main`. Builds containers and deploys to Staging.
- **E2E Tests**: Runs after successful deployment to Staging.

## Pull Request Process

1.  Ensure all tests pass.
2.  Update documentation if necessary.
3.  Open a Pull Request against the `main` branch.
4.  Wait for CI checks to pass and for a code review.
