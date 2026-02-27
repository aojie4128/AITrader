/*
指标计算模块测试程序

测试内容：
- 测试短线策略指标计算（1m, 5m, 15m）
- 测试中长线策略指标计算（1h, 4h, 1d）
- 验证各项技术指标的计算结果

运行方式：
  go run test/indicators/test_indicators.go
*/
package main

import (
	"encoding/json"
	"fmt"

	"crypto-ai-trader/binance"
	"crypto-ai-trader/config"
	"crypto-ai-trader/indicators"
	"crypto-ai-trader/utils"

	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	if err := utils.Init("logs/app.log", "debug"); err != nil {
		panic(err)
	}
	defer utils.Sync()

	utils.Info("=== 指标计算模块测试开始 ===")

	// 加载配置
	cfg, err := config.Load("configs/config.yml")
	if err != nil {
		utils.Fatal("加载配置失败", zap.Error(err))
	}

	// 创建客户端
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

	symbol := "BTCUSDT"

	// ========== 测试短线策略指标 ==========
	fmt.Println("【1. 短线策略指标测试】")
	fmt.Println("获取K线数据...")

	// 获取短线K线数据（需要足够的数据来计算指标）
	klines1m, err := client.GetKlines(symbol, "1m", 100)
	if err != nil {
		utils.Fatal("获取1m K线失败", zap.Error(err))
	}
	fmt.Printf("  1m K线: %d根\n", len(klines1m))

	klines5m, err := client.GetKlines(symbol, "5m", 100)
	if err != nil {
		utils.Fatal("获取5m K线失败", zap.Error(err))
	}
	fmt.Printf("  5m K线: %d根\n", len(klines5m))

	klines15m, err := client.GetKlines(symbol, "15m", 100)
	if err != nil {
		utils.Fatal("获取15m K线失败", zap.Error(err))
	}
	fmt.Printf("  15m K线: %d根\n", len(klines15m))

	// 计算短线指标
	fmt.Println("\n计算短线指标...")
	shortTerm := indicators.CalculateShortTermIndicators(symbol, klines1m, klines5m, klines15m)

	// 显示1分钟指标
	fmt.Println("\n【1分钟周期指标】")
	printTimeframeIndicators(shortTerm.M1)

	// 显示5分钟指标
	fmt.Println("\n【5分钟周期指标】")
	printTimeframeIndicators(shortTerm.M5)

	// 显示15分钟指标
	fmt.Println("\n【15分钟周期指标】")
	printTimeframeIndicators(shortTerm.M15)

	// 输出JSON格式（用于AI分析）
	fmt.Println("\n【短线指标JSON格式】")
	shortTermJSON, _ := json.MarshalIndent(shortTerm, "", "  ")
	fmt.Println(string(shortTermJSON))

	fmt.Println("\n" + "============================================================")

	// ========== 测试中长线策略指标 ==========
	fmt.Println("\n【2. 中长线策略指标测试】")
	fmt.Println("获取K线数据...")

	// 获取中长线K线数据
	klines1h, err := client.GetKlines(symbol, "1h", 100)
	if err != nil {
		utils.Fatal("获取1h K线失败", zap.Error(err))
	}
	fmt.Printf("  1h K线: %d根\n", len(klines1h))

	klines4h, err := client.GetKlines(symbol, "4h", 100)
	if err != nil {
		utils.Fatal("获取4h K线失败", zap.Error(err))
	}
	fmt.Printf("  4h K线: %d根\n", len(klines4h))

	klines1d, err := client.GetKlines(symbol, "1d", 100)
	if err != nil {
		utils.Fatal("获取1d K线失败", zap.Error(err))
	}
	fmt.Printf("  1d K线: %d根\n", len(klines1d))

	// 计算中长线指标
	fmt.Println("\n计算中长线指标...")
	longTerm := indicators.CalculateLongTermIndicators(symbol, klines1h, klines4h, klines1d)

	// 显示1小时指标
	fmt.Println("\n【1小时周期指标】")
	printTimeframeIndicators(longTerm.H1)

	// 显示4小时指标
	fmt.Println("\n【4小时周期指标】")
	printTimeframeIndicators(longTerm.H4)

	// 显示1天指标
	fmt.Println("\n【1天周期指标】")
	printTimeframeIndicators(longTerm.D1)

	// 输出JSON格式（用于AI分析）
	fmt.Println("\n【中长线指标JSON格式】")
	longTermJSON, _ := json.MarshalIndent(longTerm, "", "  ")
	fmt.Println(string(longTermJSON))

	utils.Info("=== 测试完成 ===")
}

// printTimeframeIndicators 打印单个时间周期的指标
func printTimeframeIndicators(ind *indicators.TimeframeIndicators) {
	if ind == nil {
		fmt.Println("  无数据")
		return
	}

	fmt.Printf("  当前价格: %.2f\n", ind.ClosePrice)
	fmt.Printf("  成交量: %.2f\n", ind.Volume)
	fmt.Println("\n  移动平均线:")
	fmt.Printf("    MA5:  %.2f\n", ind.MA5)
	fmt.Printf("    MA10: %.2f\n", ind.MA10)
	fmt.Printf("    MA20: %.2f\n", ind.MA20)
	fmt.Printf("    EMA12: %.2f\n", ind.EMA12)
	fmt.Printf("    EMA26: %.2f\n", ind.EMA26)

	fmt.Println("\n  震荡指标:")
	fmt.Printf("    RSI: %.2f\n", ind.RSI)

	if ind.KDJ != nil && len(ind.KDJ.K) > 0 {
		latest := len(ind.KDJ.K) - 1
		fmt.Printf("    KDJ: K=%.2f, D=%.2f, J=%.2f\n",
			ind.KDJ.K[latest],
			ind.KDJ.D[latest],
			ind.KDJ.J[latest])
	}

	fmt.Println("\n  趋势指标:")
	if ind.MACD != nil && len(ind.MACD.DIF) > 0 {
		latest := len(ind.MACD.DIF) - 1
		fmt.Printf("    MACD: DIF=%.4f, DEA=%.4f, Histogram=%.4f\n",
			ind.MACD.DIF[latest],
			ind.MACD.DEA[latest],
			ind.MACD.Histogram[latest])
	}

	fmt.Println("\n  波动率指标:")
	fmt.Printf("    ATR: %.2f\n", ind.ATR)
	if ind.BB != nil && len(ind.BB.Upper) > 0 {
		latest := len(ind.BB.Upper) - 1
		fmt.Printf("    布林带: 上轨=%.2f, 中轨=%.2f, 下轨=%.2f\n",
			ind.BB.Upper[latest],
			ind.BB.Middle[latest],
			ind.BB.Lower[latest])
	}
}
