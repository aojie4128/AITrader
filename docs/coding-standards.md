# 编码规范

## 文件头注释规范

每个包含函数的 Go 文件都应该在开头添加多行注释，列出所有导出的函数及其作用。

### 格式

```go
/*
Package 包名 包描述

主要功能：
- FunctionName1(params) returnType  // 函数1的作用
- FunctionName2(params) returnType  // 函数2的作用
- (r *Receiver) Method() returnType // 方法的作用
*/
package packagename
```

### 示例

```go
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
```

### 测试文件注释

测试文件应该说明测试内容和运行方式：

```go
/*
配置模块测试程序

测试内容：
- 加载主配置文件
- 验证配置解析
- 测试各项功能

运行方式：
  go run test/config/test_config.go
*/
package main
```

## 函数注释规范

每个导出的函数都应该有注释说明其用途：

```go
// Load 加载配置文件
func Load(configPath string) (*Config, error) {
    // ...
}

// GetProxyURL 获取代理URL
func (c *Config) GetProxyURL() string {
    // ...
}
```

## 代码组织规范

### 目录结构

```
module/
├── module.go           # 主要功能
├── helper.go           # 辅助功能
└── types.go            # 类型定义

test/
└── module/
    └── test_module.go  # 测试程序
```

### 文件命名

- 功能模块：`modulename.go`
- 测试文件：`test_modulename.go`
- 类型定义：`types.go`
- 辅助函数：`helper.go` 或 `utils.go`

## 配置文件规范

### 主配置文件

- 位置：`configs/config.yml`
- 内容：非敏感配置
- 版本控制：提交到 git

### 敏感配置文件

- 位置：`configs/accounts.yml`
- 内容：API密钥等敏感信息
- 版本控制：不提交到 git（添加到 .gitignore）
- 提供示例：`configs/accounts.example.yml`

## 注释语言

- 代码注释：中文
- 变量名/函数名：英文（驼峰命名）
- 配置文件注释：中文

## 错误处理

```go
// 返回详细的错误信息
if err != nil {
    return nil, fmt.Errorf("操作失败: %w", err)
}
```

## 日志输出

```go
// 使用中文日志
fmt.Println("✅ 配置加载成功")
fmt.Printf("❌ 加载失败: %v\n", err)
```
