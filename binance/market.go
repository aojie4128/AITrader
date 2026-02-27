/*
Package binance 市场数据相关API

主要功能：
- (c *Client) GetOpenInterest(symbol string) (*OpenInterest, error)                    // 获取持仓量
- (c *Client) GetFundingRateHistory(symbol string, limit int) ([]FundingRate, error)   // 获取资金费率历史
- (c *Client) GetPremiumIndex(symbol string) (*PremiumIndex, error)                    // 获取当前资金费率和标记价格
- CalculateOIChange(current, previous float64) float64                                 // 计算持仓量变化率
*/
package binance

import (
	"encoding/json"
	"fmt"
	"strconv"

	"crypto-ai-trader/utils"

	"go.uber.org/zap"
)

// OpenInterest 持仓量数据
type OpenInterest struct {
	Symbol       string `json:"symbol"`       // 交易对
	OpenInterest string `json:"openInterest"` // 持仓量（张数）
	Time         int64  `json:"time"`         // 时间戳
}

// FundingRate 资金费率数据
type FundingRate struct {
	Symbol      string `json:"symbol"`      // 交易对
	FundingRate string `json:"fundingRate"` // 资金费率
	FundingTime int64  `json:"fundingTime"` // 资金费时间
	Time        int64  `json:"time"`        // 时间戳
}

// PremiumIndex 溢价指数和资金费率
type PremiumIndex struct {
	Symbol          string `json:"symbol"`          // 交易对
	MarkPrice       string `json:"markPrice"`       // 标记价格
	IndexPrice      string `json:"indexPrice"`      // 指数价格
	LastFundingRate string `json:"lastFundingRate"` // 最新资金费率
	NextFundingTime int64  `json:"nextFundingTime"` // 下次资金费时间
	Time            int64  `json:"time"`            // 时间戳
}

// GetOpenInterest 获取持仓量
// symbol: 交易对，如 "BTCUSDT"
func (c *Client) GetOpenInterest(symbol string) (*OpenInterest, error) {
	utils.Debug("获取持仓量", zap.String("symbol", symbol))

	params := map[string]string{
		"symbol": symbol,
	}

	body, err := c.doRequest("GET", EndpointOpenInterest, params, false)
	if err != nil {
		return nil, fmt.Errorf("获取持仓量失败: %w", err)
	}

	var oi OpenInterest
	if err := json.Unmarshal(body, &oi); err != nil {
		return nil, fmt.Errorf("解析持仓量数据失败: %w", err)
	}

	utils.Info("获取持仓量成功",
		zap.String("symbol", symbol),
		zap.String("open_interest", oi.OpenInterest),
	)

	return &oi, nil
}

// GetFundingRateHistory 获取资金费率历史
// symbol: 交易对，如 "BTCUSDT"
// limit: 获取数量，默认100，最大1000
func (c *Client) GetFundingRateHistory(symbol string, limit int) ([]FundingRate, error) {
	utils.Debug("获取资金费率历史",
		zap.String("symbol", symbol),
		zap.Int("limit", limit),
	)

	params := map[string]string{
		"symbol": symbol,
	}

	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}

	body, err := c.doRequest("GET", EndpointFundingRate, params, false)
	if err != nil {
		return nil, fmt.Errorf("获取资金费率历史失败: %w", err)
	}

	var fundingRates []FundingRate
	if err := json.Unmarshal(body, &fundingRates); err != nil {
		return nil, fmt.Errorf("解析资金费率数据失败: %w", err)
	}

	utils.Info("获取资金费率历史成功",
		zap.String("symbol", symbol),
		zap.Int("count", len(fundingRates)),
	)

	return fundingRates, nil
}

// GetPremiumIndex 获取当前资金费率和标记价格
// symbol: 交易对，如 "BTCUSDT"
func (c *Client) GetPremiumIndex(symbol string) (*PremiumIndex, error) {
	utils.Debug("获取溢价指数", zap.String("symbol", symbol))

	params := map[string]string{
		"symbol": symbol,
	}

	body, err := c.doRequest("GET", EndpointPremiumIndex, params, false)
	if err != nil {
		return nil, fmt.Errorf("获取溢价指数失败: %w", err)
	}

	var premium PremiumIndex
	if err := json.Unmarshal(body, &premium); err != nil {
		return nil, fmt.Errorf("解析溢价指数数据失败: %w", err)
	}

	utils.Info("获取溢价指数成功",
		zap.String("symbol", symbol),
		zap.String("mark_price", premium.MarkPrice),
		zap.String("funding_rate", premium.LastFundingRate),
	)

	return &premium, nil
}

// CalculateOIChange 计算持仓量变化率
// current: 当前持仓量
// previous: 之前的持仓量
// 返回：变化率（百分比）
func CalculateOIChange(current, previous float64) float64 {
	if previous == 0 {
		return 0
	}
	return ((current - previous) / previous) * 100
}
