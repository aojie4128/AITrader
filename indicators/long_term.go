/*
Package indicators 中长线策略指标计算

主要功能：
- CalculateLongTermIndicators(symbol string, klines4h, klines1h, klines15m []binance.Kline) *LongTermIndicators  // 计算中长线策略指标

中长线策略：持仓2-4小时
时间周期：4h（大趋势） → 1h（主分析） → 15m（入场）
*/
package indicators

import (
	"crypto-ai-trader/binance"
	"crypto-ai-trader/utils"
	"time"

	"go.uber.org/zap"
)

// CalculateLongTermIndicators 计算中长线策略指标
// symbol: 交易对（如BTCUSDT）
// klines4h: 4小时K线数据（建议100根以上）
// klines1h: 1小时K线数据（建议100根以上）
// klines15m: 15分钟K线数据（建议100根以上）
// 返回：中长线策略指标数据
func CalculateLongTermIndicators(symbol string, klines4h, klines1h, klines15m []binance.Kline) *LongTermIndicators {
	utils.Debug("计算中长线策略指标",
		zap.String("symbol", symbol),
		zap.Int("4h_klines", len(klines4h)),
		zap.Int("1h_klines", len(klines1h)),
		zap.Int("15m_klines", len(klines15m)),
	)

	// 验证数据充足性
	if len(klines4h) < 55 || len(klines1h) < 55 || len(klines15m) < 55 {
		utils.Error("K线数据不足，无法计算指标",
			zap.Int("4h", len(klines4h)),
			zap.Int("1h", len(klines1h)),
			zap.Int("15m", len(klines15m)),
		)
		return nil
	}

	indicators := &LongTermIndicators{
		Symbol:    symbol,
		Timestamp: time.Now().Unix(),
		Timeframes: &LongTermTimeframes{
			H4:  calculateTimeframeData(klines4h, "4h"),   // 大趋势判断
			H1:  calculateTimeframeData(klines1h, "1h"),   // 主分析周期
			M15: calculateTimeframeData(klines15m, "15m"), // 入场周期
		},
	}

	utils.Info("中长线策略指标计算完成",
		zap.String("symbol", symbol),
		zap.Float64("4h_close", indicators.Timeframes.H4.ClosePrice),
		zap.Float64("1h_close", indicators.Timeframes.H1.ClosePrice),
		zap.Float64("15m_close", indicators.Timeframes.M15.ClosePrice),
	)

	return indicators
}

// CalculateLongTermIndicatorsWithMarket 计算中长线策略指标（包含市场数据）
// symbol: 交易对（如BTCUSDT）
// klines4h: 4小时K线数据（建议100根以上）
// klines1h: 1小时K线数据（建议100根以上）
// klines15m: 15分钟K线数据（建议100根以上）
// client: 币安客户端（用于获取OI和资金费率）
// oiCache: OI缓存（用于计算变化率）
// 返回：中长线策略指标数据（包含OI和资金费率）
func CalculateLongTermIndicatorsWithMarket(symbol string, klines4h, klines1h, klines15m []binance.Kline, client *binance.Client, oiCache *OICache) *LongTermIndicators {
	// 先计算基础指标
	indicators := CalculateLongTermIndicators(symbol, klines4h, klines1h, klines15m)
	if indicators == nil {
		return nil
	}

	// 获取当前价格
	currentPrice := indicators.Timeframes.M15.ClosePrice

	// 计算市场数据
	marketData := CalculateMarketData(client, symbol, currentPrice, oiCache)
	if marketData != nil {
		indicators.MarketData = marketData
	}

	utils.Info("中长线策略指标计算完成（含市场数据）",
		zap.String("symbol", symbol),
		zap.Float64("oi_current", marketData.OICurrent),
		zap.Float64("funding_rate", marketData.FundingRate),
	)

	return indicators
}
