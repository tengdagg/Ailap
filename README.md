# AILAP (AI Log Analysis Platform)

[English](./README.md) | [中文](./README_zh-CN.md)

AILAP is a modern log analysis platform built with Go and Vue 3, designed to provide unified log querying, analysis, and AI-assisted diagnostics.

![Image text](https://github.com/tengdagg/Ailap/blob/main/img/img.png)

## ✨ Features

### 📊 Log Management
- **Multi-Source Support**: Supports Loki (LogQL) and Elasticsearch (Lucene) data sources.
- **Structured Display**: Automatically parses Nginx/Access Logs, extracting key fields (IP, Latency, Status Code, etc.).
- **Smart Highlighting**: Automatically color-codes logs based on HTTP status codes (2xx Green, 4xx Orange, 5xx Red).
- **Advanced Filtering**: Supports fuzzy search and value filtering for fields like Source, Method, Status, etc.

### 🤖 AI-Assisted Diagnostics
- **Intelligent Analysis**: Integrates AI models for deep analysis of log content to quickly pinpoint root causes of anomalies.
- **Interactive Chat**: Provides a ChatGPT-like interface for context-aware questions about specific logs.
- **Dynamic Avatar**: The AI assistant's avatar automatically syncs with the currently configured default model's icon.

### 🔌 Models & Data Sources
- **Model Management**: Supports configuration for various LLM interfaces like OpenAI, Deepseek, Qwen, etc.
- **Data Source Configuration**: Visual management of log data sources with connectivity testing.

## 🛠 Tech Stack

- **Frontend**: Vue 3, Vite, Arco Design Vue, Pinia
- **Backend**: Go (Gin), GORM (SQLite), Viper, Zap
- **Deployment**: Docker Multi-stage Build

## 🚀 Quick Start

### Prerequisites
- Docker (Recommended)
- Or Node.js 18+ & Go 1.21+

### Docker Deployment (Recommended)

```bash
# Build the image
docker build -t ailap .

# Run the container (with data persistence)
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  --name ailap \
  ailap
```

Access the application at `http://localhost:8080`.

### Local Development

**Backend**
```bash
cd backend
go mod tidy
go run ./cmd
# Default port: 8080
```

**Frontend**
```bash
cd frontend
npm install
npm run dev
# Access at http://localhost:5173
```

## 🔐 Default Credentials

- **Username**: `admin`
- **Password**: `admin123`

Can be overridden via `AILAP_ADMIN_USER` and `AILAP_ADMIN_PASS` environment variables during Docker startup.
