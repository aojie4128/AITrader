/*
Package indicators 短线策略指标计算

主要功能：
- CalculateShortTermIndicators(symbol string, klines1h, klines15m, klines5m []binance.Kline) *ShortTermIndicators  // 计算短线策略指标

短线策略：持仓30-90分钟
时间周期：1h（方向过滤） → 15m（主分析） → 5m（入场）
*/
package indicators

import (
	"crypto-ai-trader/binance"
	"crypto-ai-trader/utils"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// CalculateShortTermIndicators 计算短线策略指标
// symbol: 交易对（如BTCUSDT）
// klines1h: 1小时K线数据（建议100根以上）
// klines15m: 15分钟K线数据（建议100根以上）
// klines5m: 5分钟K线数据（建议100根以上）
// 返回：短线策略指标数据
func CalculateShortTermIndicators(symbol string, klines1h, klines15m, klines5m []binance.Kline) *ShortTermIndicators {
	utils.Debug("计算短线策略指标",
		zap.String("symbol", symbol),
		zap.Int("1h_klines", len(klines1h)),
		zap.Int("15m_klines", len(klines15m)),
		zap.Int("5m_klines", len(klines5m)),
	)

	// 验证数据充足性
	if len(klines1h) < 55 || len(klines15m) < 55 || len(klines5m) < 55 {
		utils.Error("K线数据不足，无法计算指标",
			zap.Int("1h", len(klines1h)),
			zap.Int("15m", len(klines15m)),
			zap.Int("5m", len(klines5m)),
		)
		return nil
	}

	indicators := &ShortTermIndicators{
		Symbol:    symbol,
		Timestamp: time.Now().Unix(),
		Timeframes: &ShortTermTimeframes{
			H1:  calculateTimeframeData(klines1h, "1h"),   // 方向过滤
			M15: calculateTimeframeData(klines15m, "15m"), // 主分析周期
			M5:  calculateTimeframeData(klines5m, "5m"),   // 入场周期
		},
	}

	utils.Info("短线策略指标计算完成",
		zap.String("symbol", symbol),
		zap.Float64("1h_close", indicators.Timeframes.H1.ClosePrice),
		zap.Float64("15m_close", indicators.Timeframes.M15.ClosePrice),
		zap.Float64("5m_close", indicators.Timeframes.M5.ClosePrice),
	)

	return indicators
}

// CalculateShortTermIndicatorsWithMarket 计算短线策略指标（包含市场数据）
// symbol: 交易对（如BTCUSDT）
// klines1h: 1小时K线数据（建议100根以上）
// klines15m: 15分钟K线数据（建议100根以上）
// klines5m: 5分钟K线数据（建议100根以上）
// client: 币安客户端（用于获取OI和资金费率）
// oiCache: OI缓存（用于计算变化率）
// 返回：短线策略指标数据（包含OI和资金费率）
func CalculateShortTermIndicatorsWithMarket(symbol string, klines1h, klines15m, klines5m []binance.Kline, client *binance.Client, oiCache *OICache) *ShortTermIndicators {
	// 先计算基础指标
	indicators := CalculateShortTermIndicators(symbol, klines1h, klines15m, klines5m)
	if indicators == nil {
		return nil
	}

	// 获取当前价格
	currentPrice := indicators.Timeframes.M5.ClosePrice

	// 计算市场数据
	marketData := CalculateMarketData(client, symbol, currentPrice, oiCache)
	if marketData != nil {
		indicators.MarketData = marketData
	}

	utils.Info("短线策略指标计算完成（含市场数据）",
		zap.String("symbol", symbol),
		zap.Float64("oi_current", marketData.OICurrent),
		zap.Float64("funding_rate", marketData.FundingRate),
	)

	return indicators
}

// calculateTimeframeData 计算单个时间周期的指标数据
func calculateTimeframeData(klines []binance.Kline, timeframe string) *TimeframeData {
	if len(klines) == 0 {
		return nil
	}

	latest := len(klines) - 1

	// 获取价格信息（格式化为2位小数）
	closePrice, _ := strconv.ParseFloat(klines[latest].Close, 64)
	highPrice, _ := strconv.ParseFloat(klines[latest].High, 64)
	lowPrice, _ := strconv.ParseFloat(klines[latest].Low, 64)
	openPrice, _ := strconv.ParseFloat(klines[latest].Open, 64)
	volume := GetVolume(klines[latest])

	// 计算趋势指标
	ema9 := CalculateEMA(klines, 9)
	ema21 := CalculateEMA(klines, 21)
	ema55 := CalculateEMA(klines, 55)

	// 计算动能指标
	macd := CalculateMACD(klines)
	rsi := CalculateRSI(klines, 14)

	// 计算波动率指标
	bb := CalculateBollingerBands(klines, 20, 2.0)
	atr := CalculateATR(klines, 14)

	// 第二阶段指标（可选）
	var adx *float64
	var vwap *float64
	var stochRSI *StochRSIData
	if len(klines) >= 28 {
		adxValue := CalculateADX(klines, 14)
		if adxValue > 0 {
			adx = &adxValue
		}
		vwapValue := CalculateVWAP(klines)
		if vwapValue > 0 {
			vwap = &vwapValue
		}
		stochRSI = CalculateStochRSI(klines, 14)
	}

	data := &TimeframeData{
		ClosePrice: formatPrice(closePrice),
		HighPrice:  formatPrice(highPrice),
		LowPrice:   formatPrice(lowPrice),
		OpenPrice:  formatPrice(openPrice),
		EMA9:       ema9,
		EMA21:      ema21,
		EMA55:      ema55,
		MACD:       macd,
		RSI:        rsi,
		BB:         bb,
		ATR:        atr,
		Volume:     volume,
		ADX:        adx,
		VWAP:       vwap,
		StochRSI:   stochRSI,
	}

	utils.Debug("时间周期指标计算完成",
		zap.String("timeframe", timeframe),
		zap.Float64("close", data.ClosePrice),
		zap.Float64("ema9", data.EMA9),
		zap.Float64("ema21", data.EMA21),
		zap.Float64("ema55", data.EMA55),
		zap.Float64("rsi", rsi),
	)

	return data
}
