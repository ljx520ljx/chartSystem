# Chart系统API设计文档

## 1. API概述

Chart系统采用RESTful风格设计API，主要提供以下功能接口：

- 用户认证与授权
- 文件管理
- 数据处理与分析
- 配置管理
- 实时数据传输（WebSocket）

## 2. 通用规范

### 2.1 基础URL结构

```
https://api.chartsystem.example.com/v1
```

所有API端点均使用`/v1`前缀，以便未来版本升级。

### 2.2 请求格式

- 所有请求应使用JSON格式（除文件上传外）
- 请求应包含适当的`Content-Type` header
- 身份验证通过`Authorization` header传递JWT令牌

### 2.3 响应格式

标准JSON响应结构：

```json
{
  "status": "success",       // success, error
  "data": {},                // 响应数据对象或数组
  "message": "",             // 错误消息（仅在状态为error时出现）
  "errors": [],              // 详细错误信息（仅在状态为error时出现）
  "pagination": {            // 分页信息（仅在适用时出现）
    "total": 100,
    "page": 1,
    "page_size": 20,
    "total_pages": 5
  },
  "request_id": "uuid-string" // 请求追踪ID
}
```

### 2.4 状态码

- 200 OK：请求成功
- 201 Created：资源创建成功
- 204 No Content：请求成功但无返回内容
- 400 Bad Request：请求参数错误
- 401 Unauthorized：未认证
- 403 Forbidden：权限不足
- 404 Not Found：资源不存在
- 422 Unprocessable Entity：参数验证失败
- 429 Too Many Requests：请求频率超限
- 500 Internal Server Error：服务器错误

### 2.5 版本控制

- API版本通过URL路径中的版本号标识（例如：`/v1/users`）
- 主要版本号变更表示不兼容的API变化

### 2.6 限流策略

- 基于IP地址和用户ID的限流
- 默认速率：每分钟60个请求
- 文件上传API：每分钟10个请求
- 使用HTTP头部`X-RateLimit-*`提供限流信息

## 3. 认证与授权API

### 3.1 用户注册

**端点：** `POST /v1/auth/register`

**请求体：**
```json
{
  "username": "testuser",
  "email": "user@example.com",
  "password": "securepassword",
  "full_name": "Test User"
}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "user_id": "uuid-string",
    "username": "testuser",
    "email": "user@example.com",
    "full_name": "Test User",
    "created_at": "2023-06-28T10:00:00Z"
  }
}
```

### 3.2 用户登录

**端点：** `POST /v1/auth/login`

**请求体：**
```json
{
  "username": "testuser",
  "password": "securepassword"
}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "access_token": "jwt-token-string",
    "refresh_token": "refresh-token-string",
    "expires_in": 3600,
    "token_type": "Bearer",
    "user": {
      "user_id": "uuid-string",
      "username": "testuser",
      "email": "user@example.com"
    }
  }
}
```

### 3.3 刷新令牌

**端点：** `POST /v1/auth/refresh`

**请求体：**
```json
{
  "refresh_token": "refresh-token-string"
}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "access_token": "new-jwt-token-string",
    "refresh_token": "new-refresh-token-string",
    "expires_in": 3600,
    "token_type": "Bearer"
  }
}
```

### 3.4 登出

**端点：** `POST /v1/auth/logout`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应：**
```json
{
  "status": "success",
  "data": null,
  "message": "已成功退出登录"
}
```

## 4. 用户管理API

### 4.1 获取用户信息

