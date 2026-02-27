/*
Package indicators 市场数据指标计算

主要功能：
- CalculateOIMetrics(client *binance.Client, symbol string, currentPrice float64) *OIMetrics  // 计算持仓量指标
- CalculateFundingMetrics(client *binance.Client, symbol string) *FundingMetrics              // 计算资金费率指标
*/
package indicators

import (
	"crypto-ai-trader/binance"
	"crypto-ai-trader/utils"
	"strconv"

	"go.uber.org/zap"
)

// OIMetrics 持仓量指标
type OIMetrics struct {
	Current float64 // 当前持仓量（USDT价值）
}

// FundingMetrics 资金费率指标
type FundingMetrics struct {
	Current float64 // 当前资金费率(%)
	Avg3    float64 // 最近3次平均(%)
}

// OICache 持仓量缓存（用于计算变化率）
type OICache struct {
	Symbol    string    // 交易对
	History   []float64 // 历史OI值（从新到旧，最多5个）
	Timestamps []int64  // 对应的时间戳
}

// CalculateMarketData 计算市场数据（OI + 资金费率）
// client: 币安客户端
// symbol: 交易对
// currentPrice: 当前价格
// oiCache: OI缓存（可选，用于计算变化率）
// 返回：市场数据
func CalculateMarketData(client *binance.Client, symbol string, currentPrice float64, oiCache *OICache) *MarketData {
	// 获取当前OI
	oiMetrics := CalculateOIMetrics(client, symbol, currentPrice)
	if oiMetrics == nil {
		return nil
	}

	// 获取资金费率
	fundingMetrics := CalculateFundingMetrics(client, symbol)
	if fundingMetrics == nil {
		return nil
	}

	marketData := &MarketData{
		OICurrent:   formatPrice(oiMetrics.Current / 1000000), // 转换为百万美元
		FundingRate: fundingMetrics.Current,
		FundingAvg3: fundingMetrics.Avg3,
	}

	// 如果有缓存，计算OI变化率
	if oiCache != nil && len(oiCache.History) > 0 {
		marketData.OIHistory = oiCache.History
		
		// 计算不同时间段的变化率
		if len(oiCache.History) >= 2 {
			change5m := calculateOIChangeRate(oiMetrics.Current/1000000, oiCache.History[0])
			marketData.OIChange5m = &change5m
		}
		if len(oiCache.History) >= 4 {
			change15m := calculateOIChangeRate(oiMetrics.Current/1000000, oiCache.History[2])
			marketData.OIChange15m = &change15m
		}
		if len(oiCache.History) >= 5 {
			change25m := calculateOIChangeRate(oiMetrics.Current/1000000, oiCache.History[4])
			marketData.OIChange25m = &change25m
			change45m := calculateOIChangeRate(oiMetrics.Current/1000000, oiCache.History[4])
			marketData.OIChange45m = &change45m
		}
	}

	return marketData
}

// CalculateOIMetrics 计算持仓量指标
// client: 币安客户端
// symbol: 交易对
// currentPrice: 当前价格（用于计算USDT价值）
// 返回：持仓量指标数据
func CalculateOIMetrics(client *binance.Client, symbol string, currentPrice float64) *OIMetrics {
	// 获取当前持仓量
	oi, err := client.GetOpenInterest(symbol)
	if err != nil {
		utils.Error("获取持仓量失败", zap.Error(err))
		return nil
	}

	// 解析持仓量（张数）
	oiValue, err := strconv.ParseFloat(oi.OpenInterest, 64)
	if err != nil {
		utils.Error("解析持仓量失败", zap.Error(err))
		return nil
	}

	// 计算USDT价值（持仓量 * 当前价格）
	currentOIValue := oiValue * currentPrice

	metrics := &OIMetrics{
		Current: formatPrice(currentOIValue),
	}

	utils.Debug("持仓量指标计算完成",
		zap.String("symbol", symbol),
		zap.Float64("oi_value", metrics.Current),
	)

	return metrics
}

// CalculateFundingMetrics 计算资金费率指标
// client: 币安客户端
// symbol: 交易对
// 返回：资金费率指标数据
func CalculateFundingMetrics(client *binance.Client, symbol string) *FundingMetrics {
	// 获取当前资金费率
	premium, err := client.GetPremiumIndex(symbol)
	if err != nil {
		utils.Error("获取当前资金费率失败", zap.Error(err))
		return nil
	}

	currentRate, err := strconv.ParseFloat(premium.LastFundingRate, 64)
	if err != nil {
		utils.Error("解析当前资金费率失败", zap.Error(err))
		return nil
	}

	// 获取最近3次资金费率历史
	fundingRates, err := client.GetFundingRateHistory(symbol, 3)
	if err != nil {
		utils.Error("获取资金费率历史失败", zap.Error(err))
		return &FundingMetrics{
			Current: formatPercent(currentRate * 100),
			Avg3:    0,
		}
	}

	// 计算最近3次平均
	sum := 0.0
	count := 0
	for _, fr := range fundingRates {
		rate, err := strconv.ParseFloat(fr.FundingRate, 64)
		if err == nil {
			sum += rate
			count++
		}
	}

	avg3 := 0.0
	if count > 0 {
		avg3 = sum / float64(count)
	}

	metrics := &FundingMetrics{
		Current: formatPercent(currentRate * 100),
		Avg3:    formatPercent(avg3 * 100),
	}

	utils.Debug("资金费率指标计算完成",
		zap.String("symbol", symbol),
		zap.Float64("current", metrics.Current),
		zap.Float64("avg3", metrics.Avg3),
	)

	return metrics
}

