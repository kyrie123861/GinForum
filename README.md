# GinForum

GinForum 是一个基于 Go 语言开发的社区论坛后端 API 服务，采用 Gin 框架构建，支持用户注册登录、社区分类、帖子发布与浏览、投票排序等核心功能。

## 功能特性

- **用户认证**：支持用户注册和登录，基于 JWT 的令牌认证机制
- **社区管理**：支持多社区分类，提供社区列表与详情查询
- **帖子管理**：支持创建帖子、查看帖子详情、分页获取帖子列表
- **投票系统**：支持赞成/反对投票，投票分数影响帖子排序（432 分 = 1 天排序权重）
- **双排序模式**：支持按创建时间排序和按投票分数热度排序
- **分布式 ID**：基于 Sonyflake 算法生成全局唯一、时间有序的帖子/用户 ID
- **优雅关机**：支持信号捕获，实现平滑的服务关闭
- **pprof 性能分析**：内置性能剖析接口，便于性能调优
- **配置热更新**：基于 Viper 实现配置文件变更的自动重载

## 技术栈

| 类别            | 技术                    | 版本  |
| --------------- | ----------------------- | ----- |
| **语言**        | Go                      | 1.25  |
| **HTTP 框架**   | Gin                     | v1.11 |
| **数据库**      | MySQL (sqlx)            | -     |
| **缓存 / 排序** | Redis (go-redis)        | v9    |
| **配置管理**    | Viper                   | v1.21 |
| **日志**        | Zap + Lumberjack        | v1.27 |
| **认证**        | JWT (dgrijalva/jwt-go)  | v3.2  |
| **分布式 ID**   | Sonyflake               | v1.3  |
| **参数校验**    | go-playground/validator | v10   |
| **限流**        | golang.org/x/time/rate  | v0.14 |
| **性能分析**    | gin-contrib/pprof       | v1.5  |

## 项目架构

项目采用分层架构设计，层次清晰，职责分明：

```
┌─────────────────────────────────────┐
│         HTTP Layer (Gin)            │
│   routers.go — 路由注册与中间件       │
├─────────────────────────────────────┤
│      Controllers (控制器层)          │
│   参数校验、请求解析、响应封装         │
├─────────────────────────────────────┤
│         Logic (业务逻辑层)           │
│   核心业务规则、流程编排              │
├─────────────────────────────────────┤
│       DAO (数据访问层)              │
│   ├── dao/mysql — MySQL 数据库操作   │
│   └── dao/redis — Redis 缓存操作     │
└─────────────────────────────────────┘
```

### 中间件

| 中间件      | 文件                       | 说明                                             |
| ----------- | -------------------------- | ------------------------------------------------ |
| GinLogger   | `logger/logger.go`         | 自定义请求日志，记录方法、路径、耗时、状态码等   |
| GinRecovery | `logger/logger.go`         | Panic 恢复，支持断管检测与堆栈打印               |
| JWTAuth     | `middlewares/auth.go`      | JWT 令牌验证，解析 Bearer Token 并注入用户上下文 |
| RateLimit   | `middlewares/ratelimit.go` | 令牌桶限流（已实现，未默认启用）                 |

## 项目结构