**端点：** `GET /v1/users/me`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "user_id": "uuid-string",
    "username": "testuser",
    "email": "user@example.com",
    "full_name": "Test User",
    "avatar_url": "https://example.com/avatar.jpg",
    "created_at": "2023-06-28T10:00:00Z",
    "roles": ["user", "analyst"]
  }
}
```

### 4.2 更新用户信息

**端点：** `PATCH /v1/users/me`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求体：**
```json
{
  "full_name": "Updated Name",
  "avatar_url": "https://example.com/new-avatar.jpg"
}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "user_id": "uuid-string",
    "username": "testuser",
    "email": "user@example.com",
    "full_name": "Updated Name",
    "avatar_url": "https://example.com/new-avatar.jpg",
    "updated_at": "2023-06-28T11:00:00Z"
  }
}
```

## 5. 文件管理API

### 5.1 文件列表

**端点：** `GET /v1/files`

**请求头：**
```
Authorization: Bearer {access_token}
```

**查询参数：**
- `page`: 页码（默认1）
- `page_size`: 每页大小（默认20）
- `sort`: 排序字段（默认created_at）
- `order`: 排序方向（asc或desc，默认desc）
- `search`: 搜索关键词
- `file_type`: 文件类型过滤

**响应：**
```json
{
  "status": "success",
  "data": [
    {
      "file_id": "uuid-string-1",
      "name": "ecg_record_001.bin",
      "description": "心电记录-患者A",
      "file_type": "binary/ecg",
      "file_size": 1048576,
      "created_at": "2023-06-28T10:00:00Z",
      "updated_at": "2023-06-28T10:00:00Z"
    },
    {
      "file_id": "uuid-string-2",
      "name": "pressure_record_002.bin",
      "description": "血压记录-患者B",
      "file_type": "binary/pressure",
      "file_size": 524288,
      "created_at": "2023-06-27T14:30:00Z",
      "updated_at": "2023-06-27T14:30:00Z"
    }
  ],
  "pagination": {
    "total": 42,
    "page": 1,
    "page_size": 20,
    "total_pages": 3
  }
}
```

### 5.2 上传文件

**端点：** `POST /v1/files`

**请求头：**
```
Authorization: Bearer {access_token}
Content-Type: multipart/form-data
```

**表单参数：**
- `file`: 文件数据
- `name`: 文件名称（可选，默认使用原始文件名）
- `description`: 文件描述（可选）
- `file_type`: 文件类型（可选）
- `metadata`: 文件元数据（JSON字符串，可选）

**响应：**
```json
{
  "status": "success",
  "data": {
    "file_id": "uuid-string",
    "name": "ecg_record_003.bin",
    "description": "心电记录-患者C",
    "file_type": "binary/ecg",
    "file_size": 2097152,
    "created_at": "2023-06-28T15:45:00Z",
    "upload_status": "completed"
  }
}
```

### 5.3 获取文件详情

**端点：** `GET /v1/files/{file_id}`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "file_id": "uuid-string",
    "name": "ecg_record_003.bin",
    "description": "心电记录-患者C",
    "file_type": "binary/ecg",
    "file_size": 2097152,
    "format": "binary",
    "metadata": {
      "patient_id": "P12345",
      "record_date": "2023-06-15",
      "device": "ECG-Monitor-5000"
    },
    "channels": [
      {
        "channel_id": "uuid-string-1",
        "name": "ECG-I",
        "type": "ecg",
        "sampling_rate": 500
      },
      {
        "channel_id": "uuid-string-2",
        "name": "ECG-II",
        "type": "ecg",
        "sampling_rate": 500
      }
    ],
    "created_at": "2023-06-28T15:45:00Z",
    "updated_at": "2023-06-28T15:45:00Z",
    "created_by": "uuid-string"
  }
}
```

### 5.4 更新文件信息

**端点：** `PATCH /v1/files/{file_id}`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求体：**
```json
{
  "name": "更新的文件名称",
  "description": "更新的文件描述",
  "metadata": {
    "patient_id": "P12345-UPDATED",
    "notes": "新增备注信息"
  }
}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "file_id": "uuid-string",
    "name": "更新的文件名称",
    "description": "更新的文件描述",
    "metadata": {
      "patient_id": "P12345-UPDATED",
      "record_date": "2023-06-15",
      "device": "ECG-Monitor-5000",
      "notes": "新增备注信息"
    },
    "updated_at": "2023-06-28T16:30:00Z"
  }
}
```

### 5.5 删除文件

**端点：** `DELETE /v1/files/{file_id}`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应：**
```json
{
  "status": "success",
  "data": null,
  "message": "文件已成功删除"
}
```

## 6. 数据通道API

### 6.1 获取文件通道列表

