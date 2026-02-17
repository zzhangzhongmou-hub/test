```markdown
# 红岩网校作业管理系统 (Redrock Assessment System)

一个面向红岩网校内部使用的作业管理系统，支持多部门协作，实现「老登」（管理员）发布作业、「小登」（学员）提交作业的完整流程。

## 功能清单

### 基础功能（必做）

**用户模块**
- 用户注册（支持用户名、密码、昵称、部门）
- 用户登录（双 Token 机制：Access + Refresh）
- 获取当前用户信息
- 刷新 Token
- 注销账号（软删除）
- 密码加密（bcrypt 加盐哈希）

**作业模块**
- 发布作业（设置截止时间和补交策略）
- 作业列表（按部门筛选、分页查询）
- 作业详情（关联发布者信息）
- 修改作业（同部门管理员可修改，乐观锁并发控制）
- 删除作业（同部门管理员可删除）

**提交模块**
- 提交作业（自动判断是否迟交）
- 我的提交列表
- 查看部门提交（老登查看本部门学员提交）
- 批改评语（分数+评语）
- 标记优秀作业
- 优秀作业展示

**系统功能**
- JWT 中间件认证
- 权限中间件（区分老登/小登）
- 统一响应格式（code + message + data）
- 部门标签自动转换（department_label）

### 进阶功能（选做）

**邮箱与通知**
- 绑定邮箱
- 作业发布邮件通知
- 截止前 24 小时提醒（定时任务）
- 评语通知邮件

**AI 作业评价**
- 调用 DeepSeek API 进行代码评价
- 生成改进建议和分数

**考核系统**
- 考核任务类型（区分作业和考核）
- 多人阅卷（支持多位老登阅卷）
- 乐观锁并发处理（版本号控制）
- 自动分配阅卷任务

**部署相关**
- Docker 部署支持
- docker-compose（MySQL + App）
- 交叉编译配置

## 技术栈

- 语言：Go 1.21+
- Web 框架：Gin
- ORM：GORM
- 数据库：MySQL 8.0
- 认证：JWT（双 Token 机制）
- 定时任务：robfig/cron
- 邮件服务：gomail（异步队列）
- AI 服务：DeepSeek API
- 部署：Docker & Docker Compose
- 配置管理：Viper

## 项目结构

homework-system/
├── cmd/                        # 程序入口
│   └── main.go
├── configs/                    # 配置文件
│   ├── config.go              # 配置结构定义
│   └── config-example.yaml    # 配置示例
├── cron/                       # 定时任务
│   └── reminder.go            # 截止提醒
├── dao/                        # 数据访问层（DAO）
│   ├── init.go                # 数据库初始化
│   ├── user.go
│   ├── homework.go
│   ├── submission.go
│   └── exam.go
├── handler/                    # HTTP 处理器（Handler）
│   ├── user.go
│   ├── homework.go
│   ├── submission.go
│   ├── exam.go
│   └── ai.go
├── middleware/                 # 中间件
│   ├── auth.go                # JWT 认证
│   ├── cors.go                # 跨域处理
│   └── logger.go              # 日志记录
├── models/                     # 数据模型（GORM）
│   ├── user.go
│   ├── homework.go
│   ├── submission.go
│   ├── exam.go
│   ├── examReview.go
│   └── constants.go           # 常量定义
├── pkg/                        # 公共工具包
│   ├── jwt/                   # JWT 工具
│   ├── hash/                  # 密码加密
│   ├── response/              # 统一响应
│   ├── email/                 # 邮件发送
│   └── ai/                    # AI 客户端
├── router/                     # 路由定义
│   └── router.go
├── service/                    # 业务逻辑层（Service）
│   ├── user.go
│   ├── homework.go
│   ├── submission.go
│   ├── exam.go
│   └── email_service.go
├── docker-compose.yml          # Docker 编排
├── go.mod
├── go.sum
└── README.md                   # 本文件

## 快速开始

### 方式一：Docker 部署（推荐）

确保已安装 Docker 和 Docker Compose。

1. 克隆仓库
```bash
git clone <your-repo-url>
cd homework-system
```

2. 创建配置文件
```bash
cp configs/config-example.yaml configs/config.yaml
# 编辑 config.yaml，填写你的数据库和 SMTP 配置
```

3. 启动服务
```bash
docker-compose up -d
```

服务将在 http://localhost:8080 启动，MySQL 运行在 localhost:3307（映射端口）。

4. 查看日志
```bash
docker-compose logs -f app
```

5. 停止服务
```bash
docker-compose down
# 如需删除数据卷
docker-compose down -v
```

### 方式二：本地运行

环境要求：
- Go 1.21+
- MySQL 8.0

1. 创建数据库
```sql
CREATE DATABASE homework_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 配置环境
```bash
cp configs/config-example.yaml configs/config.yaml
# 修改 configs/config.yaml 中的数据库连接信息
```

3. 安装依赖
```bash
go mod tidy
```

4. 运行项目
```bash
go run cmd/main.go
```

服务默认运行在 http://localhost:8080

