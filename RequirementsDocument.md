# 📘 技术需求文档（细化版）

## 1. 技术架构要求

### 前端
- **框架**：Vue 3 + Vite  
- **UI 库**：Arco Design Vue（字节跳动 UI）  
- **状态管理**：Pinia  
- **路由**：Vue Router 4  
- **接口请求**：Axios（统一封装，支持 token 拦截、错误处理）  
- **打包部署**：前端单独构建产物（`dist/`），通过 Nginx 或容器挂载  

### 后端
- **语言**：Go  
- **Web 框架**：Gin  
- **ORM**：GORM（支持 PostgreSQL / MySQL）  
- **配置管理**：Viper  
- **日志**：Zap  
- **认证**：JWT  
- **模块划分**：
  - 日志服务模块（Loki / ELK 接入）  
  - 模型管理模块  
  - 数据源管理模块  
  - 用户 & 权限管理（可选）  

### 前后端分离
- 前端独立开发、打包 → Nginx / CDN 部署  
- 后端提供 **RESTful API** 接口  
- 前后端通过 HTTP(S) JSON API 交互  
- 后端支持 CORS  

---

## 2. 接口设计 (API 规范)

### 2.1 通用规范
- **协议**：HTTP/HTTPS  
- **数据格式**：JSON  
- **鉴权**：JWT Bearer Token  
- **返回格式**：
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
