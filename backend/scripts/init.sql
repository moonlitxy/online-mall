-- 在线商城数据库初始化脚本

-- 创建数据库
CREATE DATABASE IF NOT EXISTS online_mall DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE online_mall;

-- 用户表
CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_phone` (`phone`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 收货地址表
CREATE TABLE `addresses` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '地址ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `name` varchar(50) NOT NULL COMMENT '收货人姓名',
  `phone` varchar(20) NOT NULL COMMENT '手机号',
  `province` varchar(50) NOT NULL COMMENT '省份',
  `city` varchar(50) NOT NULL COMMENT '城市',
  `district` varchar(50) NOT NULL COMMENT '区县',
  `detail` varchar(255) NOT NULL COMMENT '详细地址',
  `postcode` varchar(10) DEFAULT NULL COMMENT '邮编',
  `tag` varchar(20) DEFAULT NULL COMMENT '地址标签',
  `is_default` tinyint(1) DEFAULT 0 COMMENT '是否默认地址',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收货地址表';

-- 商品分类表
CREATE TABLE `categories` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `parent_id` bigint(20) unsigned DEFAULT 0 COMMENT '父分类ID',
  `level` tinyint(1) DEFAULT 1 COMMENT '层级',
  `sort` int(11) DEFAULT 0 COMMENT '排序',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态：1-显示，0-隐藏',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品分类表';

-- 商品表
CREATE TABLE `products` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  `name` varchar(255) NOT NULL COMMENT '商品名称',
  `category_id` bigint(20) unsigned NOT NULL COMMENT '分类ID',
  `description` text COMMENT '商品描述',
  `price` decimal(10,2) NOT NULL COMMENT '商品价格',
  `original_price` decimal(10,2) DEFAULT NULL COMMENT '原价',
  `stock` int(11) DEFAULT 0 COMMENT '库存',
  `sales` int(11) DEFAULT 0 COMMENT '销量',
  `images` text COMMENT '商品图片',
  `video_url` varchar(255) DEFAULT NULL COMMENT '视频链接',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态：1-上架，0-下架',
  `is_hot` tinyint(1) DEFAULT 0 COMMENT '是否热门',
  `is_new` tinyint(1) DEFAULT 0 COMMENT '是否新品',
  `sort` int(11) DEFAULT 0 COMMENT '排序',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_status` (`status`),
  KEY `idx_hot_new` (`is_hot`, `is_new`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- 商品SKU表
CREATE TABLE `product_skus` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'SKU ID',
  `product_id` bigint(20) unsigned NOT NULL COMMENT '商品ID',
  `name` varchar(255) NOT NULL COMMENT 'SKU名称',
  `specifications` text NOT NULL COMMENT '规格信息',
  `price` decimal(10,2) NOT NULL COMMENT 'SKU价格',
  `stock` int(11) DEFAULT 0 COMMENT 'SKU库存',
  `sales` int(11) DEFAULT 0 COMMENT 'SKU销量',
  `image` varchar(255) DEFAULT NULL COMMENT 'SKU图片',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品SKU表';

-- 订单表
CREATE TABLE `orders` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `order_no` varchar(32) NOT NULL COMMENT '订单号',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `address_id` bigint(20) unsigned NOT NULL COMMENT '地址ID',
  `total_amount` decimal(10,2) NOT NULL COMMENT '订单总金额',
  `freight` decimal(10,2) DEFAULT 0.00 COMMENT '运费',
  `discount_amount` decimal(10,2) DEFAULT 0.00 COMMENT '优惠金额',
  `pay_amount` decimal(10,2) NOT NULL COMMENT '支付金额',
  `pay_status` tinyint(1) DEFAULT 0 COMMENT '支付状态：0-未支付，1-已支付',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `payment_method` varchar(20) DEFAULT NULL COMMENT '支付方式',
  `order_status` tinyint(1) DEFAULT 0 COMMENT '订单状态：0-待付款，1-待发货，2-待收货，3-已完成，4-已取消',
  `cancel_reason` varchar(255) DEFAULT NULL COMMENT '取消原因',
  `cancel_time` datetime DEFAULT NULL COMMENT '取消时间',
  `remark` varchar(255) DEFAULT NULL COMMENT '订单备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_order_status` (`order_status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- 订单商品表
CREATE TABLE `order_items` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '订单商品ID',
  `order_id` bigint(20) unsigned NOT NULL COMMENT '订单ID',
  `product_id` bigint(20) unsigned NOT NULL COMMENT '商品ID',
  `sku_id` bigint(20) unsigned NOT NULL COMMENT 'SKU ID',
  `product_name` varchar(255) NOT NULL COMMENT '商品名称',
  `product_image` varchar(255) DEFAULT NULL COMMENT '商品图片',
  `specifications` text COMMENT '规格信息',
  `price` decimal(10,2) NOT NULL COMMENT '商品价格',
  `quantity` int(11) NOT NULL COMMENT '购买数量',
  `total_amount` decimal(10,2) NOT NULL COMMENT '小计金额',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单商品表';

-- 购物车表
CREATE TABLE `cart_items` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '购物车项ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `product_id` bigint(20) unsigned NOT NULL COMMENT '商品ID',
  `sku_id` bigint(20) unsigned NOT NULL COMMENT 'SKU ID',
  `quantity` int(11) NOT NULL COMMENT '商品数量',
  `selected` tinyint(1) DEFAULT 1 COMMENT '是否选中',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购物车表';

-- 优惠券表
CREATE TABLE `coupons` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '优惠券ID',
  `name` varchar(100) NOT NULL COMMENT '优惠券名称',
  `type` tinyint(1) NOT NULL COMMENT '类型：1-满减券，2-折扣券',
  `value` decimal(10,2) NOT NULL COMMENT '金额或折扣率',
  `min_amount` decimal(10,2) DEFAULT 0.00 COMMENT '最低使用金额',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `end_time` datetime NOT NULL COMMENT '结束时间',
  `stock` int(11) DEFAULT 0 COMMENT '库存',
  `used_count` int(11) DEFAULT 0 COMMENT '已使用数量',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态：1-可用，0-停用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='优惠券表';

-- 用户优惠券表
CREATE TABLE `user_coupons` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户优惠券ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `coupon_id` bigint(20) unsigned NOT NULL COMMENT '优惠券ID',
  `order_id` bigint(20) unsigned DEFAULT NULL COMMENT '使用的订单ID',
  `status` tinyint(1) DEFAULT 0 COMMENT '状态：0-未使用，1-已使用，2-已过期',
  `used_time` datetime DEFAULT NULL COMMENT '使用时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_coupon` (`user_id`, `coupon_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_coupon_id` (`coupon_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户优惠券表';

-- 插入测试数据

-- 插入管理员用户
INSERT INTO `users` (`username`, `password`, `nickname`, `status`, `role`) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '管理员', 1);