## 配置文件说明

配置文件位于 configs/config.yaml，主要配置项：

```yaml
server:
  port: 8080              # 服务端口
  mode: debug             # 运行模式: debug/release

database:
  host: localhost         # 数据库地址
  port: 3306             # 数据库端口
  user: root             # 用户名
  password: "your-pass"  # 密码
  dbname: homework_system # 数据库名

jwt:
  secret: "your-secret"   # JWT 密钥（生产环境请更换）
  access_expire: 30       # Access Token 过期时间（分钟）
  refresh_expire: 10080   # Refresh Token 过期时间（分钟）

smtp:
  enable: true            # 是否启用邮件
  host: smtp.163.com      # SMTP 服务器
  port: 465               # 端口
  username: your@163.com  # 邮箱账号
  password: your-auth-code # 授权码（非密码）
  from: "红岩网校作业系统"  # 发件人名称

cron:
  enable: true            # 是否启用定时任务
  reminder_time: "0 9 * * *"  # 每天 9:00 检查截止提醒

ai:
  enable: true            # 是否启用 AI 评价
  provider: deepseek      # AI 提供商
  api_key: sk-xxx         # API 密钥
  base_url: https://api.deepseek.com/v1
  model: deepseek-coder   # 模型名称
  timeout: 10             # 超时时间（秒）
  max_tokens: 2000        # 最大 Token 数
```

## API 接口文档

- **在线文档（Apifox）**: (https://s.apifox.cn/99f01e0d-8396-4227-977d-c540d5ceeb07)

### 接口概览

**用户模块**
- POST /user/register - 注册
- POST /user/login - 登录
- POST /user/refresh - 刷新 Token
- GET /user/profile - 获取用户信息
- POST /user/bindEmail - 绑定邮箱
- DELETE /user/account - 注销账号

**作业模块**
- POST /homework - 发布作业（老登）
- GET /homework - 作业列表
- GET /homework/:id - 作业详情
- PUT /homework/:id - 修改作业（老登）
- DELETE /homework/:id - 删除作业（老登）

**提交模块**
- POST /submission - 提交作业（小登）
- GET /submission/my - 我的提交
- GET /submission/homework/:homework_id - 查看作业提交（老登）
- PUT /submission/:id/review - 批改作业（老登）
- PUT /submission/:id/excellent - 标记优秀（老登）
- GET /submission/excellent - 优秀作业列表

**考核模块**
- POST /exam - 创建考核（老登）
- GET /exam/reviews - 我的阅卷任务
- POST /exam/review - 提交阅卷

**AI 模块**
- POST /submission/:id/aiReview - AI 评价（老登）

### 认证说明

需要认证的接口请在请求头中添加：
```
Authorization: Bearer <access_token>
```

### 部门枚举值

| 枚举值 | 中文标签 |
|--------|----------|
| backend | 后端 |
| frontend | 前端 |
| sre | SRE |
| product | 产品 |
| design | 视觉设计 |
| android | Android |
| ios | iOS |

## 安全特性

- 密码安全：使用 bcrypt 加盐哈希存储
- JWT 双 Token：Access Token（短期）+ Refresh Token（长期）
- 乐观锁并发控制：使用版本号（version）防止并发修改冲突
- 软删除：用户注销使用软删除，保留数据完整性
- SQL 注入防护：使用 GORM 参数化查询

## 并发处理说明

### 作业修改并发控制
使用乐观锁（Optimistic Locking）机制：
- 每个作业记录包含 version 字段
- 修改时检查版本号，若已被他人修改则返回错误
- 确保同部门老登同时修改时的数据一致性

### 考核阅卷并发控制
- 每位考生的每份考卷可分配给多位老登
- 使用版本号控制防止重复批改
- 所有阅卷人完成阅卷后自动计算平均分

## 开发说明

### 交叉编译（Linux）
```bash
# 在 Windows/Mac 上编译 Linux 可执行文件
GOOS=linux GOARCH=amd64 go build -o homework-server cmd/main.go
```

### 数据库迁移
系统启动时会自动执行数据库迁移（AutoMigrate），无需手动执行 SQL。

## 常见问题

1. 邮件发送失败
   - 检查 SMTP 配置是否正确
   - 确认使用的是授权码而非邮箱密码（163/QQ 邮箱等）
   - 查看控制台日志中的错误信息

2. AI 评价超时
   - 默认超时时间为 8 秒
   - 如需调整，修改 configs/config.yaml 中的 ai.timeout

3. 数据库连接失败
   - Docker 部署时，确保 MySQL 完全启动后再启动 App（已配置 depends_on）
   - 本地部署时，检查 MySQL 服务是否运行

4. JWT 验证失败
   - 检查 jwt.secret 是否一致
   - 确认 Token 未过期

## Dockerfile 参考

在项目根目录创建 Dockerfile：

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o homework-server cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/homework-server .
COPY configs/config-example.yaml ./configs/config.yaml

EXPOSE 8080
CMD ["./homework-server"]
```

## 开源协议

MIT License

## 开发者

- 后端开发：Remindal
