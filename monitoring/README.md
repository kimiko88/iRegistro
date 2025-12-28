# Monitoring Stack - Setup Instructions

## Prerequisites

- Docker and Docker Compose installed
- User must be in the `docker` group OR use `sudo`

## Directory Structure

```
monitoring/
├── docker-compose.yml
├── prometheus.yml
├── alerts.yml
├── alertmanager.yml
├── k8s-prometheus-config.yaml
├── dashboards/
│   ├── dashboards.yaml
│   ├── system-health.json
│   ├── api-metrics.json
│   ├── business-kpis.json
│   └── database.json
└── datasources/
    └── prometheus.yaml
```

## Quick Start

### Option 1: With Docker (if you have permissions)

```bash
cd monitoring
docker-compose up -d
```

### Option 2: With sudo (if permission denied)

```bash
cd monitoring
sudo docker-compose up -d
```

### Option 3: Add user to docker group (recommended for development)

```bash
sudo usermod -aG docker $USER
newgrp docker  # or logout/login
```

## Verify Setup

```bash
# Check containers are running
docker-compose ps

# Check Prometheus targets
curl http://localhost:9090/api/v1/targets

# Access Grafana
# Open browser: http://localhost:3000
# Login: admin / admin
```

## Troubleshooting

### "Permission denied" error
This means your user doesn't have access to Docker socket. Use one of the solutions above.

### "monitoring directory not found"
Ensure you're in the project root directory (`/home/k/Documenti/GitHub/iRegistro`)

### Metrics not appearing in Prometheus
1. Ensure the API service is running on port 8080
2. Check `/metrics` endpoint: `curl http://localhost:8080/metrics`
3. Update `prometheus.yml` targets if using different ports

## Environment Variables

Create a `.env` file in the monitoring directory:

```env
GRAFANA_PASSWORD=your_secure_password
DB_USER=postgres
DB_PASSWORD=your_db_password
DB_NAME=registro
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

## Next Steps

After starting the monitoring stack:
1. Access Grafana at http://localhost:3000
2. Dashboards are auto-provisioned in the "Default" folder
3. Configure alert notifications in Alertmanager
4. Review and customize alert thresholds in `alerts.yml`
