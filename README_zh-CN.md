# AILAP (AI Log Analysis Platform)

[English](./README.md) | [中文](./README_zh-CN.md)

AILAP 是一个基于 Go 和 Vue 3 的现代化日志分析平台，旨在提供统一的日志查询、分析和 AI 辅助诊断功能。

![Image text](https://github.com/tengdagg/Ailap/blob/main/img/img.png)

## ✨ 功能特性

### 📊 日志管理
- **多数据源支持**：支持 Loki (LogQL) 和 Elasticsearch (Lucene) 数据源。
- **结构化展示**：自动解析 Nginx/Access Log，提取关键字段（IP、耗时、状态码等）。
- **智能高亮**：根据 HTTP 状态码自动标记颜色（2xx 绿色, 4xx 橙色, 5xx 红色）。
- **高级筛选**：支持对源地址、Method、Status 等字段的模糊搜索和值筛选。

### 🤖 AI 辅助诊断
- **智能分析**：集成 AI 模型对日志内容进行深度分析，快速定位异常根因。
- **交互式对话**：提供类似 ChatGPT 的对话界面，支持针对特定日志的上下文提问。
- **动态头像**：AI 助手头像自动同步当前使用的默认模型图标。

### 🔌 模型与数据源
- **模型管理**：支持配置 OpenAI、Deepseek、Qwen 等多种大模型接口。
- **数据源配置**：可视化管理日志数据源，支持连通性测试。

## 🛠 技术栈

- **前端**：Vue 3, Vite, Arco Design Vue, Pinia
- **后端**：Go (Gin), GORM (SQLite), Viper, Zap
- **部署**：Docker 多阶段构建

## 🚀 快速开始

### 前置要求
- Docker (推荐)
- 或 Node.js 18+ & Go 1.21+

### Docker 部署 (推荐)

```bash
# 构建镜像
docker build -t ailap .

# 运行容器 (数据持久化)
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  --name ailap \
  ailap
```

访问 `http://localhost:8080` 即可使用。

### 本地开发

**后端**
```bash
cd backend
go mod tidy
go run ./cmd
# 默认端口: 8080
```

**前端**
```bash
cd frontend
npm install
npm run dev
# 访问 http://localhost:5173
```

## 🔐 默认账号

- **用户名**: `admin`
- **密码**: `admin123`

可通过环境变量 `AILAP_ADMIN_USER` 和 `AILAP_ADMIN_PASS` 在 Docker 启动时修改。
