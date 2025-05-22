# Chart系统数据库设计文档

## 1. 数据库概述

Chart系统采用混合型数据存储架构：
- **MySQL**：存储关系型数据，包括用户信息、文件元数据、配置信息等
- **MinIO/对象存储**：存储二进制文件原始数据
- **Redis**：缓存频繁访问数据、会话管理和任务队列

## 2. MySQL 数据库设计

### 2.1 ER图

```
                    +--------------+          +---------------+
                    |     User     |----------| Role          |
                    +--------------+          +---------------+
                    | PK: id       |          | PK: id        |
                    | username     |          | name          |
                    | email        |          | description   |
                    | password     |          +---------------+
                    | created_at   |                  ^
                    | updated_at   |                  |
                    +--------------+                  |
                          ^                  +---------------+
                          |                  | Permission    |
                          |                  +---------------+
                          |                  | PK: id        |
                          |                  | name          |
                          |                  | description   |
                          |                  +---------------+
+-----------------+       |
|     Analysis    |       |
+-----------------+       |
| PK: id          |       |
| file_id (FK)    |-------+--------> +-----------------+
| name            |       |          |      File       |
| description     |       |          +-----------------+
| analysis_type   |       |          | PK: id          |
| parameters      |       |          | user_id (FK)    |
| result          |       |          | name            |
| created_at      |       |          | path            |
| created_by (FK) |-------+          | size            |
| updated_at      |                  | format          |
+-----------------+                  | created_at      |
                                     | updated_at      |
                                     +-----------------+
                                         ^    ^    ^
                                         |    |    |
                          .--------------'    |    '-------------.
                          |                   |                  |
                +------------------+          |       +------------------+
                |     Marker       |          |       |  FileProcessing  |
                +------------------+          |       +------------------+
                | PK: id           |          |       | PK: id           |
                | file_id (FK)     |          |       | file_id (FK)     |
                | channel_id (FK)  |          |       | process_type     |
                | marker_type      |          |       | parameters_json  |
                | position         |          |       | status           |
                | label            |          |       | result_path      |
                | color            |          |       | created_at       |
                | created_at       |          |       | updated_at       |
                | created_by (FK)  |          |       +------------------+
                +------------------+          |
                                              |
                                    +------------------+
                                    |   DataChannel    |
                                    +------------------+
                                    | PK: id           |
                                    | file_id (FK)     |
                                    | name             |
                                    | type             |
                                    | config_json      |
                                    | processing_config|
                                    | created_at       |
                                    +------------------+
```

### 2.2 表结构定义

#### 2.2.1 用户表 (users)

| 字段名     | 类型        | 约束                 | 描述            |
|------------|-------------|---------------------|-------------------|
| id         | CHAR(36)    | PK                  | 用户唯一标识(UUID) |
| username   | VARCHAR(50) | UNIQUE, NOT NULL    | 用户名             |
| email      | VARCHAR(100)| UNIQUE, NOT NULL    | 电子邮箱           |
| password   | VARCHAR(100)| NOT NULL            | 加密密码           |
| full_name  | VARCHAR(100)| NULL                | 用户全名           |
| avatar_url | VARCHAR(255)| NULL                | 头像URL           |
| status     | TINYINT     | NOT NULL, DEFAULT 1 | 状态(1:活跃,0:禁用)|
| last_login | DATETIME    | NULL                | 最后登录时间       |
| created_at | DATETIME    | NOT NULL            | 创建时间           |
| updated_at | DATETIME    | NOT NULL            | 更新时间           |

索引：
- 主键索引：`id`
- 唯一索引：`username`, `email`
- 普通索引：`status`

#### 2.2.2 角色表 (roles)

| 字段名      | 类型        | 约束           | 描述            |
|------------|------------|---------------|----------------|
| id         | INT        | PK, AUTO_INCREMENT | 角色唯一标识    |
| name       | VARCHAR(50)| UNIQUE, NOT NULL | 角色名称         |
| description| VARCHAR(200)| NULL         | 角色描述         |
| created_at | DATETIME   | NOT NULL      | 创建时间         |
| updated_at | DATETIME   | NOT NULL      | 更新时间         |

#### 2.2.3 用户角色关联表 (user_roles)

| 字段名      | 类型        | 约束           | 描述            |
|------------|------------|---------------|----------------|
| user_id    | CHAR(36)   | PK, FK        | 用户ID          |
| role_id    | INT        | PK, FK        | 角色ID          |
| created_at | DATETIME   | NOT NULL      | 创建时间         |

复合主键：`(user_id, role_id)`

#### 2.2.4 权限表 (permissions)

| 字段名      | 类型        | 约束           | 描述            |
|------------|------------|---------------|----------------|
| id         | INT        | PK, AUTO_INCREMENT | 权限唯一标识    |
| name       | VARCHAR(50)| UNIQUE, NOT NULL | 权限名称         |
| code       | VARCHAR(50)| UNIQUE, NOT NULL | 权限编码         |
| description| VARCHAR(200)| NULL         | 权限描述         |
| created_at | DATETIME   | NOT NULL      | 创建时间         |
| updated_at | DATETIME   | NOT NULL      | 更新时间         |

