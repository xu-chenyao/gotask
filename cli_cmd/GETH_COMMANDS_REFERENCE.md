# Geth 命令完全参考手册

## 概述

go-ethereum (geth) 是以太坊的官方命令行客户端。本文档详细介绍了所有可用命令的功能、用法和示例。

main 代码解析：
命令行参数 (os.Args)
    ↓
CLI解析 (app.Run)
    ↓
Before钩子 (maxprocs.Set, debug.Setup)
    ↓
geth()函数
    ↓
prepare() - 网络识别和缓存优化
    ↓
makeFullNode() - 创建节点实例
    ├── gethConfig{Eth, Node, Ethstats, Metrics}
    ├── node.New(&cfg.Node)
    └── 注册账户管理器后端
    ↓
startNode() - 启动节点服务
    ├── utils.StartNode() - 启动核心服务
    │   ├── P2P网络服务器
    │   ├── RPC服务器 (HTTP/WS/IPC)
    │   ├── 数据库服务
    │   └── 以太坊协议服务
    ├── 钱包事件处理
    │   ├── WalletArrived → 自动打开
    │   ├── WalletOpened → 账户派生
    │   └── WalletDropped → 资源清理
    └── 同步监控 (可选)
        └── DoneEvent → 自动退出
    ↓
stack.Wait() - 等待节点运行
    ↓
After钩子 (debug.Exit, prompt.Stdin.Close)

## 命令分类

### 1. 基本信息命令

#### `geth version`
**功能**: 显示geth版本信息和构建详情
**用法**: `geth version`
**示例输出**:
```
Geth
Version: 1.13.0-stable
Git Commit: 3c063cdec
Architecture: amd64
Go Version: go1.21.1
Operating System: linux
```

#### `geth license`
**功能**: 显示软件许可证信息 (GNU GPL v3)
**用法**: `geth license`

#### `geth help`
**功能**: 显示帮助信息和可用命令列表
**用法**: `geth help` 或 `geth --help`

### 2. 初始化和配置命令

#### `geth init`
**功能**: 使用创世块配置初始化新的区块链
**用法**: `geth init [flags] <genesis.json>`
**参数**:
- `--datadir`: 指定数据目录
**示例**:
```bash
geth --datadir ./mychain init genesis.json
```

#### `geth dumpconfig`
**功能**: 导出当前配置到TOML格式文件
**用法**: `geth dumpconfig > config.toml`
**说明**: 可以编辑此配置文件后用 `--config` 参数加载

#### `geth dumpgenesis`
**功能**: 导出当前网络的创世块配置
**用法**: `geth --datadir /path/to/data dumpgenesis`
**输出**: JSON格式的创世块配置

### 3. 账户管理命令

#### `geth account`
**功能**: 账户管理主命令，包含多个子命令

##### `geth account list`
**功能**: 列出所有本地账户地址
**用法**: `geth account list`
**示例输出**:
```
Account #0: {7df9a875a174b3bc565e6424a0050ebc1b2d1d82} keystore:///path/to/keystore/UTC--2023...
Account #1: {f466859ead1932d743d622cb74fc058882e8648a} keystore:///path/to/keystore/UTC--2023...
```

##### `geth account new`
**功能**: 创建新的以太坊账户
**用法**: `geth account new`
**交互**: 会提示输入密码
**示例**:
```bash
geth account new
# 或使用密码文件
geth account new --password /path/to/password.txt
```

##### `geth account update`
**功能**: 更新现有账户的密码
**用法**: `geth account update <address>`

##### `geth account import`
**功能**: 从私钥文件导入账户
**用法**: `geth account import <keyfile>`

### 4. 钱包管理命令

#### `geth wallet`
**功能**: 硬件钱包管理

##### `geth wallet status`
**功能**: 显示硬件钱包状态
##### `geth wallet open`
**功能**: 打开硬件钱包会话
##### `geth wallet derive`
**功能**: 从硬件钱包派生账户

### 5. 数据库管理命令

#### `geth db`
**功能**: 底层数据库操作主命令

