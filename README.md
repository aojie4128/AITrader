# 加密货币AI交易系统

一个基于Go的加密货币AI交易系统，支持多账号、多策略并行交易。

## 项目结构

```
crypto-ai-trader/
├── config/              # 配置管理
├── binance/             # 币安API封装
├── indicators/          # 技术指标计算
├── aggregator/          # 数据聚合器
├── ai/                  # AI分析
├── executor/            # 交易执行器
├── scheduler/           # 调度器
├── trading/             # 交易相关
├── database/            # 数据库
├── notification/        # 通知服务
├── server/              # HTTP服务器
├── utils/               # 公共工具
├── test/                # 测试程序
│   ├── config/          # config模块测试
│   ├── binance/         # binance模块测试
│   └── ...              # 其他模块测试
├── web/                 # React前端
├── prompts/             # AI提示词
├── config.yml           # 配置文件
└── main.go              # 主程序
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置文件

编辑 `config.yml` 文件，配置你的API密钥。

### 3. 运行测试

```bash
# 测试配置模块
cd test/config
go run test_config.go
```

## 已完成模块

- ✅ **config** - 配置管理模块
  - 读取 YAML 配置
  - 代理配置
  - 币安API配置
  - 配置验证

## 开发中模块

- ⏳ binance - 币安API封装
- ⏳ indicators - 技术指标计算
- ⏳ aggregator - 数据聚合器
- ⏳ ai - AI分析
- ⏳ executor - 交易执行器
- ⏳ scheduler - 调度器

## 测试规则

- 每个功能模块在 `test/` 下都有对应的文件夹
- 测试文件命名格式：`test_xx.go`
- 例如：`test/config/test_config.go` 测试 `config/` 模块

## 许可证

MIT