**端点：** `GET /v1/files/{file_id}/channels`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应：**
```json
{
  "status": "success",
  "data": [
    {
      "channel_id": "uuid-string-1",
      "name": "ECG-I",
      "display_name": "心电I导联",
      "type": "ecg",
      "unit": "mV",
      "sampling_rate": 500,
      "y_axis_min": -1.0,
      "y_axis_max": 1.0,
      "color": "#FF0000",
      "visible": true
    },
    {
      "channel_id": "uuid-string-2",
      "name": "ECG-II",
      "display_name": "心电II导联",
      "type": "ecg",
      "unit": "mV",
      "sampling_rate": 500,
      "y_axis_min": -1.0,
      "y_axis_max": 1.0,
      "color": "#00FF00",
      "visible": true
    }
  ]
}
```

### 6.2 更新通道配置

**端点：** `PATCH /v1/files/{file_id}/channels/{channel_id}`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求体：**
```json
{
  "display_name": "心电I导联(修改)",
  "color": "#990000",
  "visible": true,
  "y_axis_min": -2.0,
  "y_axis_max": 2.0
}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "channel_id": "uuid-string-1",
    "name": "ECG-I",
    "display_name": "心电I导联(修改)",
    "type": "ecg",
    "unit": "mV",
    "sampling_rate": 500,
    "y_axis_min": -2.0,
    "y_axis_max": 2.0,
    "color": "#990000",
    "visible": true,
    "updated_at": "2023-06-28T17:15:00Z"
  }
}
```

### 6.3 获取通道数据

**端点：** `GET /v1/files/{file_id}/channels/{channel_id}/data`

**请求头：**
```
Authorization: Bearer {access_token}
```

**查询参数：**
- `start_time`: 开始时间（秒，默认0）
- `end_time`: 结束时间（秒，可选）
- `max_points`: 最大返回点数（默认1000）
- `processing`: 处理方式（none, filter, smooth，默认none）

**响应：**
```json
{
  "status": "success",
  "data": {
    "channel_id": "uuid-string-1",
    "name": "ECG-I",
    "start_time": 0,
    "end_time": 10,
    "sampling_rate": 500,
    "point_count": 1000,
    "time_points": [0.000, 0.002, 0.004, ...],
    "values": [0.1, 0.12, 0.15, ...]
  }
}
```

## 7. 数据处理API

### 7.1 创建处理任务

**端点：** `POST /v1/files/{file_id}/processings`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求体：**
```json
{
  "process_type": "filter",
  "parameters": {
    "filter_type": "bandpass",
    "low_cutoff": 0.5,
    "high_cutoff": 40.0,
    "order": 4,
    "target_channels": ["uuid-string-1", "uuid-string-2"]
  },
  "save_result": true
}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "task_id": "uuid-string",
    "file_id": "uuid-string",
    "process_type": "filter",
    "parameters": {
      "filter_type": "bandpass",
      "low_cutoff": 0.5,
      "high_cutoff": 40.0,
      "order": 4,
      "target_channels": ["uuid-string-1", "uuid-string-2"]
    },
    "status": "queued",
    "created_at": "2023-06-28T17:30:00Z"
  }
}
```

### 7.2 获取处理任务状态

