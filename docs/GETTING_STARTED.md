# Come Avviare iRegistro

## Backend (Go API)

### 1. Configura il Database
```bash
# Avvia PostgreSQL (se non già running)
sudo systemctl start postgresql

# Oppure usa Docker
docker run -d --name registro-db \
  -e POSTGRES_USER=registro \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=registro \
  -p 5432:5432 \
  postgres:15
```

### 2. Configura le Variabili d'Ambiente
```bash
cd /home/k/Documenti/GitHub/iRegistro

# Copia il file di esempio
cp .env.example .env

# Modifica .env con i tuoi valori
nano .env
```

**Variabili minime richieste:**
```env
DATABASE_URL=postgres://registro:password@localhost:5432/registro?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-min-32-characters-long
PORT=8080
FRONTEND_URL=http://localhost:3000
```

### 3. Esegui le Migrazioni
```bash
# Compila il tool di migrazione
go build -o bin/migrate ./cmd/migrate

# Esegui tutte le migrazioni
./bin/migrate up
```

### 4. Avvia il Backend
```bash
# Opzione 1: Con go run
go run ./cmd/api/main.go

# Opzione 2: Compila ed esegui
go build -o bin/api ./cmd/api
./bin/api

# Il server sarà disponibile su http://localhost:8080
```

**Verifica che sia running:**
```bash
curl http://localhost:8080/health
# Dovresti vedere: {"status":"ok"}
```

---

## Frontend (Vue.js)

### 1. Installa le Dipendenze
```bash
cd /home/k/Documenti/GitHub/iRegistro/frontend

# Installa i pacchetti npm
npm install
```

### 2. Configura le Variabili d'Ambiente
```bash
# Crea file .env.local
cat > .env.local << EOF
VITE_API_URL=http://localhost:8080
EOF
```

### 3. Avvia il Dev Server
```bash
# Avvia in modalità development
npm run dev

# Il frontend sarà disponibile su http://localhost:3000
```

**Oppure build per produzione:**
```bash
npm run build
npm run preview
```

---

## Avvio Completo con Monitoring

### Usa Docker Compose per tutto
```bash
cd /home/k/Documenti/GitHub/iRegistro

# Crea docker-compose.yml principale (se non esiste)
# Avvia tutti i servizi
docker-compose up -d

# Oppure solo monitoring
cd monitoring
sudo docker-compose up -d
```

---

## Verifica che Tutto Funzioni

### 1. Backend API
```bash
# Health check
curl http://localhost:8080/health

# Metrics (Prometheus)
curl http://localhost:8080/metrics
```

### 2. Frontend
Apri il browser: **http://localhost:3000**

### 3. Monitoring Stack
- **Grafana**: http://localhost:3000 (se monitoring è in porta diversa: 3001)
- **Prometheus**: http://localhost:19090
- **Postgres Exporter**: http://localhost:9187/metrics

---

## Quick Start (Tutto in Un Colpo)

```bash
# Terminal 1: Backend
cd /home/k/Documenti/GitHub/iRegistro
./bin/migrate up
go run ./cmd/api/main.go

# Terminal 2: Frontend  
cd /home/k/Documenti/GitHub/iRegistro/frontend
npm run dev

# Terminal 3: Monitoring (opzionale)
cd /home/k/Documenti/GitHub/iRegistro/monitoring
sudo docker-compose up -d
```

---

## Troubleshooting

### Database non si connette
```bash
# Verifica che PostgreSQL sia running
sudo systemctl status postgresql

# Testa la connessione
psql -U registro -h localhost -d registro
```

### Porta già in uso
```bash
# Trova cosa usa la porta 8080
sudo lsof -i :8080

# Oppure cambia porta nel .env
PORT=8081
```

### Frontend non trova API
```bash
# Verifica che VITE_API_URL sia corretto
cat frontend/.env.local

# Dev tools del browser > Network tab
# Controlla le chiamate API
```

### Migrazioni falliscono
```bash
# Drop e ricrea database
dropdb registro
createdb registro

# Riesegui migrazioni
./bin/migrate up
```

---

## Accesso Iniziale

**Crea un utente test:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@scuola.it",
    "password": "password123",
    "role": "Admin",
    "school_id": 1,
    "first_name": "Admin",
    "last_name": "Test"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@scuola.it",
    "password": "password123"
  }'
```

Salva il token JWT ricevuto e usalo per le chiamate autenticate:
```bash
curl http://localhost:8080/api/protected \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Logs e Debug

### Backend Logs
```bash
# Con go run
go run ./cmd/api/main.go 2>&1 | tee api.log

# Se compilato
./bin/api 2>&1 | tee api.log
```

### Frontend Dev Logs
I log sono automaticamente nel terminal dove hai eseguito `npm run dev`

### Database Queries
Attiva query logging nel .env:
```env
LOG_LEVEL=debug
```

Buon divertimento! 🚀
