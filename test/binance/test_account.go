/*
币安账户信息测试程序

测试内容：
- 获取账户信息
- 获取账户余额
- 获取持仓信息
- 获取持仓风险

运行方式：
  go run test/binance/test_account.go
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

	utils.Info("=== 币安账户信息测试开始 ===")

	// 加载配置
	cfg, err := config.Load("configs/config.yml")
	if err != nil {
		utils.Fatal("加载配置失败", zap.Error(err))
	}

	// 获取第一个启用的账号进行测试
	enabledAccounts := cfg.GetEnabledAccounts()
	if len(enabledAccounts) == 0 {
		utils.Fatal("没有启用的账号")
	}

	acc := enabledAccounts[0]
	utils.Info("使用账号进行测试",
		zap.String("账号ID", acc.ID),
		zap.String("账号名称", acc.Name),
	)

	// 创建客户端
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

	// 1. 获取账户信息
	fmt.Println("【1. 账户信息】")
	accountInfo, err := client.GetAccountInfo()
	if err != nil {
		utils.Error("获取账户信息失败", zap.Error(err))
	} else {
		fmt.Printf("  总余额: %s USDT\n", accountInfo.TotalWalletBalance)
		fmt.Printf("  可用余额: %s USDT\n", accountInfo.AvailableBalance)
		fmt.Printf("  保证金余额: %s USDT\n", accountInfo.TotalMarginBalance)
		fmt.Printf("  未实现盈亏: %s USDT\n", accountInfo.TotalUnrealizedProfit)
		fmt.Println()
		fmt.Println("  【USDT资产详情】")
		fmt.Printf("    钱包余额: %s\n", accountInfo.Asset.WalletBalance)
		fmt.Printf("    可用余额: %s\n", accountInfo.Asset.AvailableBalance)
		fmt.Printf("    未实现盈亏: %s\n", accountInfo.Asset.UnrealizedProfit)
		fmt.Printf("    保证金余额: %s\n", accountInfo.Asset.MarginBalance)
		fmt.Printf("    最大可提现: %s\n", accountInfo.Asset.MaxWithdrawAmount)
	}
	fmt.Println()

	// 2. 获取USDT余额
	fmt.Println("【2. USDT余额】")
	balance, err := client.GetBalance()
	if err != nil {
		utils.Error("获取账户余额失败", zap.Error(err))
	} else {
		fmt.Printf("  资产: %s\n", balance.Asset)
		fmt.Printf("  余额: %s\n", balance.Balance)
		fmt.Printf("  可用余额: %s\n", balance.AvailableBalance)
		fmt.Printf("  未实现盈亏: %s\n", balance.UnrealizedProfit)
	}
	fmt.Println()

	// 3. 获取持仓信息
	fmt.Println("【3. 持仓信息】")
	positions, err := client.GetPositions()
	if err != nil {
		utils.Error("获取持仓信息失败", zap.Error(err))
	} else {
		if len(positions) == 0 {
			fmt.Println("  当前无持仓")
		} else {
			fmt.Printf("  持仓数量: %d\n", len(positions))
			for _, pos := range positions {
				fmt.Printf("  - %s:\n", pos.Symbol)
				fmt.Printf("      持仓数量: %s\n", pos.PositionAmt)
				fmt.Printf("      开仓均价: %s\n", pos.EntryPrice)
				fmt.Printf("      标记价格: %s\n", pos.MarkPrice)
				fmt.Printf("      未实现盈亏: %s\n", pos.UnRealizedProfit)
				fmt.Printf("      杠杆倍数: %s\n", pos.Leverage)
				fmt.Printf("      持仓方向: %s\n", pos.PositionSide)
				fmt.Printf("      保证金模式: %s\n", pos.MarginType)
				fmt.Println()
			}
		}
	}

	// 4. 获取持仓风险（所有交易对）
	fmt.Println("【4. 持仓风险】")
	positionRisks, err := client.GetPositionRisk("")
	if err != nil {
		utils.Error("获取持仓风险失败", zap.Error(err))
	} else {
		// 只显示有持仓的
		hasPosition := false
		for _, risk := range positionRisks {
			// 过滤掉0持仓
			if risk.PositionAmt != "0" && risk.PositionAmt != "0.0" && risk.PositionAmt != "0.00" && risk.PositionAmt != "0.000" {
				hasPosition = true
				fmt.Printf("  - %s:\n", risk.Symbol)
				fmt.Printf("      持仓数量: %s\n", risk.PositionAmt)
				fmt.Printf("      强平价格: %s\n", risk.LiquidationPrice)
				fmt.Printf("      未实现盈亏: %s\n", risk.UnRealizedProfit)
				fmt.Println()
			}
		}
		if !hasPosition {
			fmt.Println("  当前无持仓风险")
		}
	}

	utils.Info("=== 测试完成 ===")
}
