/*
Package indicators 指标数据结构定义

数据结构：
- ShortTermIndicators   // 短线策略指标（1h → 15m → 5m）
- LongTermIndicators    // 中长线策略指标（4h → 1h → 15m）
- TimeframeData         // 单个时间周期的指标数据
- MACDData              // MACD指标数据
- BBData                // 布林带数据
*/
package indicators

// ShortTermIndicators 短线策略指标（持仓30-90分钟）
// 时间周期：1h（方向过滤） → 15m（主分析） → 5m（入场）
type ShortTermIndicators struct {
	Symbol     string              `json:"symbol"`
	Timestamp  int64               `json:"timestamp"`
	MarketData *MarketData         `json:"market_data,omitempty"` // 市场数据（OI、资金费率）
	Timeframes *ShortTermTimeframes `json:"timeframes"`            // 各时间周期指标
}

// LongTermIndicators 中长线策略指标（持仓2-4小时）
// 时间周期：4h（大趋势） → 1h（主分析） → 15m（入场）
type LongTermIndicators struct {
	Symbol     string             `json:"symbol"`
	Timestamp  int64              `json:"timestamp"`
	MarketData *MarketData        `json:"market_data,omitempty"` // 市场数据（OI、资金费率）
	Timeframes *LongTermTimeframes `json:"timeframes"`            // 各时间周期指标
}

// ShortTermTimeframes 短线策略各时间周期
type ShortTermTimeframes struct {
	H1  *TimeframeData `json:"1h"`  // 1小时 - 方向过滤
	M15 *TimeframeData `json:"15m"` // 15分钟 - 主分析周期
	M5  *TimeframeData `json:"5m"`  // 5分钟 - 入场周期
}

// LongTermTimeframes 中长线策略各时间周期
type LongTermTimeframes struct {
	H4  *TimeframeData `json:"4h"`  // 4小时 - 大趋势判断
	H1  *TimeframeData `json:"1h"`  // 1小时 - 主分析周期
	M15 *TimeframeData `json:"15m"` // 15分钟 - 入场周期
}

// MarketData 市场数据（symbol级别）
type MarketData struct {
	// 持仓量数据
	OICurrent  float64   `json:"oi_current"`            // 当前持仓量（百万美元）
	OIHistory  []float64 `json:"oi_history,omitempty"`  // 历史持仓量（最近5个，从新到旧）
	OIChange5m *float64  `json:"oi_change_5m,omitempty"` // 5分钟变化率(%)
	OIChange15m *float64 `json:"oi_change_15m,omitempty"` // 15分钟变化率(%)
	OIChange25m *float64 `json:"oi_change_25m,omitempty"` // 25分钟变化率(%)
	OIChange45m *float64 `json:"oi_change_45m,omitempty"` // 45分钟变化率(%)
	OIChange75m *float64 `json:"oi_change_75m,omitempty"` // 75分钟变化率(%)
	
	// 资金费率数据
	FundingRate float64 `json:"funding_rate"` // 当前资金费率(%)
	FundingAvg3 float64 `json:"funding_avg_3"` // 最近3次平均(%)
}

// TimeframeData 单个时间周期的指标数据（第一阶段：核心指标）
type TimeframeData struct {
	// 价格信息
	ClosePrice float64 `json:"close_price"` // 收盘价
	HighPrice  float64 `json:"high_price"`  // 最高价
	LowPrice   float64 `json:"low_price"`   // 最低价
	OpenPrice  float64 `json:"open_price"`  // 开盘价

	// 趋势指标
	EMA9  float64 `json:"ema9"`  // 9周期指数移动平均线
	EMA21 float64 `json:"ema21"` // 21周期指数移动平均线
	EMA55 float64 `json:"ema55"` // 55周期指数移动平均线

	// 动能指标
	MACD *MACDData `json:"macd"` // MACD指标
	RSI  float64   `json:"rsi"`  // RSI指标(14)

	// 波动率指标
	BB  *BBData `json:"bb"`  // 布林带(20, 2)
	ATR float64 `json:"atr"` // 平均真实波幅(14)

	// 成交量
	Volume float64 `json:"volume"` // 当前成交量

	// 第二阶段扩展（预留）
	ADX      *float64      `json:"adx,omitempty"`       // 平均趋向指标
	VWAP     *float64      `json:"vwap,omitempty"`      // 成交量加权平均价
	StochRSI *StochRSIData `json:"stoch_rsi,omitempty"` // Stochastic RSI

	// 第三阶段扩展（预留）
	Ichimoku *IchimokuData `json:"ichimoku,omitempty"` // 一目均衡表
	CVD      *float64      `json:"cvd,omitempty"`      // 累积成交量差
}

// MACDData MACD指标数据
type MACDData struct {
	DIF       float64 `json:"dif"`       // 差离值（快线-慢线）
	DEA       float64 `json:"dea"`       // 信号线（DIF的EMA）
	Histogram float64 `json:"histogram"` // 柱状图（DIF-DEA）
}

// BBData 布林带数据
type BBData struct {
	Upper  float64 `json:"upper"`  // 上轨
	Middle float64 `json:"middle"` // 中轨（MA）
	Lower  float64 `json:"lower"`  // 下轨
}

// StochRSIData Stochastic RSI数据（第二阶段）
type StochRSIData struct {
	K float64 `json:"k"` // K值
	D float64 `json:"d"` // D值
}

// IchimokuData 一目均衡表数据（第三阶段）
type IchimokuData struct {
	TenkanSen   float64 `json:"tenkan_sen"`   // 转换线
	KijunSen    float64 `json:"kijun_sen"`    // 基准线
	SenkouSpanA float64 `json:"senkou_span_a"` // 先行带A
	SenkouSpanB float64 `json:"senkou_span_b"` // 先行带B
	ChikouSpan  float64 `json:"chikou_span"`  // 迟行线
}