#### 2.2.5 角色权限关联表 (role_permissions)

| 字段名       | 类型        | 约束           | 描述            |
|-------------|------------|---------------|----------------|
| role_id     | INT        | PK, FK        | 角色ID          |
| permission_id| INT       | PK, FK        | 权限ID          |
| created_at  | DATETIME   | NOT NULL      | 创建时间         |

复合主键：`(role_id, permission_id)`

#### 2.2.6 文件表 (files)

| 字段名      | 类型        | 约束           | 描述            |
|------------|------------|---------------|----------------|
| id         | CHAR(36)   | PK            | 文件唯一标识(UUID)|
| user_id    | CHAR(36)   | FK, NOT NULL  | 所有者用户ID     |
| name       | VARCHAR(255)| NOT NULL      | 文件名称         |
| description| TEXT       | NULL          | 文件描述         |
| file_path  | VARCHAR(500)| NOT NULL      | 存储路径         |
| file_type  | VARCHAR(50)| NOT NULL      | 文件类型         |
| file_size  | BIGINT     | NOT NULL      | 文件大小(字节)    |
| format     | VARCHAR(50)| NOT NULL      | 文件格式         |
| metadata   | JSON       | NULL          | 文件元数据(JSON)  |
| status     | TINYINT    | NOT NULL      | 状态(1:正常,0:删除)|
| created_at | DATETIME   | NOT NULL      | 创建时间         |
| updated_at | DATETIME   | NOT NULL      | 更新时间         |

索引：
- 主键索引：`id`
- 外键索引：`user_id`
- 普通索引：`file_type`, `status`

#### 2.2.7 数据通道表 (data_channels)

| 字段名          | 类型        | 约束           | 描述                |
|----------------|------------|---------------|---------------------|
| id             | CHAR(36)   | PK            | 通道唯一标识(UUID)    |
| file_id        | CHAR(36)   | FK, NOT NULL  | 所属文件ID           |
| name           | VARCHAR(100)| NOT NULL     | 通道名称             |
| display_name   | VARCHAR(100)| NULL         | 显示名称             |
| channel_type   | VARCHAR(50)| NOT NULL      | 通道类型             |
| unit           | VARCHAR(20)| NULL          | 单位                 |
| config         | JSON       | NOT NULL      | 通道配置(JSON)       |
| y_axis_min     | DOUBLE     | NULL          | Y轴最小值            |
| y_axis_max     | DOUBLE     | NULL          | Y轴最大值            |
| color          | VARCHAR(20)| NULL          | 显示颜色             |
| visible        | BOOLEAN    | DEFAULT TRUE  | 是否可见             |
| sampling_rate  | INT        | NOT NULL      | 采样率               |
| created_at     | DATETIME   | NOT NULL      | 创建时间             |
| updated_at     | DATETIME   | NOT NULL      | 更新时间             |

索引：
- 主键索引：`id`
- 外键索引：`file_id`
- 普通索引：`channel_type`

#### 2.2.8 文件处理表 (file_processings)

| 字段名           | 类型        | 约束           | 描述                |
|-----------------|------------|---------------|---------------------|
| id              | CHAR(36)   | PK            | 处理任务唯一标识(UUID)|
| file_id         | CHAR(36)   | FK, NOT NULL  | 文件ID              |
| process_type    | VARCHAR(50)| NOT NULL      | 处理类型             |
| parameters      | JSON       | NOT NULL      | 处理参数(JSON)       |
| status          | TINYINT    | NOT NULL      | 状态(0:等待,1:处理中,2:完成,3:失败)|
| result_file_path| VARCHAR(500)| NULL         | 结果文件路径         |
| message         | TEXT       | NULL          | 处理消息/错误信息     |
| progress        | TINYINT    | DEFAULT 0     | 处理进度(0-100)      |
| created_by      | CHAR(36)   | FK, NOT NULL  | 创建用户             |
| created_at      | DATETIME   | NOT NULL      | 创建时间             |
| updated_at      | DATETIME   | NOT NULL      | 更新时间             |
| completed_at    | DATETIME   | NULL          | 完成时间             |

索引：
- 主键索引：`id`
- 外键索引：`file_id`, `created_by`
- 普通索引：`status`, `process_type`

#### 2.2.9 标记表 (markers)

