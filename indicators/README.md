# Indicators 指标计算模块

## 功能说明

本模块用于计算技术指标，支持两套策略：

### 短线策略（持仓30-90分钟）
- **时间周期**: 1h（方向过滤） → 15m（主分析） → 5m（入场）
- **目标**: 快速进出，捕捉短期波动

### 中长线策略（持仓2-4小时）
- **时间周期**: 4h（大趋势） → 1h（主分析） → 15m（入场）
- **目标**: 波段操作，跟随趋势

## 实现阶段

### 第一阶段（已实现）- 核心指标

**趋势指标**
- EMA (9, 21, 55) - 判断趋势方向和强度

**动能指标**
- MACD (12, 26, 9) - 动能方向与背离
- RSI (14) - 超买超卖判断

**波动率指标**
- Bollinger Bands (20, 2) - 价格通道
- ATR (14) - 波动率和止损参考

**成交量**
- Volume - 验证突破有效性

### 第二阶段（待实现）- 增强胜率

- ADX - 判断趋势强度
- VWAP - 成交量加权平均价
- OI - 持仓量变化
- Funding Rate - 资金费率

### 第三阶段（待实现）- 机构级

- Ichimoku - 一目均衡表
- CVD - 累积成交量差
- Volume Profile - 成交量分布

## 文件结构

```
indicators/
├── types.go           # 数据结构定义
├── common.go          # 通用指标计算函数
├── short_term.go      # 短线策略（1h → 15m → 5m）
├── long_term.go       # 中长线策略（4h → 1h → 15m）
└── README.md          # 说明文档
```

## 使用方式

### 短线策略

```go
// 获取K线数据
klines1h, _ := client.GetKlines("BTCUSDT", "1h", 100)
klines15m, _ := client.GetKlines("BTCUSDT", "15m", 100)
klines5m, _ := client.GetKlines("BTCUSDT", "5m", 100)

// 计算短线指标
shortTerm := indicators.CalculateShortTermIndicators("BTCUSDT", klines1h, klines15m, klines5m)

// 使用指标
fmt.Printf("1h EMA55: %.2f\n", shortTerm.H1.EMA55)
fmt.Printf("15m RSI: %.2f\n", shortTerm.M15.RSI)
fmt.Printf("5m MACD: %.4f\n", shortTerm.M5.MACD.DIF)
```

### 中长线策略

```go
// 获取K线数据
klines4h, _ := client.GetKlines("BTCUSDT", "4h", 100)
klines1h, _ := client.GetKlines("BTCUSDT", "1h", 100)
klines15m, _ := client.GetKlines("BTCUSDT", "15m", 100)

// 计算中长线指标
longTerm := indicators.CalculateLongTermIndicators("BTCUSDT", klines4h, klines1h, klines15m)

// 使用指标
fmt.Printf("4h EMA55: %.2f\n", longTerm.H4.EMA55)
fmt.Printf("1h RSI: %.2f\n", longTerm.H1.RSI)
fmt.Printf("15m MACD: %.4f\n", longTerm.M15.MACD.DIF)
```

## 数据结构

### ShortTermIndicators
```json
{
  "symbol": "BTCUSDT",
  "timestamp": 1234567890,
  "1h": { /* TimeframeData */ },
  "15m": { /* TimeframeData */ },
  "5m": { /* TimeframeData */ }
}
```

### LongTermIndicators
```json
{
  "symbol": "BTCUSDT",
  "timestamp": 1234567890,
  "4h": { /* TimeframeData */ },
  "1h": { /* TimeframeData */ },
  "15m": { /* TimeframeData */ }
}
```

### TimeframeData
```json
{
  "close_price": 50000.0,
  "high_price": 50100.0,
  "low_price": 49900.0,
  "open_price": 49950.0,
  "ema9": 49980.0,
  "ema21": 49950.0,
  "ema55": 49900.0,
  "macd": {
    "dif": 10.5,
    "dea": 8.3,
    "histogram": 4.4
  },
  "rsi": 65.5,
  "bb": {
    "upper": 50200.0,
    "middle": 50000.0,
    "lower": 49800.0
  },
  "atr": 150.0,
  "volume": 1234.56
}
```

## 测试

```bash
go run test/indicators/test_indicators.go
```

## 设计原则

1. **最小指标集** - 避免指标冗余，降低过拟合风险
2. **多周期确认** - 大周期过滤方向，小周期精确入场
3. **可扩展性** - 预留第二、三阶段扩展字段
4. **核心目标** - 提高胜率，获得盈利

## 注意事项

- 建议每个时间周期至少提供100根K线数据
- 指标计算需要足够的历史数据（至少55根K线）
- 返回的指标值都是最新的（当前K线的指标值）
- 所有价格和指标值都是float64类型
