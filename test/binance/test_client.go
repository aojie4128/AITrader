/*
币安客户端测试程序

测试内容：
- 创建币安客户端
- 测试代理设置
- 测试连接（Ping）
- 测试多账号客户端创建

运行方式：
  go run test/binance/test_client.go
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

	utils.Info("=== 币安客户端测试开始 ===")

	// 加载配置
	cfg, err := config.Load("configs/config.yml")
	if err != nil {
		utils.Fatal("加载配置失败", zap.Error(err))
	}

	utils.Info("配置加载成功",
		zap.Int("账号数量", len(cfg.Accounts)),
		zap.String("币安URL", cfg.Binance.FuturesURL),
	)

	// 获取代理URL
	proxyURL := cfg.GetProxyURL()
	if proxyURL != "" {
		utils.Info("使用代理", zap.String("proxy", proxyURL))
	}

	// 测试：为每个启用的账号创建客户端
	enabledAccounts := cfg.GetEnabledAccounts()
	utils.Info("启用的账号", zap.Int("数量", len(enabledAccounts)))

	clients := make(map[string]*binance.Client)

	for _, acc := range enabledAccounts {
		utils.Info("创建客户端",
			zap.String("账号ID", acc.ID),
			zap.String("账号名称", acc.Name),
			zap.String("策略", acc.GetStrategyName()),
			zap.String("提示词", acc.GetPromptTypeName()),
		)

		// 创建客户端
		client := binance.NewClient(
			acc.APIKey,
			acc.APISecret,
			cfg.Binance.FuturesURL,
			proxyURL,
		)

		clients[acc.ID] = client

		// 测试连接
		utils.Info("测试连接", zap.String("账号", acc.Name))
		if err := client.Ping(); err != nil {
			utils.Error("连接测试失败",
				zap.String("账号", acc.Name),
				zap.Error(err),
			)
			continue
		}

		utils.Info("连接测试成功", zap.String("账号", acc.Name))
		fmt.Println()
	}

	// 总结
	utils.Info("=== 测试完成 ===")
	utils.Info("客户端创建统计",
		zap.Int("总数", len(clients)),
		zap.Int("成功", len(clients)),
	)

	// 显示所有客户端
	fmt.Println("\n创建的客户端列表：")
	for id := range clients {
		acc := cfg.GetAccountByID(id)
		fmt.Printf("  - %s (%s)\n", acc.Name, id)
	}
}
