/*
Package indicators 通用指标计算函数

主要功能：
- CalculateEMA(klines []binance.Kline, period int) float64                             // 计算EMA
- CalculateMACD(klines []binance.Kline) *MACDData                                      // 计算MACD
- CalculateRSI(klines []binance.Kline, period int) float64                             // 计算RSI
- CalculateBollingerBands(klines []binance.Kline, period int, stdDev float64) *BBData  // 计算布林带
- CalculateATR(klines []binance.Kline, period int) float64                             // 计算ATR
- CalculateADX(klines []binance.Kline, period int) float64                             // 计算ADX
- CalculateStochRSI(klines []binance.Kline, period int) *StochRSIData                  // 计算Stochastic RSI
- CalculateVWAP(klines []binance.Kline) float64                                        // 计算VWAP
- GetVolume(kline binance.Kline) float64                                               // 获取成交量
- formatPrice(value float64) float64                                                   // 格式化价格（2位小数）
- formatMACD(value float64) float64                                                    // 格式化MACD（4位小数）
- formatPercent(value float64) float64                                                 // 格式化百分比（2位小数）
*/
package indicators

import (
	"crypto-ai-trader/binance"
	"math"
	"strconv"

	"github.com/markcheno/go-talib"
)

// CalculateEMA 计算指数移动平均线（使用ta-lib）
// period: EMA周期（如9, 21, 55）
// 返回：最新的EMA值
func CalculateEMA(klines []binance.Kline, period int) float64 {
	if len(klines) < period {
		return 0
	}

	// 提取收盘价
	closes := extractCloses(klines)

	// 使用ta-lib计算EMA
	ema := talib.Ema(closes, period)

	// 返回最新值并格式化
	return formatPrice(ema[len(ema)-1])
}

// CalculateMACD 计算MACD指标（使用ta-lib）
// 使用标准参数：快线12，慢线26，信号线9
// 返回：最新的MACD数据
func CalculateMACD(klines []binance.Kline) *MACDData {
	if len(klines) < 26 {
		return nil
	}

	// 提取收盘价
	closes := extractCloses(klines)

	// 使用ta-lib计算MACD
	macd, signal, histogram := talib.Macd(closes, 12, 26, 9)

	// 获取最新值
	latest := len(macd) - 1

	return &MACDData{
		DIF:       formatMACD(macd[latest]),
		DEA:       formatMACD(signal[latest]),
		Histogram: formatMACD(histogram[latest]),
	}
}

// CalculateRSI 计算RSI指标（使用ta-lib）
// period: RSI周期（通常为14）
// 返回：最新的RSI值（0-100）
func CalculateRSI(klines []binance.Kline, period int) float64 {
	if len(klines) < period+1 {
		return 0
	}

	// 提取收盘价
	closes := extractCloses(klines)

	// 使用ta-lib计算RSI
	rsi := talib.Rsi(closes, period)

	// 返回最新值并格式化
	return formatPercent(rsi[len(rsi)-1])
}

// CalculateBollingerBands 计算布林带（使用ta-lib）
// period: 周期（通常为20）
// stdDev: 标准差倍数（通常为2）
// 返回：最新的布林带数据
func CalculateBollingerBands(klines []binance.Kline, period int, stdDev float64) *BBData {
	if len(klines) < period {
		return nil
	}

	// 提取收盘价
	closes := extractCloses(klines)

	// 使用ta-lib计算布林带
	upper, middle, lower := talib.BBands(closes, period, stdDev, stdDev, talib.SMA)

	// 获取最新值
	latest := len(upper) - 1

	return &BBData{
		Upper:  formatPrice(upper[latest]),
		Middle: formatPrice(middle[latest]),
		Lower:  formatPrice(lower[latest]),
	}
}

