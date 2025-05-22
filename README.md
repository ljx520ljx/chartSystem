# Chart系统

## 项目简介

Chart系统是一个前后端分离的医疗数据可视化平台，专注于心电图、血压等医疗监测数据的处理、分析与可视化。系统支持大型二进制数据文件解析、多通道数据展示、信号处理与分析功能，为医疗数据分析提供强大的工具支持。

## 系统架构

系统采用前后端分离架构：

- **前端**：React + TypeScript，使用D3.js和ECharts实现高性能数据可视化
- **后端**：Go 1.24.1 + Gin，模块化设计，提供RESTful API和WebSocket接口
- **数据库**：MySQL存储结构化数据，MinIO存储二进制文件，Redis用于缓存
- **部署**：Docker + Kubernetes容器化部署

![系统架构图](document/deployment/architecture.png)

## 主要功能

- 用户认证与权限管理
- 文件上传与管理
- 多通道医疗数据可视化
- 实时数据浏览与交互（缩放、平移、标记）
- 信号处理与分析（滤波、特征提取、异常检测）
- 数据导出与报告生成

## 快速开始

### 环境要求

- Go 1.24.1+
- Node.js 18+
- Docker & Docker Compose
- MySQL 8.0+
- Redis 7.0+
- MinIO (或兼容S3的对象存储)

### 后端启动

```bash
# 克隆仓库
git clone https://github.com/ljx520ljx/chartSystem.git
cd chartsystem

# 安装依赖
go mod tidy

# 配置环境变量
cp .env.example .env
# 根据需要修改.env文件中的配置

# 启动服务
go run main.go
```

### 前端启动

```bash
cd web

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

### Docker部署

```bash
# 在项目根目录下构建镜像
docker build -t chartsystem:latest .

# 使用docker-compose启动所有服务
docker-compose up -d
```

## 项目结构

```
chartsystem/
├── api/             # API路由和控制器
├── config/          # 配置管理
├── internal/        # 内部包
│   ├── middleware/  # 中间件
│   ├── model/       # 数据模型
│   ├── repository/  # 数据访问层
│   ├── service/     # 业务逻辑层
│   └── utils/       # 工具函数
├── pkg/             # 可复用包
├── web/             # 前端React项目
│   ├── public/      # 静态资源
│   ├── src/         # 源代码
│   └── package.json # 项目配置
├── document/        # 文档
│   ├── api_design.md        # API设计文档
│   ├── database_design.md   # 数据库设计文档
│   ├── design.md            # 系统设计文档
│   └── deployment/          # 部署相关文档和图片
├── main.go          # 程序入口
├── Dockerfile       # Docker构建文件
└── docker-compose.yml # Docker编排配置
```

## API文档

API文档详情请参考[API设计文档](document/api_design.md)。

## 部署

### 使用Docker部署

1. 构建Docker镜像
   ```bash
   docker build -t chartsystem:latest .
   ```

2. 运行容器
   ```bash
   docker run -p 8080:8080 -e DB_HOST=mysql -e REDIS_HOST=redis chartsystem:latest
   ```

3. 使用docker-compose (推荐)
   ```bash
   docker-compose up -d
   ```

详细的部署指南请参考[部署文档](document/deployment/deployment.md)。

## 许可证

本项目采用 MIT 许可证 - 详情请参见 [LICENSE](LICENSE) 文件 