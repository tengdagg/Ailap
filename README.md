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







