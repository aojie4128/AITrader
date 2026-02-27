/*
Package main 加密货币AI交易系统主程序

主要功能：
- 初始化系统（日志、配置、币安客户端）
- 获取交易对池
- 创建OI缓存管理器
- 启动定时任务（短线5分钟、长线15分钟更新OI）
- 计算指标并输出JSON数据
*/
package main

import (
	"crypto-ai-trader/binance"
	"crypto-ai-trader/config"
	"crypto-ai-trader/indicators"
	"crypto-ai-trader/utils"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 1. 初始化日志
	if err := utils.Init("logs/app.log", "info"); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer utils.Sync()

	utils.Info("=== 加密货币AI交易系统启动 ===")

	// 2. 加载配置
	cfg, err := config.Load("configs/config.yml")
	if err != nil {
		utils.Error("加载配置失败", zap.Error(err))
		os.Exit(1)
	}
	utils.Info("配置加载成功",
		zap.Int("accounts", len(cfg.Accounts)),
		zap.String("futures_url", cfg.Binance.FuturesURL),
	)

	// 3. 获取交易对池
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
		os.Exit(1)
	}
	utils.Info("交易对池构建完成", zap.Int("total", len(symbols)), zap.Strings("symbols", symbols))

	// 4. 创建OI缓存管理器（保存5个历史记录）
	oiCacheManager := utils.NewOICacheManager(5)
	utils.Info("OI缓存管理器创建完成")

	// 5. 为每个账号创建币安客户端
	clients := make(map[string]*binance.Client)
	for _, account := range cfg.GetEnabledAccounts() {
		client := binance.NewClient(
			cfg.Binance.FuturesURL,
			account.APIKey,
			account.APISecret,
			cfg.GetProxyURL(),
		)
		clients[account.ID] = client
		utils.Info("创建币安客户端",
			zap.String("account_id", account.ID),
			zap.String("strategy", account.Strategy),
		)
	}

	// 6. 启动定时任务
	utils.Info("启动定时任务...")
	
	// 短线策略：每5分钟更新一次OI
	shortTermTicker := time.NewTicker(5 * time.Minute)
	defer shortTermTicker.Stop()

	// 长线策略：每15分钟更新一次OI
	longTermTicker := time.NewTicker(15 * time.Minute)
	defer longTermTicker.Stop()

	// 立即执行一次
	utils.Info("执行初始数据采集...")
	for _, account := range cfg.GetEnabledAccounts() {
		client := clients[account.ID]
		if account.Strategy == "short_term" {
			processShortTermStrategy(client, symbols, oiCacheManager, account.ID)
		} else if account.Strategy == "long_term" {
			processLongTermStrategy(client, symbols, oiCacheManager, account.ID)
		}
	}

	// 监听系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 主循环
	utils.Info("系统运行中，按 Ctrl+C 退出...")
	for {
		select {
		case <-shortTermTicker.C:
			utils.Info("=== 短线策略定时任务触发 ===")
			for _, account := range cfg.GetEnabledAccounts() {
				if account.Strategy == "short_term" {
					client := clients[account.ID]
					processShortTermStrategy(client, symbols, oiCacheManager, account.ID)
				}
			}

		case <-longTermTicker.C:
			utils.Info("=== 长线策略定时任务触发 ===")
			for _, account := range cfg.GetEnabledAccounts() {
				if account.Strategy == "long_term" {
					client := clients[account.ID]
					processLongTermStrategy(client, symbols, oiCacheManager, account.ID)
				}
			}

		case sig := <-sigChan:
			utils.Info("收到退出信号", zap.String("signal", sig.String()))
			utils.Info("=== 系统正常退出 ===")
			return
		}
	}
}

