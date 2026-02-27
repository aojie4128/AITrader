/*
Package binance API端点常量

币安合约API端点定义
*/
package binance

const (
	// 基础端点
	EndpointPing       = "/fapi/v1/ping" // 测试连接
	EndpointServerTime = "/fapi/v1/time" // 获取服务器时间
	
	// 账户端点
	EndpointAccount      = "/fapi/v2/account"      // 获取账户信息
	EndpointBalance      = "/fapi/v2/balance"      // 获取账户余额
	EndpointPositionRisk = "/fapi/v2/positionRisk" // 获取持仓风险
	
	// 市场数据端点
	EndpointKlines = "/fapi/v1/klines" // 获取K线数据
)
