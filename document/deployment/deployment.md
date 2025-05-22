# Chart系统部署指南

本文档提供Chart系统的部署方法和配置说明。

## Docker部署

### 前提条件

- 已安装Docker和Docker Compose
- 可访问互联网的环境（用于拉取镜像）
- 至少4GB RAM和10GB磁盘空间

### 使用Docker Compose部署

1. 确保项目根目录下存在`docker-compose.yml`文件

2. 在项目根目录执行以下命令启动所有服务:

```bash
docker-compose up -d
```

3. 验证服务是否正常启动:

```bash
docker-compose ps
```

4. 访问服务:
   - 前端界面: http://localhost:3000
   - 后端API: http://localhost:8080

### 使用单独Docker容器部署

1. 构建后端Docker镜像:

```bash
docker build -t chartsystem:latest .
```

2. 运行MySQL:

```bash
docker run -d --name mysql -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=chartsystem -p 3306:3306 mysql:8.0
```

3. 运行Redis:

```bash
docker run -d --name redis -p 6379:6379 redis:7.0
```

4. 运行后端服务:

```bash
docker run -d --name chartsystem -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=3306 \
  -e DB_USER=root \
  -e DB_PASSWORD=password \
  -e DB_NAME=chartsystem \
  -e REDIS_HOST=host.docker.internal \
  -e REDIS_PORT=6379 \
  chartsystem:latest
```

## 手动部署

### 后端部署

1. 确保已安装Go 1.24.1或更高版本
2. 安装依赖:

```bash
go mod tidy
```

3. 配置环境变量:

```bash
cp .env.example .env
```

4. 编辑.env文件设置数据库和Redis连接信息

5. 构建并运行:

```bash
go build -o chartsystem
./chartsystem
```

### 前端部署

1. 确保已安装Node.js 18+
2. 进入前端目录:

```bash
cd web
```

3. 安装依赖:

```bash
npm install
```

4. 构建生产版本:

```bash
npm run build
```

5. 使用Nginx或其他Web服务器部署dist目录

## 配置说明

系统的主要配置参数如下:

| 参数名 | 说明 | 默认值 |
|--------|------|--------|
| DB_HOST | 数据库主机地址 | localhost |
| DB_PORT | 数据库端口 | 3306 |
| DB_USER | 数据库用户名 | root |
| DB_PASSWORD | 数据库密码 | password |
| DB_NAME | 数据库名称 | chartsystem |
| REDIS_HOST | Redis主机地址 | localhost |
| REDIS_PORT | Redis端口 | 6379 |
| API_PORT | API服务端口 | 8080 |

## 常见问题

1. 数据库连接失败
   - 检查数据库服务是否正常运行
   - 验证连接参数是否正确
   - 确认网络连接是否畅通

2. Redis连接失败
   - 检查Redis服务是否正常运行
   - 验证连接参数是否正确
   - 检查防火墙设置

3. 前端无法连接后端API
   - 确认API服务是否正常运行
   - 检查CORS配置是否正确
   - 验证前端API基础URL配置 