// CalculateATR 计算平均真实波幅（使用ta-lib）
// period: ATR周期（通常为14）
// 返回：最新的ATR值
func CalculateATR(klines []binance.Kline, period int) float64 {
	if len(klines) < period+1 {
		return 0
	}

	// 提取高、低、收盘价
	highs, lows, closes := extractHLC(klines)

	// 使用ta-lib计算ATR
	atr := talib.Atr(highs, lows, closes, period)

	// 返回最新值并格式化
	return formatPrice(atr[len(atr)-1])
}

// CalculateADX 计算平均趋向指标（使用ta-lib）
// period: ADX周期（通常为14）
// 返回：最新的ADX值
func CalculateADX(klines []binance.Kline, period int) float64 {
	if len(klines) < period*2 {
		return 0
	}

	// 提取高、低、收盘价
	highs, lows, closes := extractHLC(klines)

	// 使用ta-lib计算ADX
	adx := talib.Adx(highs, lows, closes, period)

	// 返回最新值并格式化
	return formatPercent(adx[len(adx)-1])
}

// CalculateStochRSI 计算Stochastic RSI（使用ta-lib）
// period: 周期（通常为14）
// 返回：最新的Stochastic RSI数据
func CalculateStochRSI(klines []binance.Kline, period int) *StochRSIData {
	if len(klines) < period*2 {
		return nil
	}

	// 提取收盘价
	closes := extractCloses(klines)

	// 使用ta-lib计算Stochastic RSI
	fastK, fastD := talib.StochRsi(closes, period, 5, 3, talib.SMA)

	// 获取最新值
	latest := len(fastK) - 1

	return &StochRSIData{
		K: formatPercent(fastK[latest]),
		D: formatPercent(fastD[latest]),
	}
}

// CalculateVWAP 计算成交量加权平均价
// 返回：最新的VWAP值
func CalculateVWAP(klines []binance.Kline) float64 {
	if len(klines) == 0 {
		return 0
	}

	totalPV := 0.0
	totalVolume := 0.0

	for _, kline := range klines {
		high, _ := strconv.ParseFloat(kline.High, 64)
		low, _ := strconv.ParseFloat(kline.Low, 64)
		close, _ := strconv.ParseFloat(kline.Close, 64)
		volume, _ := strconv.ParseFloat(kline.Volume, 64)

		// 典型价格 = (High + Low + Close) / 3
		typicalPrice := (high + low + close) / 3
		totalPV += typicalPrice * volume
		totalVolume += volume
	}

	if totalVolume == 0 {
		return 0
	}

	return formatPrice(totalPV / totalVolume)
}

// GetVolume 获取K线成交量
func GetVolume(kline binance.Kline) float64 {
	volume, _ := strconv.ParseFloat(kline.Volume, 64)
	return formatPrice(volume)
}

// extractCloses 提取收盘价数组（辅助函数）
func extractCloses(klines []binance.Kline) []float64 {
	closes := make([]float64, len(klines))
	for i, kline := range klines {
		closes[i], _ = strconv.ParseFloat(kline.Close, 64)
	}
	return closes
}

// extractHLC 提取高、低、收盘价数组（辅助函数）
func extractHLC(klines []binance.Kline) ([]float64, []float64, []float64) {
	highs := make([]float64, len(klines))
	lows := make([]float64, len(klines))
	closes := make([]float64, len(klines))

	for i, kline := range klines {
		highs[i], _ = strconv.ParseFloat(kline.High, 64)
		lows[i], _ = strconv.ParseFloat(kline.Low, 64)
		closes[i], _ = strconv.ParseFloat(kline.Close, 64)
	}

	return highs, lows, closes
}

// formatPrice 格式化价格（2位小数）
func formatPrice(value float64) float64 {
	return math.Round(value*100) / 100
}

// formatMACD 格式化MACD值（4位小数）
func formatMACD(value float64) float64 {
	return math.Round(value*10000) / 10000
}

// formatPercent 格式化百分比值（2位小数）
func formatPercent(value float64) float64 {
	return math.Round(value*100) / 100
}

// getLatestValue 获取数组最新值（辅助函数）
func getLatestValue(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return values[len(values)-1]
}
