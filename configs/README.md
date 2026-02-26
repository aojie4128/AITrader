# 配置文件说明

## 文件结构

```
configs/
├── config.yml              # 主配置文件（可提交到git）
├── accounts.yml            # 账号配置文件（不提交到git，包含敏感信息）
└── accounts.example.yml    # 账号配置示例（可提交到git）
```

## 首次使用

1. 复制示例文件：
```bash
cp configs/accounts.example.yml configs/accounts.yml
```

2. 编辑 `accounts.yml`，填入真实的API密钥

3. 确保 `accounts.yml` 已添加到 `.gitignore`

## 配置说明

### config.yml - 主配置

```yaml
# 代理配置
proxy:
  is_use: true              # 是否使用代理
  host: 127.0.0.1           # 代理主机
  port: 12334               # 代理端口

# 币安API配置
binance:
  futures_url: https://fapi.binance.com  # 币安合约API地址

# 账号配置文件路径（相对于config.yml的路径）
accounts_config: "accounts.yml"
```

### accounts.yml - 账号配置

```yaml
accounts:
  - id: "account_1"                    # 账号唯一标识
    name: "短线-简洁版"                # 账号名称
    strategy: "short_term"             # 策略类型：short_term 或 long_term
    prompt_type: "minimal"             # 提示词类型：minimal 或 detailed
    api_key: "YOUR_API_KEY"            # 币安API Key
    api_secret: "YOUR_API_SECRET"      # 币安API Secret
    enabled: true                      # 是否启用
```

## 策略类型

- **short_term** (短线)：快速进出，适合短期交易
- **long_term** (中长线)：趋势跟踪，适合中长期持仓

## 提示词类型

- **minimal** (简洁版)：
  - 只提供K线数据和技术指标
  - 只规定输出格式要求
  - 让AI根据数据自主判断交易信号
  - 适合测试AI的自主决策能力
  
- **detailed** (详细版)：
  - 明确写明进场条件（如：RSI<30且MACD金叉）
  - 明确写明出场条件（如：止损3%，止盈5%）
  - 提供详细的交易逻辑和规则
  - 适合执行明确的交易策略

## 账号组合

系统支持4个账号，建议配置：
1. 短线-简洁版
2. 短线-详细版
3. 中长线-简洁版
4. 中长线-详细版

## 安全提示

⚠️ **重要**：
- `accounts.yml` 包含敏感信息，切勿提交到代码仓库
- 定期更换API密钥
- 使用只读或限制权限的API密钥进行测试
- 生产环境建议使用环境变量管理密钥