```
bluebell/
├── main.go                    # 应用入口：初始化各组件，启动 HTTP 服务
├── go.mod / go.sum            # Go 模块定义与依赖校验
├── Makefile                   # 构建、运行、格式化等命令
├── conf/
│   └── config.yaml            # 应用配置文件（端口、MySQL、Redis、日志、JWT）
├── settings/
│   └── srttings.go            # 配置结构体定义 + Viper 初始化与热更新
├── logger/
│   └── logger.go              # Zap 日志初始化、Gin 日志/恢复中间件
├── router/
│   └── routers.go             # Gin 路由注册（全部 API 路由）
├── middlewares/
│   ├── auth.go                # JWT 认证中间件
│   └── ratelimit.go           # 令牌桶限流中间件
├── controllers/
│   ├── code.go                # 业务错误码枚举（1000-1008）
│   ├── response.go            # 统一 JSON 响应格式
│   ├── request.go             # 从 Gin 上下文获取当前用户、分页参数解析
│   ├── validator.go           # Gin validator 翻译器初始化（中文/英文）
│   ├── user.go                # SignUpHandler / LoginHandler
│   ├── post.go                # 帖子相关 Handler
│   ├── community.go           # 社区相关 Handler
│   └── vote.go                # 投票 Handler
├── logic/
│   ├── user.go                # 注册 / 登录业务逻辑
│   ├── post.go                # 帖子创建、查询、Redis 初始化逻辑
│   ├── community.go           # 社区查询逻辑
│   └── vote.go                # 投票逻辑（含 6 种投票状态转换）
├── dao/
│   ├── mysql/
│   │   ├── mysql.go           # MySQL 连接初始化 (sqlx)
│   │   ├── user.go            # 用户 CRUD（MD5 密码加密）
│   │   ├── post.go            # 帖子 CRUD
│   │   └── communtiy.go       # 社区 CRUD
│   └── redis/
│       ├── redis.go           # Redis 客户端初始化 (go-redis/v9)
│       ├── keys.go            # Redis Key 命名空间常量
│       ├── post.go            # 帖子排序集合操作（ZSet）、Pipeline 批处理
│       └── vote.go            # 投票分数计算与记录
├── models/
│   ├── user.go                # User 数据模型
│   ├── post.go                # Post / ApiPostDetail (复合 DTO) 模型
│   ├── community.go           # Community / CommunityDetail 模型
│   └── params.go              # 请求参数结构体（注册、登录、投票、列表查询）
└── pkg/
    ├── jwt/
    │   └── jwt.go             # JWT 生成与解析 (HS256)
    └── snowflake/
        └── snowflake.go       # Sonyflake 分布式 ID 生成器
```

## 快速开始

### 环境要求

- Go 1.25+
- MySQL 8.0+
- Redis 6.0+

### 1. 初始化数据库

```sql
-- 创建数据库
CREATE DATABASE bluebell DEFAULT CHARACTER SET utf8mb4;

-- 用户表
CREATE TABLE `user` (
    `user_id`   BIGINT NOT NULL,
    `username`  VARCHAR(64) NOT NULL,
    `password`  VARCHAR(64) NOT NULL,
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 社区表
CREATE TABLE `community` (
    `community_id`   INT NOT NULL AUTO_INCREMENT,
    `community_name` VARCHAR(128) NOT NULL,
    `introduction`   VARCHAR(256) NOT NULL DEFAULT '',
    `create_time`    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`community_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 帖子表
CREATE TABLE `post` (
    `post_id`      BIGINT NOT NULL,
    `author_id`    BIGINT NOT NULL,
    `community_id` INT NOT NULL,
    `status`       TINYINT NOT NULL DEFAULT 1,
    `title`        VARCHAR(256) NOT NULL,
    `content`      TEXT NOT NULL,
    `create_time`  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`post_id`),
    KEY `idx_community_id` (`community_id`),
    KEY `idx_author_id` (`author_id`),
    KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 2. 配置文件

修改 `conf/config.yaml` 中的数据库连接信息：

```yaml
mode: "dev"
port: 8080

mysql:
  host: 127.0.0.1
  port: 3306
  user: root
  password: your_password
  db: bluebell
  max_open_conns: 100
  max_idle_conns: 20

redis:
  host: "127.0.0.1"
  port: 6379
  password: ""
  db: 0
  PoolSize: 100

auth:
  jwt_expire: 8760  # Token 过期时间（小时）

log:
  level: "debug"
  filename: "./bluebell.log"
  max_size: 1000
  max_age: 3600
  max_backups: 5
```

### 3. 运行

```bash
# 安装依赖
go mod tidy

# 直接运行
make run

# 或编译运行
make build
./bluebell
```

