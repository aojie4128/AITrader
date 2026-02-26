# Utils 工具模块

## Logger 日志工具

基于 zap 的高性能日志系统。

### 特性

- ✅ 高性能（零内存分配）
- ✅ 结构化日志
- ✅ 多级别日志（Debug, Info, Warn, Error, Fatal）
- ✅ 双输出（控制台 + 文件）
- ✅ 控制台彩色输出
- ✅ 文件JSON格式输出
- ✅ 自动创建日志目录

### 使用方法

#### 1. 初始化

```go
import "crypto-ai-trader/utils"

// 初始化日志系统
err := utils.Init("logs/app.log", "debug")
if err != nil {
    panic(err)
}
defer utils.Sync() // 程序退出前同步日志
```

#### 2. 基本日志

```go
utils.Debug("调试信息")
utils.Info("普通信息")
utils.Warn("警告信息")
utils.Error("错误信息")
utils.Fatal("致命错误") // 会退出程序
```

#### 3. 结构化日志

```go
import "go.uber.org/zap"

// 添加字段
utils.Info("用户登录",
    zap.String("user_id", "account_1"),
    zap.String("ip", "127.0.0.1"),
    zap.Int("port", 8080),
)

// 交易日志
utils.Info("交易信号",
    zap.String("symbol", "BTCUSDT"),
    zap.String("side", "BUY"),
    zap.Float64("price", 45000.50),
    zap.Float64("quantity", 0.01),
)
```

#### 4. 常用字段类型

```go
zap.String("key", "value")          // 字符串
zap.Int("key", 123)                 // 整数
zap.Float64("key", 123.45)          // 浮点数
zap.Bool("key", true)               // 布尔值
zap.Duration("key", time.Second)    // 时间间隔
zap.Time("key", time.Now())         // 时间
zap.Error(err)                      // 错误
```

### 日志级别

- **debug**: 调试信息，开发环境使用
- **info**: 普通信息，默认级别
- **warn**: 警告信息
- **error**: 错误信息

### 输出格式

#### 控制台输出（彩色）
```
2026-02-26T22:28:23.356+0800    INFO    utils/test_logger.go:32 用户登录 {"user_id": "account_1"}
```

#### 文件输出（JSON）
```json
{"level":"INFO","time":"2026-02-26T22:28:23.356+0800","caller":"utils/test_logger.go:32","msg":"用户登录","user_id":"account_1"}
```

### 最佳实践

#### 1. 程序启动时初始化

```go
func main() {
    // 初始化日志
    if err := utils.Init("logs/app.log", "info"); err != nil {
        panic(err)
    }
    defer utils.Sync()
    
    utils.Info("程序启动")
    // ...
}
```

#### 2. 记录关键操作

```go
// 配置加载
utils.Info("配置加载成功", zap.Int("accounts", len(accounts)))

// API调用
utils.Debug("调用币安API", 
    zap.String("endpoint", "/api/v3/klines"),
    zap.String("symbol", "BTCUSDT"),
)

// 交易执行
utils.Info("下单成功",
    zap.String("order_id", "123456"),
    zap.String("symbol", "BTCUSDT"),
    zap.Float64("price", 45000),
)
```

#### 3. 错误处理

```go
if err != nil {
    utils.Error("操作失败",
        zap.String("operation", "place_order"),
        zap.Error(err),
    )
    return err
}
```

#### 4. 性能监控

```go
start := time.Now()
// ... 执行操作
duration := time.Since(start)

utils.Info("操作完成",
    zap.String("operation", "fetch_klines"),
    zap.Duration("duration", duration),
)
```

### 测试

```bash
go run test/utils/test_logger.go
```

查看日志文件：
```bash
cat logs/app.log
```

### 配置建议

- **开发环境**: `debug` 级别，查看详细信息
- **生产环境**: `info` 级别，减少日志量
- **日志文件**: 建议按日期轮转（后续可添加）

### 注意事项

1. 程序退出前务必调用 `utils.Sync()` 确保日志写入
2. 避免在高频循环中使用 `Debug` 日志
3. 敏感信息（API密钥）不要记录到日志
4. 日志文件会自动创建，无需手动创建目录
