# 📘 AILAP 技术需求文档

## 1. 项目概述
AILAP (AI Log Analysis Platform) 是一个基于 Go 和 Vue 3 的现代化日志分析平台，旨在提供统一的日志查询、分析和 AI 辅助诊断功能。支持多种日志数据源（Loki, Elasticsearch）和本地/远程 AI 模型。

## 2. 技术架构

### 2.1 前端 (Frontend)
- **框架**：Vue 3 + Vite
- **UI 组件库**：Arco Design Vue
- **状态管理**：Pinia
- **路由**：Vue Router 4
- **HTTP 客户端**：Axios (封装拦截器，处理 JWT 和错误)
- **主要特性**：
    - 响应式布局，支持亮色/暗色主题
    - 可折叠侧边栏菜单
    - 结构化日志展示与模糊搜索
    - AI 对话交互界面

### 2.2 后端 (Backend)
- **语言**：Go 1.21+
- **Web 框架**：Gin
- **数据库**：SQLite (通过 GORM 管理)
- **配置管理**：Viper (支持环境变量覆盖)
- **日志**：Zap
- **认证**：JWT (JSON Web Token)
- **主要模块**：
    - `internal/server`: HTTP 服务启动与配置
    - `internal/router`: 路由定义与静态资源托管
    - `internal/handler`: 业务逻辑处理
    - `internal/database`: 数据库初始化与迁移
    - `internal/service`: 外部服务集成 (Loki, ES, AI)

### 2.3 部署架构
- **容器化**：Docker 多阶段构建 (Multi-stage Build)
    - Stage 1: Node.js 构建前端静态资源
    - Stage 2: Go 构建后端二进制文件
    - Stage 3: Alpine 最终镜像，整合前后端
- **单端口服务**：后端统一服务 API 和前端静态资源 (`/api` 走接口，其他走静态文件)
- **数据持久化**：通过 Docker Volume 挂载 `/app/data` 目录持久化 SQLite 数据库

---

## 3. 功能模块详细说明

### 3.1 日志管理 (Logs)
- **多数据源支持**：
    - **Loki**：支持 LogQL 查询
    - **Elasticsearch**：支持 Lucene 查询语法
- **结构化展示**：
    - 自动解析 Nginx/Access Log 格式
    - 字段提取：源地址、时间、请求地址、请求方式、状态码、耗时、后端地址、客户端信息
    - **状态指示**：根据 HTTP 状态码显示不同颜色 (2xx Green, 4xx Orange, 5xx Red)
- **高级筛选**：
    - **列筛选**：支持对源地址、时间、Host、Method、Status 等字段进行模糊搜索
    - **值选择**：提供列内唯一值的下拉选择列表
- **AI 分析**：
    - 集成 AI 模型对日志内容进行智能分析和诊断
    - 支持上下文对话

### 3.2 模型管理 (Models)
- **模型配置**：支持添加、编辑、删除 AI 模型配置
- **类型支持**：Ollama, OpenAI 兼容接口
- **状态管理**：启用/禁用模型，设置默认模型
- **连通性测试**：一键测试模型接口连通性

### 3.3 数据源管理 (Data Sources)
- **源配置**：支持添加 Loki 和 Elasticsearch 数据源
- **连接测试**：验证数据源地址和认证信息的有效性

### 3.4 系统管理
- **用户认证**：基于 JWT 的登录/登出机制
- **初始化**：支持通过环境变量 `AILAP_ADMIN_USER` 和 `AILAP_ADMIN_PASS` 初始化管理员账号

---

## 4. 接口规范 (API)

### 4.1 通用响应格式
```json
{
  "code": 0,          // 0 表示成功，非 0 表示错误
  "message": "success", // 提示信息
  "data": {}          // 业务数据
}
```

### 4.2 核心接口
- `POST /api/auth/login`: 用户登录
- `GET /api/logs/query`: 日志查询 (支持分页、时间范围、数据源选择)
- `POST /api/ai/analyze-logs`: AI 日志分析
- `GET /api/models`: 获取模型列表
- `GET /api/datasources`: 获取数据源列表

---

## 5. 开发与运行

### 5.1 本地开发
- **前端**：`npm run dev` (Port 5173/5174)
- **后端**：`go run ./cmd` (Port 8080)

### 5.2 Docker 运行
```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  --name ailap \
  ailap
```
