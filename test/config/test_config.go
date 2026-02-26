/*
配置模块测试程序

测试内容：
- 加载主配置文件（configs/config.yml）
- 加载账号配置文件（configs/accounts.yml）
- 验证代理配置
- 验证币安API配置
- 验证账号配置（4个账号）
- 测试获取启用的账号
- 测试根据ID获取账号
- 测试全局配置访问

运行方式：
  go run test/config/test_config.go
*/
package main

import (
	"fmt"
	"log"

	"crypto-ai-trader/config"
)

func main() {
	fmt.Println("=== 配置模块测试 ===\n")

	// 加载配置文件
	cfg, err := config.Load("configs/config.yml")
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	fmt.Println("✅ 配置加载成功！\n")

	// 打印代理配置
	fmt.Println("【代理配置】")
	fmt.Printf("  启用状态: %v\n", cfg.Proxy.IsUse)
	fmt.Printf("  主机地址: %s\n", cfg.Proxy.Host)
	fmt.Printf("  端口号: %d\n", cfg.Proxy.Port)
	if cfg.Proxy.IsUse {
		fmt.Printf("  代理URL: %s\n", cfg.GetProxyURL())
	}
	fmt.Println()

	// 打印币安配置
	fmt.Println("【币安API配置】")
	fmt.Printf("  合约URL: %s\n", cfg.Binance.FuturesURL)
	fmt.Println()

	// 打印账号配置
	fmt.Println("【账号配置】")
	fmt.Printf("  总账号数: %d\n", len(cfg.Accounts))
	fmt.Printf("  启用账号数: %d\n", len(cfg.GetEnabledAccounts()))
	fmt.Println()

	// 打印每个账号的详细信息
	for i, acc := range cfg.Accounts {
		fmt.Printf("账号 %d: %s (%s)\n", i+1, acc.Name, acc.ID)
		fmt.Printf("  策略: %s (%s)\n", acc.GetStrategyName(), acc.Strategy)
		fmt.Printf("  提示词: %s (%s)\n", acc.GetPromptTypeName(), acc.PromptType)
		fmt.Printf("  API Key: %s...%s\n", acc.APIKey[:10], acc.APIKey[len(acc.APIKey)-10:])
		fmt.Printf("  API Secret: %s...%s\n", acc.APISecret[:10], acc.APISecret[len(acc.APISecret)-10:])
		fmt.Printf("  启用状态: %v\n", acc.Enabled)
		fmt.Println()
	}

	// 测试获取启用的账号
	fmt.Println("【启用的账号】")
	enabledAccounts := cfg.GetEnabledAccounts()
	for _, acc := range enabledAccounts {
		fmt.Printf("  - %s (%s-%s)\n", acc.Name, acc.GetStrategyName(), acc.GetPromptTypeName())
	}
	fmt.Println()

	// 测试根据ID获取账号
	fmt.Println("【根据ID获取账号】")
	acc := cfg.GetAccountByID("account_1")
	if acc != nil {
		fmt.Printf("  找到账号: %s\n", acc.Name)
		fmt.Printf("  策略: %s\n", acc.GetStrategyName())
		fmt.Printf("  提示词: %s\n", acc.GetPromptTypeName())
	} else {
		fmt.Println("  ❌ 未找到账号")
	}
	fmt.Println()

	// 测试全局配置
	fmt.Println("【全局配置测试】")
	globalCfg := config.Get()
	if globalCfg != nil {
		fmt.Printf("  ✅ 全局配置可用\n")
		fmt.Printf("  账号数量: %d\n", len(globalCfg.Accounts))
	} else {
		fmt.Println("  ❌ 全局配置不可用")
	}
	fmt.Println()

	fmt.Println("=== 测试完成 ===")
}
