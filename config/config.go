/*
Package config 配置管理模块

主要功能：
- Load(configPath string) (*Config, error)           // 加载配置文件
- Get() *Config                                       // 获取全局配置
- (c *Config) Validate() error                        // 验证配置
- (c *Config) GetProxyURL() string                    // 获取代理URL
- (c *Config) GetEnabledAccounts() []Account          // 获取所有启用的账号
- (c *Config) GetAccountByID(id string) *Account      // 根据ID获取账号
*/
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构
type Config struct {
	Proxy          ProxyConfig   `yaml:"proxy"`
	Binance        BinanceConfig `yaml:"binance"`
	AccountsConfig string        `yaml:"accounts_config"`
	Accounts       []Account     `yaml:"-"` // 从单独文件加载
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	IsUse bool   `yaml:"is_use"`
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
}

// BinanceConfig 币安API配置
type BinanceConfig struct {
	FuturesURL string `yaml:"futures_url"`
}

var globalConfig *Config

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	// 读取主配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 加载账号配置（相对于主配置文件的路径）
	if cfg.AccountsConfig != "" {
		// 获取主配置文件所在目录
		configDir := filepath.Dir(configPath)
		accountsPath := filepath.Join(configDir, cfg.AccountsConfig)
		
		accounts, err := LoadAccounts(accountsPath)
		if err != nil {
			return nil, fmt.Errorf("加载账号配置失败: %w", err)
		}
		cfg.Accounts = accounts
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

// Get 获取全局配置
func Get() *Config {
	return globalConfig
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证币安配置
	if c.Binance.FuturesURL == "" {
		return fmt.Errorf("币安合约URL不能为空")
	}

	// 验证账号配置
	if len(c.Accounts) == 0 {
		return fmt.Errorf("至少需要配置一个账号")
	}

	return nil
}

// GetProxyURL 获取代理URL
func (c *Config) GetProxyURL() string {
	if !c.Proxy.IsUse {
		return ""
	}
	return fmt.Sprintf("http://%s:%d", c.Proxy.Host, c.Proxy.Port)
}

// GetEnabledAccounts 获取所有启用的账号
func (c *Config) GetEnabledAccounts() []Account {
	var enabled []Account
	for _, acc := range c.Accounts {
		if acc.Enabled {
			enabled = append(enabled, acc)
		}
	}
	return enabled
}

// GetAccountByID 根据ID获取账号
func (c *Config) GetAccountByID(id string) *Account {
	for i := range c.Accounts {
		if c.Accounts[i].ID == id {
			return &c.Accounts[i]
		}
	}
	return nil
}
