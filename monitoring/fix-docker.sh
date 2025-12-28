#!/bin/bash
# Fix docker-compose ContainerConfig error

echo "🔧 Fixing Docker corrupted containers..."

# Stop all containers
sudo docker stop $(sudo docker ps -aq) 2>/dev/null

# Remove ALL containers (including corrupted ones)
sudo docker container prune -f

# Remove all dangling images
sudo docker image prune -f

# Clean up volumes (OPTIONAL - this will delete data!)
# sudo docker volume prune -f

# Restart Docker daemon to clear cache
echo "🔄 Restarting Docker daemon..."
sudo systemctl restart docker
sleep 3

# Navigate to monitoring directory
cd /home/k/Documenti/GitHub/iRegistro/monitoring

# Start fresh
echo "🚀 Starting monitoring stack..."
sudo docker-compose up -d

echo "✅ Done! Check status with: sudo docker-compose ps"
