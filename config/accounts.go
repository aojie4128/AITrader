/*
Package config 账号配置管理

主要功能：
- LoadAccounts(accountsPath string) ([]Account, error)  // 加载账号配置文件
- (a *Account) Validate() error                          // 验证账号配置
- (a *Account) GetStrategyName() string                  // 获取策略名称（中文）
- (a *Account) GetPromptTypeName() string                // 获取提示词类型名称（中文）
- (a *Account) GetPromptTypeDescription() string         // 获取提示词类型描述
*/
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Account 账号配置
type Account struct {
	ID         string `yaml:"id"`
	Name       string `yaml:"name"`
	Strategy   string `yaml:"strategy"`    // short_term 或 long_term
	PromptType string `yaml:"prompt_type"` // minimal 或 detailed
	APIKey     string `yaml:"api_key"`
	APISecret  string `yaml:"api_secret"`
	Enabled    bool   `yaml:"enabled"`
}

// AccountsConfig 账号配置文件结构
type AccountsConfig struct {
	Accounts []Account `yaml:"accounts"`
}

// LoadAccounts 加载账号配置文件
func LoadAccounts(accountsPath string) ([]Account, error) {
	// 读取账号配置文件
	data, err := os.ReadFile(accountsPath)
	if err != nil {
		return nil, fmt.Errorf("读取账号配置文件失败: %w", err)
	}

	// 解析YAML
	var accountsCfg AccountsConfig
	if err := yaml.Unmarshal(data, &accountsCfg); err != nil {
		return nil, fmt.Errorf("解析账号配置文件失败: %w", err)
	}

	// 验证账号配置
	for i, acc := range accountsCfg.Accounts {
		if err := acc.Validate(); err != nil {
			return nil, fmt.Errorf("账号[%d]配置无效: %w", i, err)
		}
	}

	return accountsCfg.Accounts, nil
}

// Validate 验证账号配置
func (a *Account) Validate() error {
	if a.ID == "" {
		return fmt.Errorf("账号ID不能为空")
	}
	if a.Name == "" {
		return fmt.Errorf("账号名称不能为空")
	}
	if a.Strategy != "short_term" && a.Strategy != "long_term" {
		return fmt.Errorf("策略类型无效: %s (必须是 short_term 或 long_term)", a.Strategy)
	}
	if a.PromptType != "minimal" && a.PromptType != "detailed" {
		return fmt.Errorf("提示词类型无效: %s (必须是 minimal 或 detailed)", a.PromptType)
	}
	if a.APIKey == "" {
		return fmt.Errorf("API Key不能为空")
	}
	if a.APISecret == "" {
		return fmt.Errorf("API Secret不能为空")
	}
	return nil
}

// GetStrategyName 获取策略名称（中文）
func (a *Account) GetStrategyName() string {
	switch a.Strategy {
	case "short_term":
		return "短线"
	case "long_term":
		return "中长线"
	default:
		return "未知"
	}
}

// GetPromptTypeName 获取提示词类型名称（中文）
func (a *Account) GetPromptTypeName() string {
	switch a.PromptType {
	case "minimal":
		return "简洁版"
	case "detailed":
		return "详细版"
	default:
		return "未知"
	}
}

// GetPromptTypeDescription 获取提示词类型描述
func (a *Account) GetPromptTypeDescription() string {
	switch a.PromptType {
	case "minimal":
		return "只提供数据和输出格式，让AI自主判断"
	case "detailed":
		return "明确写明进场出场条件和交易逻辑"
	default:
		return "未知类型"
	}
}
