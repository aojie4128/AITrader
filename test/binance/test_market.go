/*
币安市场数据API测试程序

测试内容：
- 获取持仓量（Open Interest）
- 获取资金费率历史
- 获取当前资金费率和标记价格
- 计算持仓量变化率

运行方式：
  go run test/binance/test_market.go
*/
package main

import (
	"fmt"
	"strconv"

	"crypto-ai-trader/binance"
	"crypto-ai-trader/config"
	"crypto-ai-trader/utils"

	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	if err := utils.Init("logs/app.log", "info"); err != nil {
		panic(err)
	}
	defer utils.Sync()

	utils.Info("=== 币安市场数据API测试开始 ===")

	// 加载配置
	cfg, err := config.Load("configs/config.yml")
	if err != nil {
		utils.Fatal("加载配置失败", zap.Error(err))
	}

	// 创建客户端（市场数据不需要签名）
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

	// ========== 1. 获取持仓量 ==========
	fmt.Println("【1. 获取持仓量（Open Interest）】")
	oi, err := client.GetOpenInterest(symbol)
	if err != nil {
		utils.Error("获取持仓量失败", zap.Error(err))
	} else {
		fmt.Printf("  交易对: %s\n", oi.Symbol)
		fmt.Printf("  持仓量（张数）: %s\n", oi.OpenInterest)
		
		// 解析数值并计算USDT价值
		oiContracts, _ := strconv.ParseFloat(oi.OpenInterest, 64)
		
		// 获取当前价格
		klines, err := client.GetKlines(symbol, "1m", 1)
		if err == nil && len(klines) > 0 {
			currentPrice, _ := strconv.ParseFloat(klines[0].Close, 64)
			oiValue := oiContracts * currentPrice
			fmt.Printf("  当前价格: $%.2f\n", currentPrice)
			fmt.Printf("  持仓量（USDT）: $%.2f\n", oiValue)
			fmt.Printf("  持仓量价值: $%.2f M\n", oiValue/1000000)
		}
	}
	fmt.Println()

	// ========== 2. 获取资金费率历史 ==========
	fmt.Println("【2. 获取资金费率历史（最近10次）】")
	fundingRates, err := client.GetFundingRateHistory(symbol, 10)
	if err != nil {
		utils.Error("获取资金费率历史失败", zap.Error(err))
	} else {
		fmt.Printf("  获取数量: %d\n", len(fundingRates))
		fmt.Println("  最近5次资金费率:")
		
		// 显示最近5次
		count := 5
		if len(fundingRates) < count {
			count = len(fundingRates)
		}
		
		for i := len(fundingRates) - count; i < len(fundingRates); i++ {
			fr := fundingRates[i]
			rate, _ := strconv.ParseFloat(fr.FundingRate, 64)
			ratePercent := rate * 100
			
			fmt.Printf("    %d. 费率: %.4f%% ", len(fundingRates)-i, ratePercent)
			if rate > 0 {
				fmt.Printf("(多头支付空头)\n")
			} else if rate < 0 {
				fmt.Printf("(空头支付多头)\n")
			} else {
				fmt.Printf("(平衡)\n")
			}
		}
		
		// 计算平均资金费率
		if len(fundingRates) >= 3 {
			sum := 0.0
			for i := len(fundingRates) - 3; i < len(fundingRates); i++ {
				rate, _ := strconv.ParseFloat(fundingRates[i].FundingRate, 64)
				sum += rate
			}
			avg := (sum / 3) * 100
			fmt.Printf("\n  最近3次平均费率: %.4f%%\n", avg)
			
			// 判断市场情绪
			if avg > 0.05 {
				fmt.Println("  → 市场情绪: 过度做多，注意回调风险")
			} else if avg < -0.05 {
				fmt.Println("  → 市场情绪: 过度做空，注意反弹风险")
			} else if avg > 0 {
				fmt.Println("  → 市场情绪: 多头占优")
			} else if avg < 0 {
				fmt.Println("  → 市场情绪: 空头占优")
			} else {
				fmt.Println("  → 市场情绪: 平衡")
			}
		}
	}
	fmt.Println()

	// ========== 3. 获取当前资金费率和标记价格 ==========
	fmt.Println("【3. 获取当前资金费率和标记价格】")
	premium, err := client.GetPremiumIndex(symbol)
	if err != nil {
		utils.Error("获取溢价指数失败", zap.Error(err))
	} else {
		fmt.Printf("  交易对: %s\n", premium.Symbol)
		fmt.Printf("  标记价格: %s\n", premium.MarkPrice)
		fmt.Printf("  指数价格: %s\n", premium.IndexPrice)
		
		rate, _ := strconv.ParseFloat(premium.LastFundingRate, 64)
		ratePercent := rate * 100
		fmt.Printf("  当前资金费率: %.4f%%\n", ratePercent)
		
		// 判断费率水平
		if rate > 0.1 {
			fmt.Println("  → 费率极高，市场过度做多")
		} else if rate > 0.05 {
			fmt.Println("  → 费率偏高，多头占优")
		} else if rate > 0 {
			fmt.Println("  → 费率正常，多头略占优")
		} else if rate > -0.05 {
			fmt.Println("  → 费率正常，空头略占优")
		} else if rate > -0.1 {
			fmt.Println("  → 费率偏低，空头占优")
		} else {
			fmt.Println("  → 费率极低，市场过度做空")
		}
	}
	fmt.Println()

	// ========== 4. 持仓量说明 ==========
	fmt.Println("【4. 持仓量变化率说明】")
	fmt.Println("  持仓量变化率需要历史数据支持")
	fmt.Println("  币安API只提供当前持仓量，不提供历史数据")
	fmt.Println("  ")
	fmt.Println("  解决方案：")
	fmt.Println("  1. 在database模块中定期存储OI数据（推荐）")
	fmt.Println("  2. 使用第三方数据源（如Coinglass）")
	fmt.Println("  ")
	fmt.Println("  当前可用的市场数据指标：")
	fmt.Println("  ✓ 当前持仓量（OI）- 判断市场规模")
	fmt.Println("  ✓ 当前资金费率 - 判断市场情绪")
	fmt.Println("  ✓ 资金费率历史平均 - 判断情绪趋势")
	fmt.Println("  ")
	fmt.Println("  这些指标已经足够用于交易决策！")
	fmt.Println()

	utils.Info("=== 测试完成 ===")
}
