# CLAUDE.md

此文件为 Claude Code (claude.ai/code) 在此代码仓库中工作时提供指导。

## Project Overview

这是一个基于 Go 1.24.0 + Gin + MySQL + Redis 技术栈构建的在线商城后端系统。采用清晰的分层架构模式，各层职责明确分离。

## Commands

### Building and Running
```bash
cd backend

# 安装/更新依赖
go mod tidy

# 构建项目
go build ./...

# 运行服务器
go run cmd/server/main.go

# 构建二进制文件
go build -o online-mall cmd/server/main.go

# 运行测试
go test ./...

# 运行性能测试
go test -bench=. -benchmem
```

### Database
```bash
# 初始化数据库（执行 SQL 脚本）
mysql -u root -p < scripts/init.sql
```

## Architecture

项目采用**分层架构**模式，包含以下关键组件：

### Layer Responsibilities

**Models 层** (`internal/models/`)
- 定义数据库表结构
- 包含 `BaseModel` 基础字段（ID、CreatedAt、UpdatedAt、DeletedAt）
- 使用 GORM 进行 ORM 操作
- 通过 `models.DB` 使用全局 `DB` 变量访问数据库
- 启动时自动迁移（见 `database.go:autoMigrate()`）

**Repository 层** (`internal/repository/`)
- 数据访问层，封装数据库操作
- 例如：`ProductRepository` 包含 `GetByID`、`GetProducts`、`Create`、`Update`、`Delete` 等方法
- 每个 repository 都有构造函数：`NewXxxRepository()`

**Service 层** (`internal/service/`)
- 业务逻辑层
- 协调 repository 操作并实现业务规则
- 例如：`ProductService.GetProducts()` 调用 `ProductRepository.GetProducts()`
- 有构造函数：`NewXxxService()`

**Controller 层** (`internal/api/controller/`)
- API 端点的 HTTP 处理器
- 使用 Gin context 处理请求和响应
- 调用 service 层执行业务操作
- **重要**：响应函数必须使用 `utils.` 前缀（例如：`utils.Success()`、`utils.ParamError()`、`utils.BadRequest()`、`utils.NotFound()`、`utils.ServerError()`）

**Routes** (`internal/api/routes/`)
- 定义 API 路由和中间件链
- 使用 `middleware.JWTAuth()` 保护需要认证的路由
- 使用 `middleware.RequireAdmin()` 限制仅管理员可访问的路由

**Middleware** (`internal/api/middleware/`)
- JWT 认证中间件
- CORS 处理
- 日志记录
- 请求验证

### Global State

- **Database**: `models.DB` - 全局 GORM DB 实例
- **Config**: `config.GlobalConfig` - 启动时加载的全局配置
- **Redis**: `utils.RedisClient` - 全局 Redis 客户端

### Configuration

配置使用 Viper 从 `configs/config.yaml` 加载。结构定义在 `internal/config/config.go` 中。主要配置部分：
- `app`: 名称、版本、调试模式、端口（默认 8080）
- `database`: MySQL 连接详情
- `redis`: Redis 连接详情
- `jwt`: JWT 密钥、过期时间、签发者

### Data Flow

```
HTTP 请求
  → 中间件（CORS、JWT、Logger）
  → Controller（请求解析、响应格式化）
  → Service（业务逻辑）
  → Repository（数据访问）
  → models.DB（GORM 操作）
  → MySQL
```

### Important Implementation Notes

1. **Model ID 类型**：所有 ID 都是 `uint64` 类型（定义在 `BaseModel` 中）
2. **JWT 中间件**：在 context 中设置 `user_id`（uint）、`username`（string）和 `role`（string）
3. **类型转换**：当从 Gin context 使用 `c.GetUint("user_id")` 时，需要转换为 `uint64` 用于 model 操作
4. **软删除**：BaseModel 包含 `DeletedAt` 用于 GORM 软删除支持
5. **JSON 字段**：某些字段存储 JSON 字符串（如 Product.Images），提供 getter/setter 方法
6. **Repository 模式**：每个实体都有自己的 repository - 不要从 service/controller 层直接访问 `models.DB`

### Current Module Status

**已完成**：
- 用户认证（登录、注册、token 刷新、登出）
- 用户资料管理
- 商品模块（Repository、Service、Controller）
- 分类模块（Repository、Service、Controller）

**待开发**（在 `routes.go` 中已注释）：
- 购物车模块
- 订单模块
- 优惠券模块
- 地址管理模块
- 搜索功能
- 文件上传功能

### Development Pattern

添加新模块（如 Cart）时，遵循以下模式：

1. **Models**：在 `internal/models/` 中定义结构体（如果不存在）
2. **Repository**：创建 `internal/repository/cart_repository.go`，包含数据访问方法
3. **Service**：创建 `internal/service/cart_service.go`，包含业务逻辑
4. **Controller**：创建 `internal/api/controller/cart.go`，包含 HTTP 处理器
5. **Routes**：在 `internal/api/routes/routes.go` 中添加路由
6. **Import utils**：在 controller 中始终导入 `"online-mall/internal/utils"`
7. **Response helpers**：使用 `utils.Success()`、`utils.ParamError()`、`utils.NotFound()` 等

### Database Auto-Migration

启动时，`models.InitDatabase()` 调用 `autoMigrate()`，为所有 model 运行 GORM AutoMigrate。这会自动创建/更新表结构。不使用手动 SQL 迁移文件来更改表结构。

### Redis Usage

直接使用 `utils.RedisClient` 或 `utils/redis.go` 中的辅助函数：
- `Set()`、`Get()`、`Del()`、`Exists()`
- `HSet()`、`HGet()`、`HGetAll()`、`HDel()`
- `LPush()`、`RPush()`、`LPop()`、`RPop()`、`LRange()`
- `Incr()`、`IncrBy()`、`Decr()`

所有 Redis 操作都需要 context 参数。

### Common Gotchas

1. **Validator 中间件**：`validate.RegisterTranslation()` 的注册函数需要 `RegisterTranslationsFunc` 类型，签名是 `func(ut.Translator) error`，而不是 `func(ut.Translator, validator.FieldError) string`
2. **Redis Options**：使用 `ConnMaxIdleTime` 而不是已弃用的 `IdleTimeout`
3. **Redis Exists**：返回 `int64` 而不是 `bool` - 使用 `result > 0` 转换
4. **Gin Context Get**：返回 `(interface{}, bool)` - 使用类型断言或 comma-ok 语法
5. **Response Imports**：Controller 必须导入 `"online-mall/internal/utils"` 才能访问响应辅助函数
6. **Unused Imports**：提交前删除未使用的导入（如果需要使用 `_` 前缀）
