# AILAP

Monorepo for a Vue 3 + Vite frontend and Go + Gin backend.

## Prerequisites
- Node.js 18+
- Go 1.21+
- Sqlite (optional initially)

## Frontend
```bash
cd frontend
npm install
npm run dev
```

## Backend
```bash
cd backend
go mod tidy
go run ./cmd
```

Default backend port: 8080. Vite dev server proxies `/api` to backend.

默认管理员账号：username admin，password admin123（可通过环境变量 AILAP_ADMIN_USER/AILAP_ADMIN_PASS 初始化时覆盖）。

## Docker

Build and run the application using Docker:

```bash
# Build the image
docker build -t ailap .

# Run the container with data persistence
# The database is stored in /app/data inside the container
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  --name ailap \
  ailap
```

Access the application at `http://localhost:8080`.
The data (SQLite DB) will be persisted in the `./data` directory on your host.

### Custom Initial Admin Credentials
You can set the initial admin username and password using environment variables:

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -e AILAP_ADMIN_USER=myadmin \
  -e AILAP_ADMIN_PASS=mypassword \
  --name ailap \
  ailap
```















