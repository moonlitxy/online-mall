package utils

import (
	"context"
	"fmt"
	"log"
	"online-mall/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient Redis客户端
var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() error {
	cfg := config.GlobalConfig.Redis

	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.PoolSize / 2,
		DialTimeout:     10 * time.Second,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
		PoolTimeout:     3 * time.Second,
		ConnMaxIdleTime: 5 * time.Minute,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect redis: %v", err)
	}

	log.Println("Redis connected successfully")
	RedisClient = rdb
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// RedisKey Redis key常量
const (
	// 用户相关
	UserInfoKey    = "user:info:%d"    // 用户信息
	UserTokenKey   = "user:token:%s"   // 用户token
	UserCartKey    = "user:cart:%d"    // 用户购物车
	UserAddressKey = "user:address:%d" // 用户地址列表

	// 商品相关
	ProductInfoKey  = "product:info:%d" // 商品信息
	CategoryTreeKey = "category:tree"   // 分类树
	HotProductsKey  = "hot:products"    // 热门商品
	NewProductsKey  = "new:products"    // 新品商品

	// 订单相关
	OrderKey      = "order:%d"       // 订单信息
	UserOrdersKey = "user:orders:%d" // 用户订单列表

	// 优惠券相关
	CouponKey      = "coupon:%d"       // 优惠券
	UserCouponsKey = "user:coupons:%d" // 用户优惠券列表

	// 缓存通用
	CachePrefix = "online-mall:" // 缓存前缀
)

// Set 设置缓存
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func Get(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// Del 删除缓存
func Del(ctx context.Context, keys ...string) error {
	return RedisClient.Del(ctx, keys...).Err()
}

// Exists 检查key是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := RedisClient.Exists(ctx, key).Result()
	return result > 0, err
}

// Expire 设置过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return RedisClient.Expire(ctx, key, expiration).Err()
}

// HSet 设置哈希
func HSet(ctx context.Context, key string, field string, value interface{}) error {
	return RedisClient.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希
func HGet(ctx context.Context, key string, field string) (string, error) {
	return RedisClient.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return RedisClient.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func HDel(ctx context.Context, key string, fields ...string) error {
	return RedisClient.HDel(ctx, key, fields...).Err()
}

// LPush 列表左推
func LPush(ctx context.Context, key string, values ...interface{}) error {
	return RedisClient.LPush(ctx, key, values...).Err()
}

// RPush 列表右推
func RPush(ctx context.Context, key string, values ...interface{}) error {
	return RedisClient.RPush(ctx, key, values...).Err()
}

// LPop 列表左弹
func LPop(ctx context.Context, key string) (string, error) {
	return RedisClient.LPop(ctx, key).Result()
}

// RPop 列表右弹
func RPop(ctx context.Context, key string) (string, error) {
	return RedisClient.RPop(ctx, key).Result()
}

// LRange 获取列表范围
func LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return RedisClient.LRange(ctx, key, start, stop).Result()
}

// LLen 获取列表长度
func LLen(ctx context.Context, key string) (int64, error) {
	return RedisClient.LLen(ctx, key).Result()
}

// Incr 增加计数
func Incr(ctx context.Context, key string) (int64, error) {
	return RedisClient.Incr(ctx, key).Result()
}

// IncrBy 增加指定数值
func IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return RedisClient.IncrBy(ctx, key, value).Result()
}

// Decr 减少计数
func Decr(ctx context.Context, key string) (int64, error) {
	return RedisClient.Decr(ctx, key).Result()
}

// SetNx 设置Nx（key不存在时设置）
func SetNx(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return RedisClient.SetNX(ctx, key, value, expiration).Result()
}

// GetSet 获取并设置
func GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	return RedisClient.GetSet(ctx, key, value).Result()
}
