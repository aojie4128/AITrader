/*
指标计算模块测试程序

测试内容：
- 测试短线策略指标计算（1h → 15m → 5m）
- 测试中长线策略指标计算（4h → 1h → 15m）
- 验证各项技术指标的计算结果
- 输出JSON格式供AI分析使用

运行方式：
  go run test/indicators/test_indicators.go
*/
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"crypto-ai-trader/binance"
	"crypto-ai-trader/config"
	"crypto-ai-trader/indicators"
	"crypto-ai-trader/utils"

	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	if err := utils.Init("logs/app.log", "info"); err != nil {
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
	fmt.Println("✓ 连接成功\n")

	symbol := "BTCUSDT"

	// ========== 测试短线策略指标 ==========
	fmt.Println("【短线策略指标测试】持仓30-90分钟")
	fmt.Println("时间周期：1h（方向过滤） → 15m（主分析） → 5m（入场）")
	fmt.Println()

	fmt.Println("正在获取K线数据...")
	klines1h_short, err := client.GetKlines(symbol, "1h", 100)
	if err != nil {
		utils.Fatal("获取1h K线失败", zap.Error(err))
	}
	fmt.Printf("  ✓ 1h K线: %d根\n", len(klines1h_short))

	klines15m_short, err := client.GetKlines(symbol, "15m", 100)
	if err != nil {
		utils.Fatal("获取15m K线失败", zap.Error(err))
	}
	fmt.Printf("  ✓ 15m K线: %d根\n", len(klines15m_short))

	klines5m, err := client.GetKlines(symbol, "5m", 100)
	if err != nil {
		utils.Fatal("获取5m K线失败", zap.Error(err))
	}
	fmt.Printf("  ✓ 5m K线: %d根\n", len(klines5m))

	fmt.Println("\n正在计算短线指标...")
	shortTerm := indicators.CalculateShortTermIndicators(symbol, klines1h_short, klines15m_short, klines5m)
	if shortTerm == nil {
		utils.Fatal("短线指标计算失败")
	}
	fmt.Println("  ✓ 计算完成")

	// 显示短线指标
	fmt.Println("\n【1小时周期 - 方向过滤】")
	printTimeframeIndicators(shortTerm.Timeframes.H1)

	fmt.Println("\n【15分钟周期 - 主分析】")
	printTimeframeIndicators(shortTerm.Timeframes.M15)

	fmt.Println("\n【5分钟周期 - 入场】")
	printTimeframeIndicators(shortTerm.Timeframes.M5)

	// 输出JSON格式
	fmt.Println("\n【短线指标JSON格式（供AI分析）】")
	shortTermJSON, _ := json.MarshalIndent(shortTerm, "", "  ")
	fmt.Println(string(shortTermJSON))

	fmt.Println("\n" + "================================================================================")

	// ========== 测试短线策略指标（含市场数据） ==========
	fmt.Println("\n【短线策略指标测试（含市场数据）】")
	
	// 模拟OI缓存（实际应用中应该从数据库或缓存中读取）
	oiCache := &indicators.OICache{
		History:    []float64{5363.02, 5350.15, 5340.28, 5330.42, 5320.55},
		Timestamps: []int64{time.Now().Unix(), time.Now().Unix() - 300, time.Now().Unix() - 600, time.Now().Unix() - 900, time.Now().Unix() - 1200},
	}
	
	fmt.Println("正在计算短线指标（含市场数据）...")
	shortTermWithMarket := indicators.CalculateShortTermIndicatorsWithMarket(symbol, klines1h_short, klines15m_short, klines5m, client, oiCache)
	if shortTermWithMarket == nil {
		utils.Fatal("短线指标（含市场数据）计算失败")
	}
	fmt.Println("  ✓ 计算完成")
	
	// 显示市场数据
	if shortTermWithMarket.MarketData != nil {
		fmt.Println("\n【市场数据】")
		md := shortTermWithMarket.MarketData
		fmt.Printf("  当前持仓量: $%.2f M\n", md.OICurrent)
		if len(md.OIHistory) > 0 {
			fmt.Printf("  历史持仓量: %v\n", md.OIHistory)
		}
		if md.OIChange5m != nil {
			fmt.Printf("  5分钟变化: %.2f%%\n", *md.OIChange5m)
		}
		if md.OIChange15m != nil {
			fmt.Printf("  15分钟变化: %.2f%%\n", *md.OIChange15m)
		}
		if md.OIChange25m != nil {
			fmt.Printf("  25分钟变化: %.2f%%\n", *md.OIChange25m)
		}
		fmt.Printf("  当前资金费率: %.4f%%\n", md.FundingRate)
		fmt.Printf("  资金费率平均: %.4f%%\n", md.FundingAvg3)
	}
	
	// 输出完整JSON
	fmt.Println("\n【短线指标JSON格式（含市场数据）】")
	shortTermWithMarketJSON, _ := json.MarshalIndent(shortTermWithMarket, "", "  ")
	fmt.Println(string(shortTermWithMarketJSON))

	fmt.Println("\n" + "================================================================================")
	fmt.Println()

	// ========== 测试中长线策略指标 ==========
	fmt.Println("【中长线策略指标测试】持仓2-4小时")
	fmt.Println("时间周期：4h（大趋势） → 1h（主分析） → 15m（入场）")
	fmt.Println()

	fmt.Println("正在获取K线数据...")
	klines4h, err := client.GetKlines(symbol, "4h", 100)
	if err != nil {
		utils.Fatal("获取4h K线失败", zap.Error(err))
	}
	fmt.Printf("  ✓ 4h K线: %d根\n", len(klines4h))

	klines1h_long, err := client.GetKlines(symbol, "1h", 100)
	if err != nil {
		utils.Fatal("获取1h K线失败", zap.Error(err))
	}
	fmt.Printf("  ✓ 1h K线: %d根\n", len(klines1h_long))

	klines15m_long, err := client.GetKlines(symbol, "15m", 100)
	if err != nil {
		utils.Fatal("获取15m K线失败", zap.Error(err))
	}
	fmt.Printf("  ✓ 15m K线: %d根\n", len(klines15m_long))

	fmt.Println("\n正在计算中长线指标...")
	longTerm := indicators.CalculateLongTermIndicators(symbol, klines4h, klines1h_long, klines15m_long)
	if longTerm == nil {
		utils.Fatal("中长线指标计算失败")
	}
	fmt.Println("  ✓ 计算完成")

	// 显示中长线指标
	fmt.Println("\n【4小时周期 - 大趋势】")
	printTimeframeIndicators(longTerm.Timeframes.H4)

	fmt.Println("\n【1小时周期 - 主分析】")
	printTimeframeIndicators(longTerm.Timeframes.H1)

	fmt.Println("\n【15分钟周期 - 入场】")
	printTimeframeIndicators(longTerm.Timeframes.M15)

	// 输出JSON格式
	fmt.Println("\n【中长线指标JSON格式（供AI分析）】")
	longTermJSON, _ := json.MarshalIndent(longTerm, "", "  ")
	fmt.Println(string(longTermJSON))

	fmt.Println("\n" + "================================================================================")
	utils.Info("=== 测试完成 ===")
}

