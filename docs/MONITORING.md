# Monitoring & Observability

This document describes the monitoring stack for iRegistro production environment.

## Architecture

The monitoring stack consists of:
- **Prometheus**: Metrics collection and storage
- **Grafana**: Visualization and dashboarding
- **Alertmanager**: Alert routing and notification
- **Postgres Exporter**: Database metrics

## Quick Start

### Local Development

```bash
cd monitoring
docker-compose up -d
```

Access points:
- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090  
- **Alertmanager**: http://localhost:9093

### Kubernetes Production

```bash
kubectl apply -f monitoring/k8s-prometheus-config.yaml
```

## Metrics Exposed

### HTTP/API Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `http_requests_total` | Counter | Total HTTP requests |
| `http_request_duration_seconds` | Histogram | Request latency distribution |
| `http_request_errors_total` | Counter | HTTP errors (4xx/5xx) |
| `active_websocket_connections` | Gauge | Active WebSocket connections |

### Database Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `db_query_duration_seconds` | Histogram | Query execution time |
| `db_connection_pool_size` | Gauge | Connection pool size |
| `db_slow_queries_total` | Counter | Queries >500ms |

### Business Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `marks_added_total` | Counter | Total marks added |
| `absences_recorded_total` | Counter | Total absences recorded |
| `documents_generated_total` | Counter | Total documents generated |
| `documents_generation_errors_total` | Counter | Document generation failures |
| `users_login_total` | Counter | Total user logins |
| `unauthorized_access_attempts_total` | Counter | **Security**: unauthorized access |

### System Metrics

Automatically collected by Prometheus Go client:
- `go_goroutines`
- `go_memstats_alloc_bytes`
- `process_cpu_seconds_total`

## Dashboards

### 1. System Health
**File**: `dashboards/system-health.json`

Monitors:
- CPU usage
- Memory allocation
- Goroutine count
- Request rate

### 2. API Metrics
**File**: `dashboards/api-metrics.json`

Monitors:
- P95 request duration
- Error rates
- Throughput by method
- Status code distribution
- WebSocket connections

### 3. Business KPIs
**File**: `dashboards/business-kpis.json`

Monitors:
- Marks added per hour
- Absences recorded per hour
- Documents generated
- User logins
- **Security violations**

### 4. Database
**File**: `dashboards/database.json`

Monitors:
- Query duration
- Connection pool size
- Slow queries

## Alert Rules

### Critical Alerts

| Alert | Condition | Threshold |
|-------|-----------|-----------|
| `DatabaseDown` | Database unreachable | 1 minute |
| `MultiTenantViolation` | Unauthorized access detected | Immediate |

### Warning Alerts

| Alert | Condition | Threshold |
|-------|-----------|-----------|
| `HighErrorRate` | Error rate > 5% | 5 minutes |
| `SlowResponses` | P95 latency > 1s | 10 minutes |
| `OutOfMemory` | Memory > 900MB | 5 minutes |
| `HighWebSocketConnections` | WS connections > 1000 | 5 minutes |

### Info Alerts

| Alert | Condition | Threshold |
|-------|-----------|-----------|
| `NoMarksAdded` | No marks in 2 hours | Informational |
| `DocumentGenerationFailed` | Failure rate > 10% | 5 minutes |

## Alert Destinations

Configure in `alertmanager.yml`:

```yaml
receivers:
  - name: 'critical-alerts'
    slack_configs:
      - channel: '#critical-alerts'
    # Add PagerDuty for production
```

## Security Monitoring

### Multi-Tenant Isolation

The `unauthorized_access_attempts_total` metric tracks attempts to access data from wrong schools. This triggers the `MultiTenantViolation` alert immediately.

**Investigation Steps**:
1. Check Grafana Business KPIs dashboard
2. Query Prometheus: `unauthorized_access_attempts_total`
3. Review application logs for user_id and school_id
4. Escalate to security team if confirmed breach

## Query Examples

### Top 10 Slowest Endpoints
```promql
topk(10, histogram_quantile(0.95, 
  rate(http_request_duration_seconds_bucket[5m])
))
```

### Error Rate by Endpoint
```promql
rate(http_request_errors_total[5m]) / 
rate(http_requests_total[5m])
```

### Daily Active Users
```promql
increase(users_login_total[24h])
```

## Troubleshooting

### Metrics Not Appearing

1. Verify `/metrics` endpoint: `curl http://localhost:8080/metrics`
2. Check Prometheus targets: http://localhost:9090/targets
3. Verify middleware is registered in router

### Alerts Not Firing

1. Check alert rules: http://localhost:9090/alerts
2. Verify Alertmanager config
3. Test notification channel (Slack webhook)

### High Memory Usage

Check Go heap profile:
```bash
go tool pprof http://localhost:8080/debug/pprof/heap
```

## Production Best Practices

1. **Retention**: Configure Prometheus retention (default 15 days)
2. **Backup**: Enable persistent volumes for Prometheus/Grafana
3. **Security**: Protect `/metrics` endpoint with authentication
4. **Alerts**: Route critical alerts to PagerDuty/OpsGenie
5. **Dashboards**: Set up separate dashboards per staging/prod environment

## Scaling

For multi-instance deployments:
- Use Prometheus federation
- Add `instance` label to all metrics
- Configure service discovery (Kubernetes)

## References

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)
- [Go Client Best Practices](https://prometheus.io/docs/instrumenting/clientlibs/)
