/*
Package utils 交易对管理

主要功能：
- GetSymbolPoolFromConfig(cfg *config.Config) ([]string, error)  // 从配置获取交易对池
- fetchExternalSymbols(url string) ([]string, error)             // 从外部API获取交易对
*/
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// ExternalSymbolsResponse 外部API响应结构
type ExternalSymbolsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		TopCoins    []CoinInfo `json:"top_coins"`
		BottomCoins []CoinInfo `json:"bottom_coins"`
	} `json:"data"`
}

// CoinInfo 币种信息
type CoinInfo struct {
	Pair  string  `json:"pair"`  // 交易对，如 "BTCUSDT"
	Score float64 `json:"score"` // 评分
}

// GetSymbolPool 获取交易对池（简化版本，避免循环依赖）
// defaultSymbols: 默认交易对列表
// excludeSymbols: 排除的交易对列表
// externalURL: 外部API地址
// useExternal: 是否使用外部API
// minScore: 最低评分要求
// 返回：交易对列表
func GetSymbolPool(defaultSymbols, excludeSymbols []string, externalURL string, useExternal bool, minScore float64) ([]string, error) {
	symbolMap := make(map[string]bool)
	
	// 1. 添加默认交易对
	for _, symbol := range defaultSymbols {
		symbolMap[symbol] = true
		Debug("添加默认交易对", zap.String("symbol", symbol))
	}
	
	// 2. 从外部API获取交易对
	if useExternal && externalURL != "" {
		externalSymbols, err := fetchExternalSymbols(externalURL, minScore)
		if err != nil {
			Warn("获取外部交易对失败", zap.Error(err))
		} else {
			for _, symbol := range externalSymbols {
				symbolMap[symbol] = true
			}
			Info("从外部API获取交易对", 
				zap.Int("count", len(externalSymbols)),
				zap.Float64("min_score", minScore),
			)
		}
	}
	
	// 3. 移除排除的交易对
	for _, symbol := range excludeSymbols {
		delete(symbolMap, symbol)
		Debug("排除交易对", zap.String("symbol", symbol))
	}
	
	// 转换为列表
	symbols := make([]string, 0, len(symbolMap))
	for symbol := range symbolMap {
		symbols = append(symbols, symbol)
	}
	
	Info("交易对池构建完成", zap.Int("total", len(symbols)))
	return symbols, nil
}

// fetchExternalSymbols 从外部API获取交易对
// url: 外部API地址
// minScore: 最低评分要求
func fetchExternalSymbols(url string, minScore float64) ([]string, error) {
	Debug("请求外部API", 
		zap.String("url", url),
		zap.Float64("min_score", minScore),
	)
	
	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// 发送请求
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	
	// 解析JSON
	var response ExternalSymbolsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w", err)
	}
	
	if !response.Success {
		return nil, fmt.Errorf("API返回失败")
	}
	
	// 提取交易对（只获取评分大于minScore的）
	symbolMap := make(map[string]bool)
	filteredCount := 0
	
	// 从top_coins提取
	for _, coin := range response.Data.TopCoins {
		if coin.Pair != "" {
			if coin.Score >= minScore {
				symbolMap[coin.Pair] = true
			} else {
				filteredCount++
				Debug("过滤低评分币种", 
					zap.String("pair", coin.Pair),
					zap.Float64("score", coin.Score),
				)
			}
		}
	}
	
	// 从bottom_coins提取
	for _, coin := range response.Data.BottomCoins {
		if coin.Pair != "" {
			if coin.Score >= minScore {
				symbolMap[coin.Pair] = true
			} else {
				filteredCount++
				Debug("过滤低评分币种", 
					zap.String("pair", coin.Pair),
					zap.Float64("score", coin.Score),
				)
			}
		}
	}
	
	// 转换为列表
	symbols := make([]string, 0, len(symbolMap))
	for symbol := range symbolMap {
		symbols = append(symbols, symbol)
	}
	
	Info("外部API返回交易对", 
		zap.Int("total", len(symbols)),
		zap.Int("filtered", filteredCount),
		zap.Float64("min_score", minScore),
	)
	return symbols, nil
}
