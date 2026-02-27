# 合约多周期交易策略规范 v2（扩展指标版）

## 1. 设计原则

* 保持核心指标最小化
* 高维数据作为过滤层
* 支持 AI 决策模型输入
* 避免指标功能重复

策略结构：

```
Market Regime → Trend → Setup → Trigger → Risk
```

---

# 2. 时间框架设计

| 周期    | 角色    | 主要用途          |
| ----- | ----- | ------------- |
| Daily | 宏观趋势  | 判断市场阶段（趋势/震荡） |
| 4H    | 大趋势   | 决定多空方向        |
| 1H    | 主交易周期 | 寻找波段机会        |
| 15M   | 结构周期  | 判断回调与突破       |
| 5M    | 入场周期  | 精确触发          |

---

# 3. 核心指标层（交易决策必须）

## 3.1 趋势

* EMA (9 / 21 / 55)
* Ichimoku Cloud（趋势过滤）

用途：

* 判断趋势方向
* 判断趋势强度
* 过滤震荡行情

---

## 3.2 动能

* MACD
* RSI
* ADX

用途：

* MACD → 动能方向与背离
* RSI → 强弱判断
* ADX → 判断是否存在趋势

---

## 3.3 波动率

* Bollinger Bands
* ATR

用途：

* 判断波动环境
* 设置止损止盈
* 判断是否进入扩张行情

---

# 4. 触发指标层（入场优化）

仅在小周期使用（15M / 5M）

* Stochastic RSI
* Parabolic SAR

用途：

* 精确入场点
* 捕捉回调结束

---

# 5. 结构分析层

## 5.1 价格结构

* Swing High / Low
* 支撑阻力位
* VWAP
* Fibonacci
* 24h / 7d Range

用途：

* 识别关键交易区域
* 计算风险收益比
* 判断突破有效性

---

# 6. 资金流指标（趋势确认）

## 6.1 必选

* Volume
* CVD
* OI（变化率）
* Funding Rate

## 6.2 可选增强

* MFI
* Taker Buy/Sell Ratio
* Volume Profile

用途：

* 判断趋势是否有资金支持
* 判断逼空 / 多杀多风险

---

# 7. 宏观环境层（是否允许交易）

仅用于风控过滤，不参与具体信号

* BTC 市场趋势
* DXY 美元指数
* SPX 指数
* 市场波动率（如 VIX）
* 交易时段（亚洲 / 欧洲 / 美盘）
* 市场温度指标

用途：

```
IF 宏观风险高 → 降低仓位或停止交易
```

---

# 8. 指标分层总结

## 核心决策层（必须）

* EMA
* Ichimoku
* MACD
* RSI
* ADX
* Bollinger
* ATR

---

## 入场触发层

* Stochastic RSI
* Parabolic SAR

---

## 结构层

* VWAP
* Swing
* S/R
* Fibonacci
* Range

---

## 资金流层

* Volume
* CVD
* OI
* Funding Rate

---

## 宏观过滤层

* BTC Trend
* DXY
* SPX
* VIX
* Session

---

# 9. AI模型输入字段（建议）

```json
{
  "trend_score": 0,
  "momentum_score": 0,
  "volatility_state": "low | normal | high",
  "liquidity_state": "weak | neutral | strong",
  "market_regime": "trend | range | breakout",
  "macro_risk": "low | medium | high"
}
```

---

# 10. 实现优先级（强烈建议）

## 第一阶段（必须实现）

* EMA
* MACD
* RSI
* Bollinger
* ATR
* Volume

## 第二阶段（增强胜率）

* ADX
* VWAP
* OI
* Funding Rate

## 第三阶段（机构级）

* Ichimoku
* CVD
* Volume Profile
* 宏观数据

---

# 11. 为什么不建议一次性全加

* 指标数量 ↑ → 信号冲突 ↑
* 过拟合风险 ↑
* 延迟 ↑
* 系统复杂度 ↑

目标不是指标最多，而是：

```
信息增益最大
```

---

# 12. v2 策略定位

✔ AI辅助交易
✔ 多因子模型
✔ 可扩展架构
✔ 适用于自动化交易系统

---

Version: 2.0
Type: Multi-Factor Trading Spec
