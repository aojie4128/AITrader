/*
Package binance K线数据相关API

主要功能：
- (c *Client) GetKlines(symbol, interval string, limit int) ([]Kline, error)  // 获取K线数据
*/
package binance

import (
	"encoding/json"
	"fmt"
	"strconv"

	"crypto-ai-trader/utils"

	"go.uber.org/zap"
)

// Kline K线数据
type Kline struct {
	OpenTime                 int64  `json:"openTime"`                 // 开盘时间
	Open                     string `json:"open"`                     // 开盘价
	High                     string `json:"high"`                     // 最高价
	Low                      string `json:"low"`                      // 最低价
	Close                    string `json:"close"`                    // 收盘价
	Volume                   string `json:"volume"`                   // 成交量
	CloseTime                int64  `json:"closeTime"`                // 收盘时间
	QuoteAssetVolume         string `json:"quoteAssetVolume"`         // 成交额
	NumberOfTrades           int64  `json:"numberOfTrades"`           // 成交笔数
	TakerBuyBaseAssetVolume  string `json:"takerBuyBaseAssetVolume"`  // 主动买入成交量
	TakerBuyQuoteAssetVolume string `json:"takerBuyQuoteAssetVolume"` // 主动买入成交额
}

// GetKlines 获取K线数据
// symbol: 交易对，如 "BTCUSDT"
// interval: K线周期，如 "1m", "5m", "15m", "1h", "4h", "1d"
// limit: 获取数量，默认500，最大1500
func (c *Client) GetKlines(symbol, interval string, limit int) ([]Kline, error) {
	utils.Debug("获取K线数据",
		zap.String("symbol", symbol),
		zap.String("interval", interval),
		zap.Int("limit", limit),
	)

	// 构建参数
	params := map[string]string{
		"symbol":   symbol,
		"interval": interval,
	}

	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}

	// 发送请求
	body, err := c.doRequest("GET", EndpointKlines, params, false)
	if err != nil {
		return nil, fmt.Errorf("获取K线数据失败: %w", err)
	}

	// 解析响应
	// 币安返回的是二维数组格式
	var rawKlines [][]interface{}
	if err := json.Unmarshal(body, &rawKlines); err != nil {
		return nil, fmt.Errorf("解析K线数据失败: %w", err)
	}

	// 转换为结构体
	klines := make([]Kline, 0, len(rawKlines))
	for _, raw := range rawKlines {
		if len(raw) < 11 {
			continue
		}

		kline := Kline{
			OpenTime:                 int64(raw[0].(float64)),
			Open:                     raw[1].(string),
			High:                     raw[2].(string),
			Low:                      raw[3].(string),
			Close:                    raw[4].(string),
			Volume:                   raw[5].(string),
			CloseTime:                int64(raw[6].(float64)),
			QuoteAssetVolume:         raw[7].(string),
			NumberOfTrades:           int64(raw[8].(float64)),
			TakerBuyBaseAssetVolume:  raw[9].(string),
			TakerBuyQuoteAssetVolume: raw[10].(string),
		}

		klines = append(klines, kline)
	}

	utils.Info("获取K线数据成功",
		zap.String("symbol", symbol),
		zap.String("interval", interval),
		zap.Int("count", len(klines)),
	)

	return klines, nil
}