| 字段名        | 类型        | 约束           | 描述                |
|--------------|------------|---------------|---------------------|
| id           | CHAR(36)   | PK            | 标记唯一标识(UUID)    |
| file_id      | CHAR(36)   | FK, NOT NULL  | 文件ID              |
| channel_id   | CHAR(36)   | FK, NULL      | 通道ID(NULL表示全局标记)|
| marker_type  | VARCHAR(50)| NOT NULL      | 标记类型             |
| position     | DOUBLE     | NOT NULL      | 位置(时间点)         |
| label        | VARCHAR(255)| NULL         | 标签文本             |
| description  | TEXT       | NULL          | 描述                |
| color        | VARCHAR(20)| NULL          | 颜色                |
| created_by   | CHAR(36)   | FK, NOT NULL  | 创建用户             |
| created_at   | DATETIME   | NOT NULL      | 创建时间             |
| updated_at   | DATETIME   | NOT NULL      | 更新时间             |

索引：
- 主键索引：`id`
- 外键索引：`file_id`, `channel_id`, `created_by`
- 普通索引：`marker_type`, `position`

#### 2.2.10 分析结果表 (analyses)

| 字段名        | 类型        | 约束           | 描述                |
|--------------|------------|---------------|---------------------|
| id           | CHAR(36)   | PK            | 分析结果唯一标识(UUID)|
| file_id      | CHAR(36)   | FK, NOT NULL  | 文件ID              |
| name         | VARCHAR(100)| NOT NULL     | 分析名称             |
| description  | TEXT       | NULL          | 描述                |
| analysis_type| VARCHAR(50)| NOT NULL      | 分析类型             |
| parameters   | JSON       | NOT NULL      | 分析参数(JSON)       |
| result       | JSON       | NULL          | 分析结果(JSON)       |
| created_by   | CHAR(36)   | FK, NOT NULL  | 创建用户             |
| created_at   | DATETIME   | NOT NULL      | 创建时间             |
| updated_at   | DATETIME   | NOT NULL      | 更新时间             |

索引：
- 主键索引：`id`
- 外键索引：`file_id`, `created_by`
- 普通索引：`analysis_type`

## 3. 对象存储设计 (MinIO)

### 3.1 Bucket设计

| Bucket名称         | 用途                      | 访问控制         |
|-------------------|--------------------------|-----------------|
| files             | 存储用户上传的原始文件       | 私有            |
| processed-files   | 存储处理后的数据文件        | 私有            |
| exports           | 存储导出的报告和图表        | 私有，带过期链接  |
| temp              | 临时文件存储               | 私有，带生命周期策略|
| system-files      | 系统级配置和资源文件        | 私有            |

### 3.2 对象命名规范

- 原始文件：`files/{user-id}/{file-id}/{filename}`
- 处理文件：`processed-files/{user-id}/{file-id}/{process-type}/{timestamp}-{filename}`
- 导出文件：`exports/{user-id}/{file-id}/{export-type}/{timestamp}-{filename}`
- 临时文件：`temp/{user-id}/{request-id}/{filename}`

## 4. Redis设计

### 4.1 数据结构

| Key模式                         | 类型   | 用途                      | 过期时间 |
|--------------------------------|--------|--------------------------|----------|
| `session:{session-id}`         | Hash   | 存储会话信息              | 24小时   |
| `token:{token}`                | String | 用户令牌到用户ID的映射     | 与令牌一致 |
| `user:{user-id}:sessions`      | Set    | 用户活跃会话列表          | 无       |
| `file:{file-id}:metadata`      | Hash   | 文件元数据缓存            | 1小时    |
| `file:{file-id}:chunks:{idx}`  | String | 文件数据块缓存            | 30分钟   |
| `task:{task-id}`               | Hash   | 任务状态信息              | 24小时   |
| `rate-limit:{ip}:{endpoint}`   | String | API限流计数器             | 1分钟    |
| `notification:{user-id}`       | List   | 用户通知队列              | 30天     |

### 4.2 任务队列

使用Redis List实现任务队列：

| 队列名称                   | 用途                          |
|---------------------------|------------------------------|
| `queue:file-processing`   | 文件处理任务队列               |
| `queue:analysis`          | 数据分析任务队列               |
| `queue:export`            | 数据导出任务队列               |
| `queue:notification`      | 系统通知队列                  |

## 5. 数据库迁移与备份策略

### 5.1 迁移策略
- 使用golang-migrate工具管理数据库版本
- 所有架构变更通过迁移脚本实现
- CI/CD流程中包含自动化迁移步骤

### 5.2 备份策略
- MySQL：每日全量备份，每小时二进制日志备份
- MinIO：对象版本控制 + 每周全量备份
- Redis：RDB快照(每小时) + AOF日志

## 6. 性能优化策略

### 6.1 索引策略
- 所有外键字段创建索引
- 查询频繁的字段创建索引
- 使用前缀索引减少索引大小
- 定期分析查询性能，调整索引

### 6.2 分区策略
- 文件表按用户ID范围分区
- 大型分析结果表按时间范围分区
- 标记表按文件ID哈希分区

### 6.3 缓存策略
- 频繁访问的文件元数据缓存到Redis
- 活跃用户的文件列表缓存
- 常用配置项缓存
- 分析结果缓存 