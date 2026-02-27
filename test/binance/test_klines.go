/*
币安K线数据测试程序

测试内容：
- 获取不同周期的K线数据
- 测试不同交易对
- 验证数据完整性

运行方式：
  go run test/binance/test_klines.go
*/
package main

import (
	"fmt"

	"crypto-ai-trader/binance"
	"crypto-ai-trader/config"
	"crypto-ai-trader/utils"

	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	if err := utils.Init("logs/app.log", "debug"); err != nil {
		panic(err)
	}
	defer utils.Sync()

	utils.Info("=== 币安K线数据测试开始 ===")

	// 加载配置
	cfg, err := config.Load("configs/config.yml")
	if err != nil {
		utils.Fatal("加载配置失败", zap.Error(err))
	}

	// 创建客户端（K线数据不需要签名，使用任意账号即可）
	acc := cfg.GetEnabledAccounts()[0]
	client := binance.NewClient(
		acc.APIKey,
		acc.APISecret,
		cfg.Binance.FuturesURL,
		cfg.GetProxyURL(),
	)

	// 测试连接
	utils.Info("测试连接...")
	if err := client.Ping(); err != nil {
		utils.Fatal("连接失败", zap.Error(err))
	}
	utils.Info("连接成功")
	fmt.Println()

	// 1. 获取BTCUSDT 1分钟K线（最近10根）
	fmt.Println("【1. BTCUSDT 1分钟K线（最近10根）】")
	klines1m, err := client.GetKlines("BTCUSDT", "1m", 10)
	if err != nil {
		utils.Error("获取K线失败", zap.Error(err))
	} else {
		fmt.Printf("  获取数量: %d\n", len(klines1m))
		if len(klines1m) > 0 {
			// 显示最新一根K线
			latest := klines1m[len(klines1m)-1]
			fmt.Println("  最新K线:")
			fmt.Printf("    开盘价: %s\n", latest.Open)
			fmt.Printf("    最高价: %s\n", latest.High)
			fmt.Printf("    最低价: %s\n", latest.Low)
			fmt.Printf("    收盘价: %s\n", latest.Close)
			fmt.Printf("    成交量: %s\n", latest.Volume)
			fmt.Printf("    成交笔数: %d\n", latest.NumberOfTrades)
		}
	}
	fmt.Println()

	// 2. 获取ETHUSDT 5分钟K线（最近20根）
	fmt.Println("【2. ETHUSDT 5分钟K线（最近20根）】")
	klines5m, err := client.GetKlines("ETHUSDT", "5m", 20)
	if err != nil {
		utils.Error("获取K线失败", zap.Error(err))
	} else {
		fmt.Printf("  获取数量: %d\n", len(klines5m))
		if len(klines5m) > 0 {
			latest := klines5m[len(klines5m)-1]
			fmt.Println("  最新K线:")
			fmt.Printf("    开盘价: %s\n", latest.Open)
			fmt.Printf("    收盘价: %s\n", latest.Close)
			fmt.Printf("    成交量: %s\n", latest.Volume)
		}
	}
	fmt.Println()

	// 3. 获取BTCUSDT 1小时K线（最近24根，即24小时）
	fmt.Println("【3. BTCUSDT 1小时K线（最近24根）】")
	klines1h, err := client.GetKlines("BTCUSDT", "1h", 24)
	if err != nil {
		utils.Error("获取K线失败", zap.Error(err))
	} else {
		fmt.Printf("  获取数量: %d\n", len(klines1h))
		if len(klines1h) >= 2 {
			// 显示第一根和最后一根
			first := klines1h[0]
			latest := klines1h[len(klines1h)-1]
			fmt.Println("  第一根K线:")
			fmt.Printf("    收盘价: %s\n", first.Close)
			fmt.Println("  最新K线:")
			fmt.Printf("    收盘价: %s\n", latest.Close)
			
			// 计算价格变化
			fmt.Printf("  24小时价格变化: %s -> %s\n", first.Close, latest.Close)
		}
	}
	fmt.Println()

	// 4. 获取BTCUSDT 4小时K线（最近10根）
	fmt.Println("【4. BTCUSDT 4小时K线（最近10根）】")
	klines4h, err := client.GetKlines("BTCUSDT", "4h", 10)
	if err != nil {
		utils.Error("获取K线失败", zap.Error(err))
	} else {
		fmt.Printf("  获取数量: %d\n", len(klines4h))
		if len(klines4h) > 0 {
			latest := klines4h[len(klines4h)-1]
			fmt.Println("  最新K线:")
			fmt.Printf("    开盘价: %s\n", latest.Open)
			fmt.Printf("    收盘价: %s\n", latest.Close)
			fmt.Printf("    成交量: %s\n", latest.Volume)
		}
	}
	fmt.Println()

	// 5. 获取BTCUSDT 1天K线（最近7根）
	fmt.Println("【5. BTCUSDT 1天K线（最近7根）】")
	klines1d, err := client.GetKlines("BTCUSDT", "1d", 7)
	if err != nil {
		utils.Error("获取K线失败", zap.Error(err))
	} else {
		fmt.Printf("  获取数量: %d\n", len(klines1d))
		fmt.Println("  最近7天收盘价:")
		for i, kline := range klines1d {
			fmt.Printf("    第%d天: %s\n", i+1, kline.Close)
		}
	}
	fmt.Println()

	utils.Info("=== 测试完成 ===")
}