**端点：** `GET /v1/processings/{task_id}`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "task_id": "uuid-string",
    "file_id": "uuid-string",
    "process_type": "filter",
    "parameters": {
      "filter_type": "bandpass",
      "low_cutoff": 0.5,
      "high_cutoff": 40.0,
      "order": 4,
      "target_channels": ["uuid-string-1", "uuid-string-2"]
    },
    "status": "processing",
    "progress": 45,
    "message": "处理中...",
    "created_at": "2023-06-28T17:30:00Z",
    "updated_at": "2023-06-28T17:31:30Z"
  }
}
```

### 7.3 获取文件处理任务列表

**端点：** `GET /v1/files/{file_id}/processings`

**请求头：**
```
Authorization: Bearer {access_token}
```

**查询参数：**
- `page`: 页码（默认1）
- `page_size`: 每页大小（默认20）
- `status`: 任务状态过滤

**响应：**
```json
{
  "status": "success",
  "data": [
    {
      "task_id": "uuid-string-1",
      "process_type": "filter",
      "status": "completed",
      "created_at": "2023-06-28T17:30:00Z",
      "completed_at": "2023-06-28T17:35:00Z"
    },
    {
      "task_id": "uuid-string-2",
      "process_type": "fft",
      "status": "processing",
      "progress": 60,
      "created_at": "2023-06-28T17:40:00Z"
    }
  ],
  "pagination": {
    "total": 5,
    "page": 1,
    "page_size": 20,
    "total_pages": 1
  }
}
```

## 8. 标记API

### 8.1 创建标记

**端点：** `POST /v1/files/{file_id}/markers`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求体：**
```json
{
  "channel_id": "uuid-string-1", // 可选，null表示全局标记
  "marker_type": "event",
  "position": 15.5, // 时间点（秒）
  "label": "R波峰",
  "description": "明显的R波",
  "color": "#FF0000"
}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "marker_id": "uuid-string",
    "file_id": "uuid-string",
    "channel_id": "uuid-string-1",
    "marker_type": "event",
    "position": 15.5,
    "label": "R波峰",
    "description": "明显的R波",
    "color": "#FF0000",
    "created_by": "uuid-string",
    "created_at": "2023-06-28T17:45:00Z"
  }
}
```

### 8.2 获取文件标记

**端点：** `GET /v1/files/{file_id}/markers`

**请求头：**
```
Authorization: Bearer {access_token}
```

**查询参数：**
- `channel_id`: 按通道ID过滤（可选）
- `marker_type`: 按标记类型过滤（可选）
- `start`: 开始位置（可选）
- `end`: 结束位置（可选）

**响应：**
```json
{
  "status": "success",
  "data": [
    {
      "marker_id": "uuid-string-1",
      "file_id": "uuid-string",
      "channel_id": "uuid-string-1",
      "marker_type": "event",
      "position": 15.5,
      "label": "R波峰",
      "color": "#FF0000",
      "created_by": "uuid-string",
      "created_at": "2023-06-28T17:45:00Z"
    },
    {
      "marker_id": "uuid-string-2",
      "file_id": "uuid-string",
      "channel_id": null,
      "marker_type": "note",
      "position": 20.2,
      "label": "心律不齐",
      "color": "#00FF00",
      "created_by": "uuid-string",
      "created_at": "2023-06-28T17:46:30Z"
    }
  ]
}
```

### 8.3 删除标记

**端点：** `DELETE /v1/files/{file_id}/markers/{marker_id}`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应：**
```json
{
  "status": "success",
  "data": null,
  "message": "标记已成功删除"
}
```

## 9. 分析API

### 9.1 创建分析任务

**端点：** `POST /v1/files/{file_id}/analyses`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求体：**
```json
{
  "name": "心率变异性分析",
  "analysis_type": "hrv",
  "description": "标准HRV分析",
  "parameters": {
    "channel_id": "uuid-string-1",
    "start_time": 10,
    "end_time": 160,
    "detailed_metrics": true
  }
}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "analysis_id": "uuid-string",
    "file_id": "uuid-string",
    "name": "心率变异性分析",
    "analysis_type": "hrv",
    "status": "queued",
    "created_by": "uuid-string",
    "created_at": "2023-06-28T18:00:00Z"
  }
}
```

### 9.2 获取分析结果

**端点：** `GET /v1/files/{file_id}/analyses/{analysis_id}`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "analysis_id": "uuid-string",
    "file_id": "uuid-string",
    "name": "心率变异性分析",
    "analysis_type": "hrv",
    "description": "标准HRV分析",
    "parameters": {
      "channel_id": "uuid-string-1",
      "start_time": 10,
      "end_time": 160,
      "detailed_metrics": true
    },
    "status": "completed",
    "result": {
      "mean_rr": 825.2,
      "sdnn": 42.5,
      "rmssd": 38.7,
      "pnn50": 24.3,
      "lf_power": 1245.8,
      "hf_power": 857.3,
      "lf_hf_ratio": 1.45
    },
    "created_by": "uuid-string",
    "created_at": "2023-06-28T18:00:00Z",
    "completed_at": "2023-06-28T18:01:45Z"
  }
}
```

### 9.3 获取文件所有分析

**端点：** `GET /v1/files/{file_id}/analyses`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应：**
```json
{
  "status": "success",
  "data": [
    {
      "analysis_id": "uuid-string-1",
      "file_id": "uuid-string",
      "name": "心率变异性分析",
      "analysis_type": "hrv",
      "status": "completed",
      "created_at": "2023-06-28T18:00:00Z",
      "completed_at": "2023-06-28T18:01:45Z"
    },
    {
      "analysis_id": "uuid-string-2",
      "file_id": "uuid-string",
      "name": "频谱分析",
      "analysis_type": "spectrum",
      "status": "processing",
      "created_at": "2023-06-28T18:05:00Z"
    }
  ]
}
```

