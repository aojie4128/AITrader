/*
Package binance 账户信息相关API

主要功能：
- (c *Client) GetAccountInfo() (*AccountInfo, error)           // 获取账户信息
- (c *Client) GetBalance() (*Balance, error)                   // 获取USDT余额
- (c *Client) GetPositions() ([]Position, error)               // 获取持仓信息
- (c *Client) GetPositionRisk(symbol string) ([]PositionRisk, error)  // 获取持仓风险
*/
package binance

import (
	"encoding/json"
	"fmt"

	"crypto-ai-trader/utils"

	"go.uber.org/zap"
)

// AccountInfo 账户信息
type AccountInfo struct {
	TotalWalletBalance    string     `json:"totalWalletBalance"`    // 账户总余额
	TotalUnrealizedProfit string     `json:"totalUnrealizedProfit"` // 未实现盈亏
	TotalMarginBalance    string     `json:"totalMarginBalance"`    // 保证金余额
	AvailableBalance      string     `json:"availableBalance"`      // 可用余额
	Asset                 Asset      `json:"-"`                     // USDT资产（从assets中提取）
	Positions             []Position `json:"positions"`             // 持仓列表
}

// accountInfoResponse API响应结构（用于解析）
type accountInfoResponse struct {
	TotalWalletBalance    string     `json:"totalWalletBalance"`
	TotalUnrealizedProfit string     `json:"totalUnrealizedProfit"`
	TotalMarginBalance    string     `json:"totalMarginBalance"`
	AvailableBalance      string     `json:"availableBalance"`
	Assets                []Asset    `json:"assets"`
	Positions             []Position `json:"positions"`
}

// Asset 资产信息（USDT）
type Asset struct {
	Asset                  string `json:"asset"`                  // 资产名称（USDT）
	WalletBalance          string `json:"walletBalance"`          // 钱包余额
	UnrealizedProfit       string `json:"unrealizedProfit"`       // 未实现盈亏
	MarginBalance          string `json:"marginBalance"`          // 保证金余额
	MaintMargin            string `json:"maintMargin"`            // 维持保证金
	InitialMargin          string `json:"initialMargin"`          // 起始保证金
	PositionInitialMargin  string `json:"positionInitialMargin"`  // 持仓起始保证金
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"` // 挂单起始保证金
	MaxWithdrawAmount      string `json:"maxWithdrawAmount"`      // 最大可提现金额
	CrossWalletBalance     string `json:"crossWalletBalance"`     // 全仓账户余额
	CrossUnPnl             string `json:"crossUnPnl"`             // 全仓持仓未实现盈亏
	AvailableBalance       string `json:"availableBalance"`       // 可用余额
}

// Position 持仓信息
type Position struct {
	Symbol           string `json:"symbol"`           // 交易对
	PositionAmt      string `json:"positionAmt"`      // 持仓数量
	EntryPrice       string `json:"entryPrice"`       // 开仓均价
	MarkPrice        string `json:"markPrice"`        // 标记价格
	UnRealizedProfit string `json:"unRealizedProfit"` // 未实现盈亏
	LiquidationPrice string `json:"liquidationPrice"` // 强平价格
	Leverage         string `json:"leverage"`         // 杠杆倍数
	MaxNotionalValue string `json:"maxNotionalValue"` // 最大名义价值
	MarginType       string `json:"marginType"`       // 保证金模式
	IsolatedMargin   string `json:"isolatedMargin"`   // 逐仓保证金
	IsAutoAddMargin  string `json:"isAutoAddMargin"`  // 是否自动追加保证金
	PositionSide     string `json:"positionSide"`     // 持仓方向
	Notional         string `json:"notional"`         // 名义价值
	IsolatedWallet   string `json:"isolatedWallet"`   // 逐仓钱包余额
	UpdateTime       int64  `json:"updateTime"`       // 更新时间
}

// Balance 余额信息（单个资产）
type Balance struct {
	Asset            string `json:"asset"`            // 资产（USDT）
	Balance          string `json:"balance"`          // 余额
	AvailableBalance string `json:"availableBalance"` // 可用余额
	UnrealizedProfit string `json:"unrealizedProfit"` // 未实现盈亏
}

