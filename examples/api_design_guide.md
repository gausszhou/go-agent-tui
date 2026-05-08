# HTTP API 设计指南

## RESTful 基础

### 资源命名

```
GET    /users              # 获取用户列表
GET    /users/123          # 获取单个用户
POST   /users              # 创建用户
PUT    /users/123          # 更新用户（完整）
PATCH  /users/123          # 更新用户（部分）
DELETE /users/123          # 删除用户
```

### HTTP 状态码

| 状态码 | 含义 | 使用场景 |
|--------|------|----------|
| 200 | OK | 成功获取/更新资源 |
| 201 | Created | 成功创建资源 |
| 204 | No Content | 成功删除资源 |
| 400 | Bad Request | 请求参数错误 |
| 401 | Unauthorized | 未认证 |
| 403 | Forbidden | 无权限 |
| 404 | Not Found | 资源不存在 |
| 500 | Internal Server Error | 服务器错误 |

## 请求格式

### 请求头

```http
Content-Type: application/json
Accept: application/json
Authorization: Bearer <token>
```

### 请求体示例

```json
{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30
}
```

## 响应格式

### 成功响应

```json
{
    "data": {
        "id": "123",
        "name": "John Doe",
        "email": "john@example.com"
    },
    "meta": {
        "timestamp": "2024-01-01T00:00:00Z"
    }
}
```

### 列表响应

```json
{
    "data": [
        {"id": "1", "name": "User 1"},
        {"id": "2", "name": "User 2"}
    ],
    "meta": {
        "total": 100,
        "page": 1,
        "per_page": 20
    }
}
```

### 错误响应

```json
{
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Invalid email format",
        "details": [
            {
                "field": "email",
                "message": "must be a valid email address"
            }
        ]
    }
}
```

## 分页

### 请求

```
GET /users?page=2&per_page=20
```

### 响应

```json
{
    "data": [...],
    "pagination": {
        "current_page": 2,
        "per_page": 20,
        "total": 100,
        "total_pages": 5
    }
}
```

## 排序和过滤

### 排序

```
GET /users?sort=name&order=asc
GET /users?sort=created_at&order=desc
```

### 过滤

```
GET /users?status=active&age_gte=18
GET /users?name_like=john&role=admin
```

## 认证

### Bearer Token

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### API Key

```
X-API-Key: your-api-key-here
```

## 版本控制

```
GET /v1/users
GET /v2/users
```

### 响应头版本

```http
Accept: application/vnd.myapp.v1+json
```

## 速率限制

### 响应头

```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1609459200
```

### 超出限制

```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json

{
    "error": {
        "code": "RATE_LIMIT_EXCEEDED",
        "message": "Too many requests",
        "retry_after": 60
    }
}
```

## 缓存

### ETag

```http
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"
```

### Last-Modified

```http
Last-Modified: Wed, 01 Jan 2024 00:00:00 GMT
```

### 条件请求

```http
If-None-Match: "33a64df551425fcc55e4d42a148795d9f25f89d4"
If-Modified-Since: Wed, 01 Jan 2024 00:00:00 GMT
```

## 超媒体

```json
{
    "data": {
        "id": "123",
        "name": "John"
    },
    "links": {
        "self": "/users/123",
        "posts": "/users/123/posts",
        "profile": "/profiles/456"
    }
}
```

## 示例代码

### Go HTTP 服务器

```go
type UserHandler struct {
    service UserService
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }

    users, err := h.service.ListUsers(r.Context(), page)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Response{Data: users})
}
```

### 前端调用

```javascript
async function fetchUsers(page = 1) {
    const response = await fetch(`/api/users?page=${page}`, {
        headers: {
            'Authorization': `Bearer ${token}`,
            'Accept': 'application/json'
        }
    });

    if (!response.ok) {
        throw new Error('Failed to fetch users');
    }

    return response.json();
}
```