-- 插入商品分类
INSERT INTO `categories` (`name`, `parent_id`, `level`, `sort`, `status`) VALUES
('电子产品', 0, 1, 1, 1),
('服装鞋帽', 0, 1, 2, 1),
('食品饮料', 0, 1, 3, 1),
('家用电器', 0, 1, 4, 1);

-- 插入测试商品
INSERT INTO `products` (`name`, `category_id`, `description`, `price`, `original_price`, `stock`, `sales`, `status`, `is_hot`, `is_new`, `sort`) VALUES
('苹果 iPhone 14 Pro', 1, '苹果最新款手机', 7999.00, 8999.00, 100, 50, 1, 1, 1, 1),
('小米13', 1, '小米旗舰手机', 3999.00, 4299.00, 200, 100, 1, 1, 1, 2),
('男士运动鞋', 2, '舒适透气', 299.00, 399.00, 50, 20, 1, 0, 1, 1),
('女士连衣裙', 2, '时尚潮流', 199.00, 299.00, 30, 15, 1, 0, 1, 2),
('可口可乐', 3, '经典口味', 3.00, 4.00, 500, 300, 1, 1, 0, 1),
('百事可乐', 3, '清爽畅饮', 3.00, 4.00, 500, 280, 1, 0, 0, 2);

-- 插入测试优惠券
INSERT INTO `coupons` (`name`, `type`, `value`, `min_amount`, `start_time`, `end_time`, `stock`, `status`) VALUES
('新人专享券', 1, 50.00, 299.00, '2024-01-01 00:00:00', '2024-12-31 23:59:59', 1000, 1),
('满减优惠券', 1, 100.00, 599.00, '2024-01-01 00:00:00', '2024-12-31 23:59:59', 500, 1),
('折扣券', 2, 0.90, 99.00, '2024-01-01 00:00:00', '2024-12-31 23:59:59', 200, 1);