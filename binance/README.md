# Binance API 封装

币安合约API的Go语言封装。

## 特性

- ✅ 支持代理
- ✅ 自动签名
- ✅ 多账号支持
- ✅ 请求日志记录
- ✅ 错误处理
- ✅ 超时控制

## 使用方法

### 1. 创建客户端

```go
import "crypto-ai-trader/binance"

// 创建客户端
client := binance.NewClient(
    "your_api_key",
    "your_api_secret",
    "https://fapi.binance.com",
    "http://127.0.0.1:7890", // 代理URL，不需要则传空字符串
)
```

### 2. 测试连接

```go
// Ping测试
err := client.Ping()
if err != nil {
    log.Fatal(err)
}
```

### 3. 多账号使用

```go
// 为每个账号创建独立的客户端
clients := make(map[string]*binance.Client)

for _, acc := range accounts {
    client := binance.NewClient(
        acc.APIKey,
        acc.APISecret,
        baseURL,
        proxyURL,
    )
    clients[acc.ID] = client
}
```

## Client 结构

```go
type Client struct {
    apiKey     string           // API密钥
    apiSecret  string           // API密钥
    baseURL    string           // 基础URL
    httpClient *http.Client     // HTTP客户端
}
```

## 主要方法

### Client 方法

**NewClient**
创建新的币安客户端

```go
func NewClient(apiKey, apiSecret, baseURL string, proxyURL string) *Client
```

**SetProxy**
设置代理

```go
func (c *Client) SetProxy(proxyURL string)
```

**Ping**
测试连接

```go
func (c *Client) Ping() error
```

**GetServerTime**
获取服务器时间

```go
func (c *Client) GetServerTime() (int64, error)
```

### Account 方法

**GetAccountInfo**
获取账户信息（包含USDT资产和持仓列表）

```go
func (c *Client) GetAccountInfo() (*AccountInfo, error)
```

**GetBalance**
获取USDT余额

```go
func (c *Client) GetBalance() (*Balance, error)
```

**GetPositions**
获取持仓信息（过滤0持仓）

```go
func (c *Client) GetPositions() ([]Position, error)
```

**GetPositionRisk**
获取持仓风险

```go
func (c *Client) GetPositionRisk(symbol string) ([]PositionRisk, error)
```

## API端点

所有API端点定义在 `endpoints.go` 中：

```go
const (
    EndpointPing         = "/fapi/v1/ping"         // 测试连接
    EndpointServerTime   = "/fapi/v1/time"         // 获取服务器时间
    EndpointAccount      = "/fapi/v2/account"      // 获取账户信息
    EndpointBalance      = "/fapi/v2/balance"      // 获取账户余额
    EndpointPositionRisk = "/fapi/v2/positionRisk" // 获取持仓风险
    EndpointKlines       = "/fapi/v1/klines"       // 获取K线数据
    // ...更多端点
)
```

## 签名机制

所有需要签名的请求会自动：
1. 添加时间戳参数
2. 按字母顺序排序参数
3. 使用HMAC SHA256生成签名
4. 添加签名到请求参数

## 代理支持

支持HTTP代理，适用于需要代理访问币安API的场景：

```go
client.SetProxy("http://127.0.0.1:7890")
```

## 错误处理

所有API调用都会返回详细的错误信息：

```go
err := client.Ping()
if err != nil {
    // 处理错误
    log.Printf("API调用失败: %v", err)
}
```

## 日志记录

客户端会自动记录：
- 客户端创建
- API请求（Debug级别）
- API响应（Debug级别）
- 错误信息（Error级别）

## 测试

```bash
go run test/binance/test_client.go
```

## 后续功能

- [x] 获取账户信息
- [x] 获取账户余额
- [x] 获取持仓信息
- [x] 获取持仓风险
- [x] 获取K线数据
- [ ] 下单功能
- [ ] 撤单功能
- [ ] 获取订单信息

## 注意事项

1. API密钥请妥善保管，不要提交到代码仓库
2. 建议使用只读或限制权限的API密钥进行测试
3. 生产环境建议启用IP白名单
4. 注意API调用频率限制