##### `geth db stats`
**功能**: 显示数据库统计信息
**用法**: `geth db stats`
**输出**: 数据库大小、键数量等统计信息

##### `geth db compact`
**功能**: 压缩数据库以减少存储空间
**用法**: `geth db compact`

##### `geth db get`
**功能**: 从数据库获取原始数据
**用法**: `geth db get <hex-key>`

##### `geth db put`
**功能**: 向数据库写入原始数据
**用法**: `geth db put <hex-key> <hex-value>`

##### `geth db delete`
**功能**: 从数据库删除键值对
**用法**: `geth db delete <hex-key>`

##### `geth db check`
**功能**: 检查数据库完整性
**用法**: `geth db check`

### 6. 区块链数据操作命令

#### `geth import`
**功能**: 从文件导入区块链数据
**用法**: `geth import <blockchain-file>`
**支持格式**: RLP编码的区块链数据
**示例**:
```bash
geth --datadir ./mychain import blockchain.rlp
```

#### `geth export`
**功能**: 导出区块链数据到文件
**用法**: `geth export [first-block] [last-block] <filename>`
**示例**:
```bash
# 导出全部区块
geth --datadir ./mychain export blockchain.rlp
# 导出指定范围
geth --datadir ./mychain export 0 100 first_100_blocks.rlp
```

#### `geth dump`
**功能**: 转储指定区块的世界状态
**用法**: `geth dump [block-number-or-hash]`
**输出**: JSON格式的账户状态数据

#### `geth removedb`
**功能**: 删除区块链数据库
**用法**: `geth removedb`
**警告**: 此操作不可逆，会删除所有区块链数据

### 7. 快照管理命令

#### `geth snapshot`
**功能**: 状态快照管理主命令

##### `geth snapshot prune-state`
**功能**: 基于快照修剪历史状态数据
**用法**: `geth snapshot prune-state <state-root>`
**说明**: 删除不属于指定状态根的历史数据

##### `geth snapshot verify-state`
**功能**: 基于快照重新计算状态根进行验证
**用法**: `geth snapshot verify-state <state-root>`

##### `geth snapshot check-dangling-storage`
**功能**: 检查是否有悬空的快照存储数据
**用法**: `geth snapshot check-dangling-storage <state-root>`

##### `geth snapshot inspect-account`
**功能**: 检查特定账户在所有快照层中的信息
**用法**: `geth snapshot inspect-account <address>`

##### `geth snapshot traverse-state`
**功能**: 遍历状态树并执行快速验证
**用法**: `geth snapshot traverse-state <state-root>`

##### `geth snapshot dump`
**功能**: 使用快照作为后端转储区块状态
**用法**: `geth snapshot dump [block-hash-or-number]`

### 8. 控制台和脚本命令

#### `geth console`
**功能**: 启动交互式JavaScript控制台
**用法**: `geth console`
**功能特性**:
- 内置web3.js API
- 账户管理
- 智能合约交互
- 区块链查询

**示例命令**:
```javascript
> eth.accounts
> eth.getBalance(eth.accounts[0])
> personal.newAccount("password")
> eth.sendTransaction({from: eth.accounts[0], to: "0x...", value: web3.toWei(1, "ether")})
```

#### `geth attach`
**功能**: 连接到正在运行的geth实例
**用法**: `geth attach [ipc-endpoint]`
**连接方式**:
- IPC: `geth attach /path/to/geth.ipc`
- HTTP: `geth attach http://localhost:8545`
- WebSocket: `geth attach ws://localhost:8546`

#### `geth js` (已废弃)
**功能**: 执行JavaScript脚本文件
**替代方案**: 使用 `geth console` 或 `geth attach`

### 9. 网络和节点运行

#### 启动完整节点
**用法**: `geth` (不带参数)
**说明**: 启动完整的以太坊节点，同步主网数据

#### 主要网络选项
- `--mainnet`: 连接以太坊主网 (默认)
- `--goerli`: 连接Goerli测试网
- `--sepolia`: 连接Sepolia测试网
- `--dev`: 开发者模式，自动挖矿的私有网络

