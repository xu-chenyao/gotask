# 个人博客系统后端 API

基于 Go + Gin + GORM 开发的个人博客系统后端，实现了用户认证、文章管理和评论功能。

## 功能特性

- ✅ 用户注册和登录
- ✅ JWT 身份认证
- ✅ 文章 CRUD 操作
- ✅ 评论功能
- ✅ 权限控制（只有作者可以修改/删除文章）
- ✅ 分页查询
- ✅ 统一错误处理
- ✅ 日志记录

## 技术栈

- **Go 1.23.0** - 编程语言
- **Gin** - Web 框架
- **GORM** - ORM 库
- **MySQL** - 数据库
- **JWT** - 身份认证
- **Bcrypt** - 密码加密
- **Logrus** - 日志库

## 项目结构

```
task4/
├── config/          # 配置文件
│   └── database.go  # 数据库配置
├── controllers/     # 控制器
│   ├── auth.go      # 认证控制器
│   ├── post.go      # 文章控制器
│   └── comment.go   # 评论控制器
├── middleware/      # 中间件
│   └── auth.go      # JWT 认证中间件
├── models/          # 数据模型
│   ├── user.go      # 用户模型
│   ├── post.go      # 文章模型
│   └── comment.go   # 评论模型
├── routes/          # 路由配置
│   └── routes.go    # 路由定义
├── utils/           # 工具函数
│   └── jwt.go       # JWT 工具
├── main.go          # 主程序入口
├── go.mod           # Go 模块文件
└── README.md        # 项目说明
```

## 数据库设计

### users 表
- `id` - 主键
- `username` - 用户名（唯一）
- `password` - 密码（加密）
- `email` - 邮箱（唯一）
- `created_at` - 创建时间
- `updated_at` - 更新时间

### posts 表
- `id` - 主键
- `title` - 文章标题
- `content` - 文章内容
- `user_id` - 用户ID（外键）
- `created_at` - 创建时间
- `updated_at` - 更新时间

### comments 表
- `id` - 主键
- `content` - 评论内容
- `user_id` - 用户ID（外键）
- `post_id` - 文章ID（外键）
- `created_at` - 创建时间

## 快速开始

### 1. 环境要求

- Go 1.23.0+
- MySQL 8.0+
- Git

### 2. 安装依赖

```bash
go mod init task4
go mod tidy
docker exec mysql-local mysql -uroot -p123456 -e "CREATE DATABASE IF NOT EXISTS blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### 3. 数据库配置

确保 MySQL 服务正在运行，并创建数据库：

```sql
CREATE DATABASE blog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

修改 `config/database.go` 中的数据库连接信息：

```go
const (
    DBUser     = "root"        // 数据库用户名
    DBPassword = "123456"      // 数据库密码
    DBHost     = "localhost"   // 数据库主机
    DBPort     = "3306"        // 数据库端口
    DBName     = "blog_db"     // 数据库名称
)
```

### 4. 启动服务

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

## API 接口文档

### 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: Bearer Token (JWT)

### 认证接口

#### 用户注册
```http
POST /api/v1/auth/register
Content-Type: application/json

{
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com"
}
```

#### 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
    "username": "testuser",
    "password": "password123"
}
```

### 文章接口

#### 获取文章列表
```http
GET /api/v1/posts?page=1&page_size=10
```

#### 获取单个文章
```http
GET /api/v1/posts/{id}
```

#### 创建文章（需要认证）
```http
POST /api/v1/posts
Authorization: Bearer {token}
Content-Type: application/json

{
    "title": "文章标题",
    "content": "文章内容"
}
```

#### 更新文章（需要认证，只有作者可操作）
```http
PUT /api/v1/posts/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
    "title": "更新后的标题",
    "content": "更新后的内容"
}
```

#### 删除文章（需要认证，只有作者可操作）
```http
DELETE /api/v1/posts/{id}
Authorization: Bearer {token}
```

### 评论接口

#### 获取文章评论
```http
GET /api/v1/comments/post/{post_id}?page=1&page_size=20
```

#### 创建评论（需要认证）
```http
POST /api/v1/comments
Authorization: Bearer {token}
Content-Type: application/json

{
    "content": "评论内容",
    "post_id": 1
}
```

### 健康检查

```http
GET /health
```

## 测试用例

### 1. 用户注册测试

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com"
  }'
```

### 2. 用户登录测试

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 3. 创建文章测试

```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "我的第一篇文章",
    "content": "这是文章内容"
  }'
```

### 4. 获取文章列表测试

```bash
curl -X GET http://localhost:8080/api/v1/posts
```

## 部署说明

### 生产环境配置

1. 修改 JWT 密钥（在 `utils/jwt.go` 中）
2. 配置生产环境数据库
3. 设置环境变量
4. 使用反向代理（如 Nginx）



## 开发说明
### 添加新功能

1. 在 `models/` 中定义数据模型
2. 在 `controllers/` 中实现业务逻辑
3. 在 `routes/` 中配置路由
4. 更新数据库迁移

