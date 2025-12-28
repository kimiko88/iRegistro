# Backend Testing Guide

This project follows a strict TDD approach with multiple layers of testing.

## Test Suites

### 1. Unit Tests (`test-unit`)
- **Location**: `internal/application/**/*_test.go`
- **Scope**: Service layer logic.
- **Dependencies**: Everything is mocked (Repositories, Notifiers).
- **Run**: `make test-unit`

### 2. Integration Tests (`test-integration`)
- **Location**: `tests/integration/`
- **Scope**: Repository layer and Database interactions.
- **Dependencies**: Uses `testcontainers-go` to spin up a real PostgreSQL 15 instance.
- **Run**: `make test-integration` (Requires Docker)

### 3. E2E Tests (`test-e2e`)
- **Location**: `tests/e2e/`
- **Scope**: Full API testing against a running server.
- **Run**: `make test-e2e`

## Test Data
- **Seed**: `tests/fixtures/seed_test_data.sql` contains standard data for E2E tests.

## Continuous Integration
- Tests are run automatically on PRs.
- Coverage target: > 85%.