## 10. 导出API

### 10.1 导出数据

**端点：** `POST /v1/files/{file_id}/exports`

**请求头：**
```
Authorization: Bearer {access_token}
```

**请求体：**
```json
{
  "export_type": "csv",
  "channels": ["uuid-string-1", "uuid-string-2"],
  "start_time": 0,
  "end_time": 60,
  "include_markers": true,
  "include_processed_data": false
}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "export_id": "uuid-string",
    "file_id": "uuid-string",
    "export_type": "csv",
    "status": "processing",
    "created_at": "2023-06-28T18:30:00Z"
  }
}
```

### 10.2 获取导出状态

**端点：** `GET /v1/exports/{export_id}`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应：**
```json
{
  "status": "success",
  "data": {
    "export_id": "uuid-string",
    "file_id": "uuid-string",
    "export_type": "csv",
    "status": "completed",
    "download_url": "https://api.chartsystem.example.com/v1/exports/uuid-string/download",
    "expires_at": "2023-07-05T18:30:00Z",
    "created_at": "2023-06-28T18:30:00Z",
    "completed_at": "2023-06-28T18:32:15Z"
  }
}
```

### 10.3 下载导出文件

**端点：** `GET /v1/exports/{export_id}/download`

**请求头：**
```
Authorization: Bearer {access_token}
```

**响应：**
文件下载（二进制数据流）

## 11. WebSocket API

### 11.1 实时数据流

**端点：** `wss://api.chartsystem.example.com/v1/ws/files/{file_id}/stream`

**查询参数：**
- `token`: JWT访问令牌
- `channels`: 通道ID列表（逗号分隔）
- `rate`: 数据发送频率（赫兹）

**消息格式（服务器到客户端）：**
```json
{
  "type": "data",
  "timestamp": 1624825200000,
  "data": {
    "uuid-string-1": [0.1, 0.12, 0.15, ...],
    "uuid-string-2": [1.1, 1.05, 0.98, ...]
  }
}
```

**控制消息（客户端到服务器）：**
```json
{
  "action": "pause"  // pause, resume, rate_change
}
```

### 11.2 任务状态通知

**端点：** `wss://api.chartsystem.example.com/v1/ws/notifications`

**查询参数：**
- `token`: JWT访问令牌

**消息格式（服务器到客户端）：**
```json
{
  "type": "task_update",
  "task_id": "uuid-string",
  "task_type": "processing",
  "status": "completed",
  "progress": 100,
  "message": "处理已完成",
  "timestamp": 1624825260000
}
```

## 12. 系统信息API

### 12.1 获取支持的文件格式

**端点：** `GET /v1/system/file-formats`

**响应：**
```json
{
  "status": "success",
  "data": {
    "supported_formats": [
      {
        "format_id": "ecg_binary",
        "name": "心电二进制格式",
        "extension": ".bin",
        "mime_type": "application/octet-stream",
        "description": "标准心电二进制格式"
      },
      {
        "format_id": "edf",
        "name": "European Data Format",
        "extension": ".edf",
        "mime_type": "application/x-edf",
        "description": "欧洲数据格式，多通道生理信号记录"
      }
    ]
  }
}
```

### 12.2 获取支持的处理算法

**端点：** `GET /v1/system/processing-algorithms`

**响应：**
```json
{
  "status": "success",
  "data": {
    "algorithms": [
      {
        "algorithm_id": "filter",
        "name": "信号滤波",
        "description": "数字滤波处理",
        "parameters": [
          {
            "name": "filter_type",
            "type": "string",
            "enum": ["lowpass", "highpass", "bandpass", "notch"],
            "description": "滤波器类型"
          },
          {
            "name": "cutoff",
            "type": "number",
            "description": "截止频率",
            "required": true
          }
        ]
      },
      {
        "algorithm_id": "fft",
        "name": "快速傅里叶变换",
        "description": "FFT频域分析",
        "parameters": [
          {
            "name": "window_size",
            "type": "integer",
            "description": "窗口大小",
            "default": 1024
          }
        ]
      }
    ]
  }
}
``` 