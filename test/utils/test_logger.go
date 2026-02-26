/*
日志模块测试程序

测试内容：
- 初始化日志系统
- 测试不同级别的日志输出
- 测试结构化日志字段
- 测试文件输出
- 测试控制台彩色输出

运行方式：
  go run test/utils/test_logger.go
*/
package main

import (
	"crypto-ai-trader/utils"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 初始化日志系统
	err := utils.Init("logs/app.log", "debug")
	if err != nil {
		panic(err)
	}
	defer utils.Sync()

	utils.Info("=== 日志模块测试开始 ===")

	// 测试不同级别的日志
	utils.Debug("这是一条调试日志")
	utils.Info("这是一条信息日志")
	utils.Warn("这是一条警告日志")
	utils.Error("这是一条错误日志")

	// 测试结构化日志
	utils.Info("用户登录",
		zap.String("user_id", "account_1"),
		zap.String("ip", "127.0.0.1"),
		zap.Int("port", 8080),
	)

	// 测试交易相关日志
	utils.Info("交易信号",
		zap.String("symbol", "BTCUSDT"),
		zap.String("side", "BUY"),
		zap.Float64("price", 45000.50),
		zap.Float64("quantity", 0.01),
		zap.String("strategy", "short_term"),
	)

	// 测试错误日志
	utils.Error("API调用失败",
		zap.String("api", "binance"),
		zap.String("endpoint", "/api/v3/klines"),
		zap.String("error", "connection timeout"),
		zap.Duration("duration", 30*time.Second),
	)

	// 测试警告日志
	utils.Warn("账户余额不足",
		zap.String("account_id", "account_1"),
		zap.Float64("balance", 100.50),
		zap.Float64("required", 500.00),
	)

	// 测试嵌套字段
	utils.Info("AI分析完成",
		zap.String("account", "短线-简洁版"),
		zap.Object("result", zapObject{
			Signal:     "BUY",
			Confidence: 0.85,
			Reason:     "RSI超卖且MACD金叉",
		}),
	)

	utils.Info("=== 日志模块测试完成 ===")
	utils.Info("日志文件位置: logs/app.log")
}

// zapObject 实现 zapcore.ObjectMarshaler 接口
type zapObject struct {
	Signal     string
	Confidence float64
	Reason     string
}

func (o zapObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("signal", o.Signal)
	enc.AddFloat64("confidence", o.Confidence)
	enc.AddString("reason", o.Reason)
	return nil
}