#### RPC服务选项
- `--http`: 启用HTTP-RPC服务器
- `--http.addr`: HTTP-RPC服务器地址 (默认: localhost)
- `--http.port`: HTTP-RPC服务器端口 (默认: 8545)
- `--ws`: 启用WebSocket-RPC服务器
- `--graphql`: 启用GraphQL服务器

### 10. 挖矿相关命令

#### 挖矿选项
- `--mine`: 启用挖矿
- `--miner.threads`: 挖矿线程数
- `--miner.etherbase`: 挖矿收益地址
- `--miner.gasprice`: 最小gas价格

**示例**:
```bash
geth --mine --miner.threads=4 --miner.etherbase=0x123...
```

### 11. 版本和安全命令

#### `geth version-check`
**功能**: 检查当前版本是否存在已知安全漏洞
**用法**: `geth version-check`
**参数**:
- `--check.url`: 自定义漏洞数据源URL
- `--check.version`: 指定要检查的版本

### 12. 实验性功能命令

#### `geth verkle`
**功能**: Verkle树相关操作 (实验性)

##### `geth verkle verify`
**功能**: 验证MPT到Verkle树的转换
**用法**: `geth verkle verify <state-root>`

##### `geth verkle dump`
**功能**: 转储Verkle树到DOT文件
**用法**: `geth verkle dump <state-root> <key1> [key2...]`

### 13. 历史数据管理

#### `geth import-history`
**功能**: 导入历史数据文件
**用法**: `geth import-history <history-file>`

#### `geth export-history`
**功能**: 导出历史数据
**用法**: `geth export-history [first] [last] <filename>`

#### `geth prune-history`
**功能**: 修剪历史数据以节省存储空间
**用法**: `geth prune-history`

#### `geth download-era`
**功能**: 下载era历史数据文件
**用法**: `geth download-era <era-number>`

## 常用参数和标志

### 数据目录
- `--datadir`: 指定数据目录路径
- `--keystore`: 指定keystore目录路径

### 网络配置
- `--networkid`: 指定网络ID
- `--port`: P2P网络监听端口
- `--maxpeers`: 最大对等节点数

### 日志和调试
- `--verbosity`: 日志级别 (0-5)
- `--log.json`: JSON格式日志
- `--pprof`: 启用pprof HTTP服务器
- `--trace`: 启用执行跟踪

### 性能调优
- `--cache`: 分配给内部缓存的内存 (MB)
- `--cache.database`: 数据库缓存大小
- `--cache.trie`: Trie缓存大小
- `--cache.gc`: GC期间保留的Trie节点百分比

## 配置文件示例

### 创世块配置 (genesis.json)
```json
{
  "config": {
    "chainId": 1337,
    "homesteadBlock": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0
  },
  "alloc": {
    "0x7df9a875a174b3bc565e6424a0050ebc1b2d1d82": {
      "balance": "300000"
    }
  },
  "coinbase": "0x0000000000000000000000000000000000000000",
  "difficulty": "0x20000",
  "extraData": "",
  "gasLimit": "0x2fefd8",
  "nonce": "0x0000000000000042",
  "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "timestamp": "0x00"
}
```

### TOML配置文件示例 (config.toml)
```toml
[Eth]
NetworkId = 1337
SyncMode = "full"
DatabaseCache = 512
TrieCleanCache = 256
TrieDirtyCache = 256
TrieTimeout = "100ms"
EnablePreimageRecording = false

[Node]
DataDir = "./data"
KeyStoreDir = "./keystore"
HTTPHost = "localhost"
HTTPPort = 8545
HTTPCors = ["*"]
HTTPVirtualHosts = ["localhost"]

[Node.P2P]
MaxPeers = 50
NoDiscovery = false
BootstrapNodes = []
ListenAddr = ":30303"
```

## 常用工作流程