// PositionRisk 持仓风险
type PositionRisk struct {
	Symbol           string `json:"symbol"`           // 交易对
	PositionAmt      string `json:"positionAmt"`      // 持仓数量
	EntryPrice       string `json:"entryPrice"`       // 开仓均价
	MarkPrice        string `json:"markPrice"`        // 标记价格
	UnRealizedProfit string `json:"unRealizedProfit"` // 未实现盈亏
	LiquidationPrice string `json:"liquidationPrice"` // 强平价格
	Leverage         string `json:"leverage"`         // 杠杆倍数
	MaxNotionalValue string `json:"maxNotionalValue"` // 最大名义价值
	MarginType       string `json:"marginType"`       // 保证金模式
	IsolatedMargin   string `json:"isolatedMargin"`   // 逐仓保证金
	IsAutoAddMargin  string `json:"isAutoAddMargin"`  // 是否自动追加保证金
	PositionSide     string `json:"positionSide"`     // 持仓方向
	Notional         string `json:"notional"`         // 名义价值
	IsolatedWallet   string `json:"isolatedWallet"`   // 逐仓钱包余额
	UpdateTime       int64  `json:"updateTime"`       // 更新时间
}

// GetAccountInfo 获取账户信息
func (c *Client) GetAccountInfo() (*AccountInfo, error) {
	utils.Debug("获取账户信息")

	body, err := c.doRequest("GET", EndpointAccount, nil, true)
	if err != nil {
		return nil, fmt.Errorf("获取账户信息失败: %w", err)
	}

	var resp accountInfoResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析账户信息失败: %w", err)
	}

	// 提取USDT资产
	var usdtAsset Asset
	for _, asset := range resp.Assets {
		if asset.Asset == "USDT" {
			usdtAsset = asset
			break
		}
	}

	accountInfo := &AccountInfo{
		TotalWalletBalance:    resp.TotalWalletBalance,
		TotalUnrealizedProfit: resp.TotalUnrealizedProfit,
		TotalMarginBalance:    resp.TotalMarginBalance,
		AvailableBalance:      resp.AvailableBalance,
		Asset:                 usdtAsset,
		Positions:             resp.Positions,
	}

	utils.Info("获取账户信息成功",
		zap.String("total_balance", accountInfo.TotalWalletBalance),
		zap.String("available_balance", accountInfo.AvailableBalance),
		zap.String("unrealized_profit", accountInfo.TotalUnrealizedProfit),
		zap.String("usdt_balance", usdtAsset.WalletBalance),
		zap.Int("positions_count", len(accountInfo.Positions)),
	)

	return accountInfo, nil
}

// GetBalance 获取USDT余额
func (c *Client) GetBalance() (*Balance, error) {
	utils.Debug("获取账户余额")

	body, err := c.doRequest("GET", EndpointBalance, nil, true)
	if err != nil {
		return nil, fmt.Errorf("获取账户余额失败: %w", err)
	}

	var balances []Balance
	if err := json.Unmarshal(body, &balances); err != nil {
		return nil, fmt.Errorf("解析账户余额失败: %w", err)
	}

	// 查找USDT余额
	for _, balance := range balances {
		if balance.Asset == "USDT" {
			utils.Info("获取USDT余额成功",
				zap.String("balance", balance.Balance),
				zap.String("available", balance.AvailableBalance),
			)
			return &balance, nil
		}
	}

	return nil, fmt.Errorf("未找到USDT余额")
}

// GetPositions 获取持仓信息
func (c *Client) GetPositions() ([]Position, error) {
	utils.Debug("获取持仓信息")

	// 通过账户信息获取持仓
	accountInfo, err := c.GetAccountInfo()
	if err != nil {
		return nil, err
	}

	// 过滤出有持仓的交易对
	var positions []Position
	for _, pos := range accountInfo.Positions {
		if pos.PositionAmt != "0" && pos.PositionAmt != "0.0" && pos.PositionAmt != "0.00" && pos.PositionAmt != "0.000" {
			positions = append(positions, pos)
		}
	}

	utils.Info("获取持仓信息成功", zap.Int("count", len(positions)))

	return positions, nil
}

// GetPositionRisk 获取持仓风险
func (c *Client) GetPositionRisk(symbol string) ([]PositionRisk, error) {
	utils.Debug("获取持仓风险", zap.String("symbol", symbol))

	params := make(map[string]string)
	if symbol != "" {
		params["symbol"] = symbol
	}

	body, err := c.doRequest("GET", EndpointPositionRisk, params, true)
	if err != nil {
		return nil, fmt.Errorf("获取持仓风险失败: %w", err)
	}

	var positionRisks []PositionRisk
	if err := json.Unmarshal(body, &positionRisks); err != nil {
		return nil, fmt.Errorf("解析持仓风险失败: %w", err)
	}

	utils.Info("获取持仓风险成功", zap.Int("count", len(positionRisks)))

	return positionRisks, nil
}
