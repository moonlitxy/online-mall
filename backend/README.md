# 在线商城后端项目

## 项目简介
基于Go 1.24.0 + Gin + MySQL + Redis技术栈开发的在线商城后端系统。

## 技术栈
- **后端框架**: Gin (v1.9.1)
- **数据库**: MySQL 8.0
- **缓存**: Redis 6.0
- **ORM**: GORM (v1.25.5)
- **JWT**: golang-jwt/jwt (v5.x)
- **验证器**: validator (v10.x)
- **配置管理**: viper (v1.16.0)
- **密码加密**: bcrypt

## 目录结构
```
backend/
├── cmd/                    # 应用入口
│   └── server/            # 主程序
├── internal/              # 内部包
│   ├── api/               # API处理器
│   │   ├── controller/     # 控制器层
│   │   ├── middleware/    # 中间件
│   │   └── routes/        # 路由配置
│   ├── service/           # 业务逻辑层
│   ├── repository/        # 数据访问层
│   ├── models/           # 数据模型
│   ├── config/           # 配置管理
│   ├── utils/            # 工具函数
│   └── pkg/              # 内部包
├── configs/              # 配置文件
│   ├── config.yaml
│   └── config.dev.yaml
├── scripts/              # 脚本
│   ├── start.bat         # Windows启动脚本
│   ├── start.sh          # Linux启动脚本
│   └── init.sql         # 数据库初始化脚本
├── logs/                # 日志目录
├── uploads/             # 文件上传目录
├── static/             # 静态文件目录
├── go.mod
├── go.sum
└── main.go
```

## 环境要求
- Go 1.17+
- MySQL 8.0+
- Redis 6.0+
- Windows/Linux/macOS

## 快速开始

### 1. 克隆项目
```bash
cd backend
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 配置数据库
- 修改 `configs/config.yaml` 中的数据库配置
- 确保MySQL和Redis服务已启动

### 4. 初始化数据库
```bash
# 在MySQL中执行 scripts/init.sql
mysql -u root -p < scripts/init.sql
```

### 5. 启动服务

#### Windows
```bash
scripts\start.bat
```

#### Linux/macOS
```bash
chmod +x scripts/start.sh
./scripts/start.sh
```

#### 直接运行
```bash
go run cmd/server/main.go
```

### 6. 访问服务
- API地址: http://localhost:8080
- 健康检查: http://localhost:8080/health

## 配置说明

### 数据库配置
```yaml
database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: 123456
  dbname: online_mall
```

### Redis配置
```yaml
redis:
  host: localhost
  port: 6379
  password:
  db: 0
  pool_size: 100
```

### JWT配置
```yaml
jwt:
  secret: online-mall-jwt-secret-key-2024
  expire_hours: 24
  issuer: online-mall
```

## API接口文档

### 认证相关
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/refresh-token` - 刷新token
- `POST /api/auth/logout` - 用户登出

### 用户管理
- `GET /api/users/profile` - 获取用户信息
- `PUT /api/users/profile` - 更新用户信息
- `PUT /api/users/password` - 修改密码

### 商品管理
- `GET /api/products` - 商品列表
- `GET /api/products/:id` - 商品详情
- `GET /api/products/:id/skus` - 商品SKU列表

### 购物车管理
- `GET /api/cart` - 购物车列表
- `POST /api/cart` - 添加到购物车
- `PUT /api/cart/:id/quantity` - 更新商品数量
- `PUT /api/cart/:id/selected` - 选择/取消选择
- `DELETE /api/cart/:id` - 删除购物车商品

### 订单管理
- `GET /api/orders` - 订单列表
- `GET /api/orders/:id` - 订单详情
- `POST /api/orders` - 创建订单
- `PUT /api/orders/:id/cancel` - 取消订单
- `PUT /api/orders/:id/receive` - 确认收货

### 地址管理
- `GET /api/addresses` - 地址列表
- `POST /api/addresses` - 添加地址
- `PUT /api/addresses/:id` - 更新地址
- `DELETE /api/addresses/:id` - 删除地址
- `PUT /api/addresses/:id/default` - 设置默认地址

### 优惠券管理
- `GET /api/coupons` - 优惠券列表
- `POST /api/coupons/:id/receive` - 领取优惠券
- `GET /api/my/coupons` - 我的优惠券

## 开发说明

### 代码规范
- 遵循Go语言官方代码规范
- 使用gofmt格式化代码
- 注释清晰完整

### 提交规范
```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码格式调整
refactor: 重构
test: 测试相关
chore: 构建/工具相关
```

## 测试

### 单元测试
```bash
go test ./...
```

### 性能测试
```bash
go test -bench=. -benchmem
```

## 部署

### 构建二进制文件
```bash
go build -o online-mall cmd/server/main.go
```

### 运行二进制文件
```bash
./online-mall
```

### Docker部署
```bash
# 构建镜像
docker build -t online-mall-backend .

# 运行容器
docker run -p 8080:8080 online-mall-backend
```

## 常见问题

### 1. 数据库连接失败
- 检查MySQL服务是否启动
- 检查配置文件中的数据库连接信息
- 确认数据库已创建

### 2. Redis连接失败
- 检查Redis服务是否启动
- 检查配置文件中的Redis连接信息

### 3. 依赖下载失败
- 配置Go代理: `go env -w GOPROXY=https://goproxy.cn,direct`
- 使用`go mod download`下载依赖

## 贡献指南
欢迎提交Issue和Pull Request！

## 许可证
MIT License

## 联系方式
- 作者: Online Mall Team
- 邮箱: support@online-mall.com