### 1. 设置私有测试网络
```bash
# 1. 创建创世块配置
# 2. 初始化网络
geth --datadir ./private-network init genesis.json

# 3. 创建账户
geth --datadir ./private-network account new

# 4. 启动网络
geth --datadir ./private-network --networkid 1337 --http --http.api personal,eth,net,web3 --mine --miner.etherbase 0x123... console
```

### 2. 连接到测试网
```bash
# 连接到Goerli测试网
geth --goerli --http --http.api eth,net,web3 --sync-mode light

# 连接后使用控制台
geth --goerli attach
```

### 3. 数据库维护
```bash
# 查看数据库统计
geth --datadir ./mychain db stats

# 压缩数据库
geth --datadir ./mychain db compact

# 检查数据库完整性
geth --datadir ./mychain db check
```

### 4. 备份和恢复
```bash
# 导出区块链数据
geth --datadir ./mychain export backup.rlp

# 在新环境导入数据
geth --datadir ./new-chain init genesis.json
geth --datadir ./new-chain import backup.rlp
```

## 故障排除

### 常见问题
1. **数据库损坏**: 使用 `geth db check` 检查，必要时重新同步
2. **磁盘空间不足**: 使用 `geth snapshot prune-state` 清理历史状态
3. **同步慢**: 考虑使用 `--sync-mode fast` 或 `--sync-mode snap`
4. **连接问题**: 检查防火墙设置和 `--port` 参数

### 常见问题
1. **geth命令不存在**: 
   ```bash
   # 安装 geth
   # macOS: brew install ethereum
   # Ubuntu: apt-get install ethereum
   # 或从源码编译: make geth
   ```

2. **端口被占用**:
   ```bash
   # 使用不同端口
   geth --http --http.port 8546
   ```

3. **磁盘空间不足**:
   ```bash
   # 使用轻量级同步模式
   geth --syncmode light
   ```

4. **数据库损坏**:
   ```bash
   # 检查数据库
   geth --datadir ./mychain db check
   
   # 如果损坏，重新同步
   geth removedb
   ```

## 安全提示

⚠️ **重要安全提醒**:
- 生产环境中不要使用 `--allow-insecure-unlock`
- 定期备份 keystore 目录
- 使用强密码保护账户
- 限制 RPC 接口的网络访问
- 定期更新 geth 版本

### 日志级别说明
- 0: Critical
- 1: Error  
- 2: Warning
- 3: Info (默认)
- 4: Debug
- 5: Trace

## 安全建议

1. **密码管理**: 使用强密码，安全存储密码文件
2. **网络访问**: 限制RPC端口访问，使用防火墙
3. **私钥保护**: 定期备份keystore，使用硬件钱包
4. **版本更新**: 定期检查和更新geth版本
5. **监控**: 启用日志记录和监控

## 性能优化

1. **缓存设置**: 根据可用内存调整 `--cache` 参数
2. **数据库调优**: 适当设置数据库缓存大小
3. **网络优化**: 调整 `--maxpeers` 和连接参数
4. **存储**: 使用SSD存储提高I/O性能

### 内存和缓存设置
```bash
# 增加缓存大小（适用于大内存服务器）
geth --cache 2048 --http

# 限制内存使用（适用于资源受限环境）
geth --cache 256 --http
```

### 网络优化
```bash
# 增加对等节点数量
geth --maxpeers 100 --http

# 指定引导节点
geth --bootnodes "enode://..." --http
```

## 日志和调试

### 启用详细日志
```bash
# 设置日志级别
geth --verbosity 4 --http

# 输出JSON格式日志
geth --log.json --http
```

### 性能分析
```bash
# 启用pprof服务器
geth --pprof --http

# 访问性能分析界面
# http://localhost:6060/debug/pprof/
```
## 参考资源

- [官方文档](https://geth.ethereum.org/docs/)
- [GitHub仓库](https://github.com/ethereum/go-ethereum)
- [以太坊黄皮书](https://ethereum.github.io/yellowpaper/paper.pdf)
- [JSON-RPC API文档](https://ethereum.org/en/developers/docs/apis/json-rpc/)

---

**更新时间**: 2024年
**版本**: 适用于 geth v1.13+
**维护者**: AI Assistant 