// processShortTermStrategy 处理短线策略
func processShortTermStrategy(client *binance.Client, symbols []string, oiCacheManager *utils.OICacheManager, accountID string) {
	utils.Info("处理短线策略", zap.String("account_id", accountID), zap.Int("symbols", len(symbols)))

	for _, symbol := range symbols {
		// 获取K线数据
		klines1h, err := client.GetKlines(symbol, "1h", 100)
		if err != nil {
			utils.Error("获取1h K线失败", zap.String("symbol", symbol), zap.Error(err))
			continue
		}

		klines15m, err := client.GetKlines(symbol, "15m", 100)
		if err != nil {
			utils.Error("获取15m K线失败", zap.String("symbol", symbol), zap.Error(err))
			continue
		}

		klines5m, err := client.GetKlines(symbol, "5m", 100)
		if err != nil {
			utils.Error("获取5m K线失败", zap.String("symbol", symbol), zap.Error(err))
			continue
		}

		// 获取OI缓存
		oiCache := oiCacheManager.Get(symbol)
		if oiCache == nil {
			oiCache = &utils.OICache{
				Symbol:     symbol,
				History:    []float64{},
				Timestamps: []int64{},
			}
		}

		// 转换为indicators.OICache类型
		indicatorOICache := &indicators.OICache{
			Symbol:     oiCache.Symbol,
			History:    oiCache.History,
			Timestamps: oiCache.Timestamps,
		}

		// 计算指标（包含市场数据）
		result := indicators.CalculateShortTermIndicatorsWithMarket(
			symbol,
			klines1h,
			klines15m,
			klines5m,
			client,
			indicatorOICache,
		)

		if result == nil {
			utils.Error("计算短线指标失败", zap.String("symbol", symbol))
			continue
		}

		// 更新OI缓存
		if result.MarketData != nil {
			oiCacheManager.Update(symbol, result.MarketData.OICurrent, time.Now().Unix())
		}

		// 输出JSON（可以发送给AI或保存到文件）
		outputIndicators(result, accountID, "short_term")
	}
}

// processLongTermStrategy 处理长线策略
func processLongTermStrategy(client *binance.Client, symbols []string, oiCacheManager *utils.OICacheManager, accountID string) {
	utils.Info("处理长线策略", zap.String("account_id", accountID), zap.Int("symbols", len(symbols)))

	for _, symbol := range symbols {
		// 获取K线数据
		klines4h, err := client.GetKlines(symbol, "4h", 100)
		if err != nil {
			utils.Error("获取4h K线失败", zap.String("symbol", symbol), zap.Error(err))
			continue
		}

		klines1h, err := client.GetKlines(symbol, "1h", 100)
		if err != nil {
			utils.Error("获取1h K线失败", zap.String("symbol", symbol), zap.Error(err))
			continue
		}

		klines15m, err := client.GetKlines(symbol, "15m", 100)
		if err != nil {
			utils.Error("获取15m K线失败", zap.String("symbol", symbol), zap.Error(err))
			continue
		}

		// 获取OI缓存
		oiCache := oiCacheManager.Get(symbol)
		if oiCache == nil {
			oiCache = &utils.OICache{
				Symbol:     symbol,
				History:    []float64{},
				Timestamps: []int64{},
			}
		}

		// 转换为indicators.OICache类型
		indicatorOICache := &indicators.OICache{
			Symbol:     oiCache.Symbol,
			History:    oiCache.History,
			Timestamps: oiCache.Timestamps,
		}

		// 计算指标（包含市场数据）
		result := indicators.CalculateLongTermIndicatorsWithMarket(
			symbol,
			klines4h,
			klines1h,
			klines15m,
			client,
			indicatorOICache,
		)

		if result == nil {
			utils.Error("计算长线指标失败", zap.String("symbol", symbol))
			continue
		}

		// 更新OI缓存
		if result.MarketData != nil {
			oiCacheManager.Update(symbol, result.MarketData.OICurrent, time.Now().Unix())
		}

		// 输出JSON（可以发送给AI或保存到文件）
		outputIndicators(result, accountID, "long_term")
	}
}

// outputIndicators 输出指标数据（JSON格式）
func outputIndicators(data interface{}, accountID, strategy string) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		utils.Error("序列化JSON失败", zap.Error(err))
		return
	}

	utils.Info("指标数据",
		zap.String("account_id", accountID),
		zap.String("strategy", strategy),
		zap.String("json", string(jsonData)),
	)

	// TODO: 这里可以将JSON数据发送给AI进行分析
	// 或者保存到文件、数据库等
}
