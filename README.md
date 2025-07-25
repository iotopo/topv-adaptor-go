# TopV Adaptor Go 实现

这是一个基于 Go 的 HTTP 服务器实现，提供了实时数据查询、历史数据查询、设备标签查询、测点标签查询以及基于 NATS 的实时数据推送功能。

## 功能特性

1. **查询实时数据** - `GET /api/find_last`
2. **查询历史数据** - `POST /api/query_history`
3. **查询设备标签** - `GET /api/query_devices`
4. **查询测点标签** - `GET /api/query_points`
5. **设置数据值** - `POST /api/set_value`
6. **基于 NATS 推送实时数据** - 自动推送模拟数据到 NATS

## 依赖要求

- Go 1.19+
- NATS 服务器 (默认连接 `nats://127.0.0.1:4222`)

## 构建和运行

### 1. 构建项目

```bash
go mod tidy
go build -o topv-adaptor-go main.go push.go
```

### 2. 运行项目

```bash
# 直接运行
go run main.go push.go

# 或者运行编译后的二进制文件
./topv-adaptor-go
```

### 3. 启动 NATS 服务器

确保 NATS 服务器在 `127.0.0.1:4222` 运行：

```bash
# 使用 Docker
docker run -p 4222:4222 nats:latest

# 或者直接安装 NATS
nats-server
```

## API 接口

### 1. 查询实时数据

**请求：**
```http
GET /api/find_last
Content-Type: application/json

{
  "projectID": "project1",
  "tag": "group1.dev1.a",
  "device": false
}
```

**响应：**
```json
{
  "tag": "group1.dev1.a",
  "timestamp": "2024-01-01T12:00:00Z",
  "value": "12.3",
  "quality": 1
}
```

### 2. 查询历史数据

**请求：**
```http
POST /api/query_history
Content-Type: application/json

{
  "projectID": "project1",
  "tag": ["group1.dev1.a"],
  "start": "2024-01-01T00:00:00Z",
  "end": "2024-01-01T23:59:59Z"
}
```

**响应：**
```json
{
  "results": [
    {
      "tag": "group1.dev1.a",
      "values": [
        {
          "value": "12.3",
          "time": "2024-01-01T12:00:00Z"
        }
      ]
    }
  ]
}
```

### 3. 查询设备标签

**请求：**
```http
GET /api/query_devices
Content-Type: application/json

{
  "projectID": "project1"
}
```

**响应：**
```json
[
  {
    "tag": "group1",
    "name": "group1",
    "children": [
      {
        "tag": "group1.dev1",
        "name": "dev1",
        "isDevice": true
      }
    ],
    "isDevice": false
  }
]
```

### 4. 查询测点标签

**请求：**
```http
GET /api/query_points
Content-Type: application/json

{
  "projectID": "project1",
  "parentTag": "group1.dev1"
}
```

**响应：**
```json
[
  {
    "tag": "group1.dev1.a",
    "name": "a"
  },
  {
    "tag": "group1.dev1.b",
    "name": "b"
  },
  {
    "tag": "group1.dev1.c",
    "name": "c"
  }
]
```

### 5. 设置数据值

**请求：**
```http
POST /api/set_value
Content-Type: application/json

{
  "projectID": "project1",
  "tag": "group1.dev1.a",
  "value": "25.5",
  "time": 1640995200000
}
```

**响应：**
```json
{
  "code": "success",
  "msg": "Value set successfully"
}
```

## NATS 推送

服务启动后会自动开始推送实时数据到 NATS，主题格式为：`rtdb.{projectID}.{tag}`。

> 单项目版本， projectID 为固定值：`iotopo`

推送的数据格式：
```json
{
  "tag": "group1.dev1.a",
  "timestamp": "2024-01-01T12:00:00Z",
  "value": 45.67,
  "quality": 1
}
```

推送逻辑：
- 每秒推送一次数据
- 生成 3 个设备组，每组 10 个设备
- 每个设备推送一个测点数据
- 数值为 1-100 之间的随机数

## 配置

- HTTP 服务器端口：8080
- NATS 服务器地址：`nats://127.0.0.1:4222`
- 实时数据推送间隔：1秒

## 项目结构

```
topv-adaptor-go/
├── main.go          # 主程序，包含 HTTP 服务器和 API 处理
├── push.go          # NATS 推送服务
├── go.mod           # Go 模块定义
├── go.sum           # 依赖校验和
└── README.md        # 项目说明文档
```

## 数据模型

### ValueItem
```go
type ValueItem struct {
    Tag     string    `json:"tag"`
    Time    time.Time `json:"timestamp"`
    Value   any       `json:"value"`
    Quality int       `json:"quality"`
}
```

### DataItem
```go
type DataItem struct {
    Value any       `json:"value,omitempty"`
    Time  time.Time `json:"time,omitempty"`
}
```

### Result
```go
type Result struct {
    Tag    string     `json:"tag"`
    Values []DataItem `json:"values"`
}
```

### HistoryResponse
```go
type HistoryResponse struct {
    Results []Result `json:"results"`
    Msg     string   `json:"msg,omitempty"`
    Code    string   `json:"code,omitempty"`
}
```

### TagPoint
```go
type TagPoint struct {
    Tag  string `json:"tag,omitempty"`
    Name string `json:"name,omitempty"`
}
```

### Device
```go
type Device struct {
    ParentTag string    `json:"-"`
    Tag       string    `json:"tag,omitempty"`
    Name      string    `json:"name,omitempty"`
    Children  []*Device `json:"children,omitempty"`
    IsDevice  bool      `json:"isDevice,omitempty"`
}
```

## 开发说明

### 添加新的 API 端点

1. 在 `main.go` 中定义处理函数
2. 在 `main()` 函数中注册路由
3. 实现相应的业务逻辑

### 修改 NATS 推送

1. 编辑 `push.go` 中的 `realPush()` 函数
2. 修改推送的数据格式或频率
3. 重启服务以应用更改

## 测试

可以使用 curl 命令测试 API：

```bash
# 测试查询设备
curl -X GET http://localhost:8080/api/query_devices \
  -H "Content-Type: application/json" \
  -d '{"projectID":"project1"}'

# 测试查询实时数据
curl -X GET http://localhost:8080/api/find_last \
  -H "Content-Type: application/json" \
  -d '{"projectID":"project1","tag":"group1.dev1.a","device":false}'
```

## 部署

### Docker 部署

```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o topv-adaptor-go .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/topv-adaptor-go .
CMD ["./topv-adaptor-go"]
```

### 系统服务

创建 systemd 服务文件 `/etc/systemd/system/topv-adaptor.service`：

```ini
[Unit]
Description=TopV Adaptor Go Service
After=network.target

[Service]
Type=simple
User=topv
WorkingDirectory=/opt/topv-adaptor
ExecStart=/opt/topv-adaptor/topv-adaptor-go
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
``` 