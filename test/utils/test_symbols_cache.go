/*
交易对池和OI缓存测试程序

测试内容：
- 测试交易对池获取
- 测试OI缓存管理

运行方式：
  go run test/utils/test_symbols_cache.go
*/
package main

import (
	"fmt"
	"time"

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

	utils.Info("=== 交易对池和OI缓存测试开始 ===")

	// 加载配置
	cfg, err := config.Load("configs/config.yml")
	if err != nil {
		utils.Fatal("加载配置失败", zap.Error(err))
	}

	// ========== 1. 测试交易对池 ==========
	fmt.Println("【1. 交易对池测试】")
	
	minScore := cfg.SymbolPool.ExternalSymbols.MinScore
	if minScore == 0 {
		minScore = 75 // 默认75分
	}
	
	symbols, err := utils.GetSymbolPool(
		cfg.SymbolPool.DefaultSymbols,
		cfg.SymbolPool.ExcludeSymbols,
		cfg.SymbolPool.ExternalSymbols.URL,
		cfg.SymbolPool.ExternalSymbols.IsUse,
		minScore,
	)
	
	if err != nil {
		utils.Error("获取交易对池失败", zap.Error(err))
	} else {
		fmt.Printf("  ✓ 交易对总数: %d\n", len(symbols))
		fmt.Println("  前10个交易对:")
		count := 10
		if len(symbols) < count {
			count = len(symbols)
		}
		for i := 0; i < count; i++ {
			fmt.Printf("    %d. %s\n", i+1, symbols[i])
		}
	}
	fmt.Println()

	// ========== 2. 测试OI缓存管理器 ==========
	fmt.Println("【2. OI缓存管理器测试】")
	
	// 创建缓存管理器
	cacheManager := utils.NewOICacheManager(5)
	fmt.Println("  ✓ 缓存管理器创建成功")
	
	// 模拟更新OI数据
	fmt.Println("\n  模拟更新OI数据...")
	testSymbol := "BTCUSDT"
	
	// 添加5个历史数据点
	for i := 0; i < 5; i++ {
		oi := 5300.0 + float64(i)*10.0
		timestamp := time.Now().Unix() - int64(i*300) // 每5分钟一个数据点
		cacheManager.Update(testSymbol, oi, timestamp)
		fmt.Printf("    添加数据点 %d: OI=%.2f M, 时间=%d\n", i+1, oi, timestamp)
	}
	
	// 获取缓存
	fmt.Println("\n  获取缓存数据...")
	cache := cacheManager.Get(testSymbol)
	if cache != nil {
		fmt.Printf("  ✓ 交易对: %s\n", cache.Symbol)
		fmt.Printf("  ✓ 历史记录数: %d\n", len(cache.History))
		fmt.Println("  ✓ 历史OI值:")
		for i, oi := range cache.History {
			fmt.Printf("    %d. %.2f M (时间戳: %d)\n", i+1, oi, cache.Timestamps[i])
		}
	}
	
	// 测试缓存统计
	fmt.Println("\n  缓存统计信息:")
	stats := cacheManager.GetStats()
	fmt.Printf("  ✓ 缓存的交易对数: %v\n", stats["symbol_count"])
	fmt.Printf("  ✓ 总记录数: %v\n", stats["total_records"])
	fmt.Printf("  ✓ 最大记录数: %v\n", stats["max_size"])
	
	// 测试多个交易对
	fmt.Println("\n  添加更多交易对...")
	cacheManager.Update("ETHUSDT", 1200.5, time.Now().Unix())
	cacheManager.Update("BNBUSDT", 450.3, time.Now().Unix())
	
	allSymbols := cacheManager.GetSymbols()
	fmt.Printf("  ✓ 已缓存的交易对: %v\n", allSymbols)
	
	// 测试过期检查
	fmt.Println("\n  测试过期检查...")
	isExpired := cacheManager.IsExpired(testSymbol, 3600) // 1小时
	fmt.Printf("  ✓ %s 是否过期(1小时): %v\n", testSymbol, isExpired)
	
	fmt.Println()
	utils.Info("=== 测试完成 ===")
}
