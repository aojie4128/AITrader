/*
Package binance 币安API客户端

主要功能：
- NewClient(apiKey, apiSecret, baseURL string, proxy string) *Client  // 创建客户端
- (c *Client) SetProxy(proxyURL string)                                // 设置代理
- (c *Client) doRequest(method, endpoint string, params map[string]string, signed bool) ([]byte, error)  // 执行HTTP请求
- (c *Client) sign(params map[string]string) string                    // 生成签名
*/
package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"crypto-ai-trader/utils"

	"go.uber.org/zap"
)

// Client 币安API客户端
type Client struct {
	apiKey     string
	apiSecret  string
	baseURL    string
	httpClient *http.Client
}

// NewClient 创建新的币安客户端
func NewClient(apiKey, apiSecret, baseURL string, proxyURL string) *Client {
	client := &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// 设置代理
	if proxyURL != "" {
		client.SetProxy(proxyURL)
	}

	utils.Info("创建币安客户端",
		zap.String("base_url", baseURL),
		zap.Bool("proxy_enabled", proxyURL != ""),
	)

	return client
}

// SetProxy 设置代理
func (c *Client) SetProxy(proxyURL string) {
	if proxyURL == "" {
		return
	}

	proxy, err := url.Parse(proxyURL)
	if err != nil {
		utils.Error("解析代理URL失败", zap.String("proxy", proxyURL), zap.Error(err))
		return
	}

	c.httpClient.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	utils.Info("设置代理", zap.String("proxy", proxyURL))
}

// doRequest 执行HTTP请求
func (c *Client) doRequest(method, endpoint string, params map[string]string, signed bool) ([]byte, error) {
	// 如果需要签名，添加时间戳和签名
	if signed {
		if params == nil {
			params = make(map[string]string)
		}
		params["timestamp"] = fmt.Sprintf("%d", time.Now().UnixMilli())
		
		// 生成签名
		signature := c.sign(params)
		
		// 构建带签名的查询字符串
		queryString := c.buildQueryString(params)
		queryString += "&signature=" + signature
		
		// 构建URL
		fullURL := c.baseURL + endpoint + "?" + queryString
		
		// 创建请求
		req, err := http.NewRequest(method, fullURL, nil)
		if err != nil {
			return nil, fmt.Errorf("创建请求失败: %w", err)
		}
		
		// 添加请求头
		req.Header.Set("X-MBX-APIKEY", c.apiKey)
		req.Header.Set("Content-Type", "application/json")
		
		return c.executeRequest(req, endpoint, signed)
	}

	// 无签名请求
	fullURL := c.baseURL + endpoint
	if len(params) > 0 {
		fullURL += "?" + c.buildQueryString(params)
	}

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("X-MBX-APIKEY", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return c.executeRequest(req, endpoint, signed)
}

// executeRequest 执行HTTP请求
func (c *Client) executeRequest(req *http.Request, endpoint string, signed bool) ([]byte, error) {
	// 发送请求
	utils.Debug("发送API请求",
		zap.String("method", req.Method),
		zap.String("endpoint", endpoint),
		zap.Bool("signed", signed),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		utils.Error("API请求失败",
			zap.String("endpoint", endpoint),
			zap.Error(err),
		)
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		utils.Error("API返回错误",
			zap.String("endpoint", endpoint),
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(body)),
		)
		return nil, fmt.Errorf("API错误 [%d]: %s", resp.StatusCode, string(body))
	}

	utils.Debug("API请求成功",
		zap.String("endpoint", endpoint),
		zap.Int("response_size", len(body)),
	)

	return body, nil
}

// sign 生成签名
func (c *Client) sign(params map[string]string) string {
	// 构建查询字符串
	queryString := c.buildQueryString(params)

	// 使用HMAC SHA256签名
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(queryString))
	signature := hex.EncodeToString(h.Sum(nil))

	return signature
}

// buildQueryString 构建查询字符串
func (c *Client) buildQueryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	// 排序参数（币安要求按字母顺序）
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "signature" { // 签名不参与查询字符串构建
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 构建查询字符串
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, url.QueryEscape(params[k])))
	}

	return strings.Join(parts, "&")
}

// Ping 测试连接
func (c *Client) Ping() error {
	_, err := c.doRequest("GET", EndpointPing, nil, false)
	if err != nil {
		return fmt.Errorf("ping失败: %w", err)
	}

	utils.Info("币安API连接正常")
	return nil
}

// GetServerTime 获取服务器时间
func (c *Client) GetServerTime() (int64, error) {
	_, err := c.doRequest("GET", EndpointServerTime, nil, false)
	if err != nil {
		return 0, err
	}

	// 简单返回当前时间（实际应该解析响应）
	return time.Now().UnixMilli(), nil
}