// printTimeframeIndicators 打印单个时间周期的指标
func printTimeframeIndicators(data *indicators.TimeframeData) {
	if data == nil {
		fmt.Println("  无数据")
		return
	}

	// 价格信息
	fmt.Printf("  价格: 开=%.2f 高=%.2f 低=%.2f 收=%.2f\n",
		data.OpenPrice, data.HighPrice, data.LowPrice, data.ClosePrice)
	fmt.Printf("  成交量: %.2f\n", data.Volume)

	// 趋势指标
	fmt.Println("\n  【趋势指标】")
	fmt.Printf("    EMA9:  %.2f\n", data.EMA9)
	fmt.Printf("    EMA21: %.2f\n", data.EMA21)
	fmt.Printf("    EMA55: %.2f\n", data.EMA55)

	// 判断趋势
	if data.ClosePrice > data.EMA55 {
		fmt.Printf("    → 价格在EMA55上方，多头趋势\n")
	} else {
		fmt.Printf("    → 价格在EMA55下方，空头趋势\n")
	}

	// 动能指标
	fmt.Println("\n  【动能指标】")
	fmt.Printf("    RSI: %.2f", data.RSI)
	if data.RSI > 70 {
		fmt.Printf(" (超买)\n")
	} else if data.RSI < 30 {
		fmt.Printf(" (超卖)\n")
	} else if data.RSI > 50 {
		fmt.Printf(" (偏强)\n")
	} else {
		fmt.Printf(" (偏弱)\n")
	}

	if data.MACD != nil {
		fmt.Printf("    MACD: DIF=%.4f, DEA=%.4f, Histogram=%.4f\n",
			data.MACD.DIF, data.MACD.DEA, data.MACD.Histogram)
		if data.MACD.DIF > data.MACD.DEA {
			fmt.Printf("    → MACD金叉，多头动能\n")
		} else {
			fmt.Printf("    → MACD死叉，空头动能\n")
		}
	}

	// 波动率指标
	fmt.Println("\n  【波动率指标】")
	fmt.Printf("    ATR: %.2f\n", data.ATR)
	if data.BB != nil {
		fmt.Printf("    布林带: 上=%.2f 中=%.2f 下=%.2f\n",
			data.BB.Upper, data.BB.Middle, data.BB.Lower)

		// 判断价格位置
		bbWidth := data.BB.Upper - data.BB.Lower
		pricePosition := (data.ClosePrice - data.BB.Lower) / bbWidth * 100
		fmt.Printf("    → 价格位于布林带 %.1f%% 位置", pricePosition)

		if data.ClosePrice > data.BB.Upper {
			fmt.Printf(" (突破上轨)\n")
		} else if data.ClosePrice < data.BB.Lower {
			fmt.Printf(" (突破下轨)\n")
		} else if pricePosition > 80 {
			fmt.Printf(" (接近上轨)\n")
		} else if pricePosition < 20 {
			fmt.Printf(" (接近下轨)\n")
		} else {
			fmt.Printf(" (中性区域)\n")
		}
	}
}
