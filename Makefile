.PHONY: up down run-backend run-frontend test test-unit test-integration test-e2e test-security test-benchmark coverage

# Start infrastructure
up:
	docker-compose up -d

# Stop infrastructure
down:
	docker-compose down

# Run backend
run-backend:
	go run cmd/api/main.go

# Run frontend
run-frontend:
	cd frontend && npm install && npm run dev

# Run all unit tests (default)
test: test-unit

# Run unit tests
test-unit:
	go test -v -short ./internal/...

# Run integration tests (requires Docker)
test-integration:
	go test -v ./tests/integration/...

# Run E2E tests (requires running server)
test-e2e:
	go test -v ./tests/e2e/...

# Run security tests
test-security:
	go test -v ./tests/security/...

# Run benchmarks
test-benchmark:
	go test -bench=. ./tests/benchmarks/...

# Run coverage report
coverage:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run all tests (slow)
test-all: test-unit test-integration test-e2e test-security