// CalculateOIChangeWithHistory 计算持仓量变化率（需要历史数据）
// currentOI: 当前持仓量
// historicalOI: 历史持仓量数据（按时间倒序）
// interval: 时间间隔（1h, 4h, 24h）
// 返回：变化率(%)
func CalculateOIChangeWithHistory(currentOI float64, historicalOI []float64, interval string) float64 {
	if len(historicalOI) == 0 {
		return 0
	}

	var previousOI float64
	switch interval {
	case "1h":
		if len(historicalOI) >= 1 {
			previousOI = historicalOI[0]
		}
	case "4h":
		if len(historicalOI) >= 4 {
			previousOI = historicalOI[3]
		}
	case "24h":
		if len(historicalOI) >= 24 {
			previousOI = historicalOI[23]
		}
	default:
		return 0
	}

	if previousOI == 0 {
		return 0
	}

	change := ((currentOI - previousOI) / previousOI) * 100
	return formatPercent(change)
}

// ShouldTradeBasedOnFunding 根据资金费率判断是否适合交易
// fundingRate: 当前资金费率(%)
// direction: 交易方向（"long" 或 "short"）
// 返回：是否适合交易，原因
func ShouldTradeBasedOnFunding(fundingRate float64, direction string) (bool, string) {
	// 资金费率阈值
	const (
		extremeHigh = 0.1  // 极高阈值
		high        = 0.05 // 偏高阈值
		low         = -0.05 // 偏低阈值
		extremeLow  = -0.1  // 极低阈值
	)

	if direction == "long" {
		// 做多时，资金费率过高不适合
		if fundingRate > extremeHigh {
			return false, "资金费率极高，市场过度做多，不适合开多"
		}
		if fundingRate > high {
			return false, "资金费率偏高，多头拥挤，建议谨慎"
		}
		return true, "资金费率正常，可以做多"
	}

	if direction == "short" {
		// 做空时，资金费率过低不适合
		if fundingRate < extremeLow {
			return false, "资金费率极低，市场过度做空，不适合开空"
		}
		if fundingRate < low {
			return false, "资金费率偏低，空头拥挤，建议谨慎"
		}
		return true, "资金费率正常，可以做空"
	}

	return false, "未知交易方向"
}

// AnalyzeOIAndPrice 分析持仓量和价格的关系
// priceChange: 价格变化率(%)
// oiChange: 持仓量变化率(%)
// 返回：市场状态描述
func AnalyzeOIAndPrice(priceChange, oiChange float64) string {
	if priceChange > 0 && oiChange > 0 {
		return "价格上涨+OI增加：真实多头趋势，新资金进场"
	}
	if priceChange > 0 && oiChange < 0 {
		return "价格上涨+OI减少：空头平仓推动，可能反转"
	}
	if priceChange < 0 && oiChange > 0 {
		return "价格下跌+OI增加：真实空头趋势，新空单进场"
	}
	if priceChange < 0 && oiChange < 0 {
		return "价格下跌+OI减少：多头平仓推动，可能反转"
	}
	return "价格和OI无明显变化，市场震荡"
}

// calculateOIChangeRate 计算OI变化率
func calculateOIChangeRate(current, previous float64) float64 {
	if previous == 0 {
		return 0
	}
	change := ((current - previous) / previous) * 100
	return formatPercent(change)
}

// UpdateOICache 更新OI缓存
// cache: 现有缓存
// newOI: 新的OI值（百万美元）
// timestamp: 时间戳
// maxSize: 最大缓存数量（建议5个）
// 返回：更新后的缓存
func UpdateOICache(cache *OICache, newOI float64, timestamp int64, maxSize int) *OICache {
	if cache == nil {
		cache = &OICache{
			History:    []float64{},
			Timestamps: []int64{},
		}
	}

	// 添加新值到开头
	cache.History = append([]float64{newOI}, cache.History...)
	cache.Timestamps = append([]int64{timestamp}, cache.Timestamps...)

	// 保持最大数量
	if len(cache.History) > maxSize {
		cache.History = cache.History[:maxSize]
		cache.Timestamps = cache.Timestamps[:maxSize]
	}

	return cache
}