服务启动后将监听 `http://localhost:8080`。

## API 接口文档

所有接口前缀：`/api/v1`

### 公开接口（无需认证）

| 方法 | 路径      | 说明                       |
| ---- | --------- | -------------------------- |
| POST | `/signup` | 用户注册                   |
| POST | `/login`  | 用户登录（返回 JWT Token） |

### 受保护接口（需要 JWT 认证）

> 请求头：`Authorization: Bearer <token>`

| 方法 | 路径                | 说明                                                     |
| ---- | ------------------- | -------------------------------------------------------- |
| GET  | `/community`        | 获取社区列表                                             |
| GET  | `/community/:id`    | 获取社区详情                                             |
| POST | `/post`             | 创建帖子                                                 |
| GET  | `/post/:id`         | 获取帖子详情（含作者、社区、投票数）                     |
| GET  | `/posts/`           | 获取帖子列表（按创建时间，MySQL 直查分页）               |
| GET  | `/posts2`           | 获取帖子列表（按时间/分数排序，Redis 排序 + MySQL 补全） |
| POST | `/vote`             | 投票（赞成 1 / 反对 -1 / 取消 0）                        |
| GET  | `/init/redis/posts` | 初始化帖子到 Redis（测试用）                             |

### 请求示例

**注册**

```json
POST /api/v1/signup
{
    "username": "alice",
    "password": "123456",
    "re_password": "123456"
}
```

**登录**

```json
POST /api/v1/login
{
    "username": "alice",
    "password": "123456"
}
```

**创建帖子**

```json
POST /api/v1/post
Authorization: Bearer <token>
{
    "title": "Hello Bluebell",
    "content": "这是我的第一篇帖子",
    "community_id": 1
}
```

**获取帖子列表（按热度）**

```
GET /api/v1/posts2?page=1&size=10&order=score
```

**投票**

```json
POST /api/v1/vote
Authorization: Bearer <token>
{
    "post_id": 1234567890,
    "direction": 1
}
```

### 统一响应格式

```json
{
    "code": 1000,
    "msg": "success",
    "data": {}
}
```

### 错误码

| 错误码 | 说明             |
| ------ | ---------------- |
| 1000   | 成功             |
| 1001   | 请求参数错误     |
| 1002   | 用户名已存在     |
| 1003   | 用户名不存在     |
| 1004   | 用户名或密码错误 |
| 1005   | 服务繁忙         |
| 1006   | 无效的 Token     |
| 1007   | 认证格式有误     |
| 1008   | 未登录           |

## 核心设计

### 投票排序算法

- 每投一票影响 **432** 分（86400 / 200）
- **200 张赞成票** = 帖子获得额外 **1 天** 的热度曝光
- 帖子创建时间作为初始分数，投票后分数动态变化
- 投票有效期：**帖子发布后 7 天内**可投票，超期锁定

### Redis 数据结构

| Key                            | 类型 | 说明                       |
| ------------------------------ | ---- | -------------------------- |
| `bluebell:post:time`           | ZSet | 帖子 ID → 创建时间戳       |
| `bluebell:post:score`          | ZSet | 帖子 ID → 投票分数         |
| `bluebell:post:voted:{postID}` | ZSet | 用户 ID → 投票类型（1/-1） |

### 双层帖子列表

- `/posts/`：直接 MySQL `ORDER BY create_time DESC` 分页查询，简单直接
- `/posts2`：先从 Redis ZSet 按分数/时间排序获取帖子 ID 列表，再用 `WHERE post_id IN (...)` 批量查询 MySQL，使用 `FIND_IN_SET` 保持 Redis 排序顺序，适合高并发场景

## Makefile 命令

```bash
make          # 格式化代码 + 编译
make build    # 编译（交叉编译 Linux amd64）
make run      # 直接运行
make clean    # 清除编译产物
make gotool   # 运行 go fmt 和 go vet
make help     # 查看帮助
```
