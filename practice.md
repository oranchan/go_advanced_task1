# 本地运行 Geth 与同步指引（含私有链搭建）

适用环境：macOS，默认 shell 为 zsh。

## 1. 安装 Geth

- 使用 Homebrew（推荐）
```zsh
brew update
brew install ethereum   # 安装 geth 可执行文件
geth version
```
- 从源码编译（本仓库已包含 go-ethereum 源码）
```zsh
cd go-ethereum
make geth
./build/bin/geth version
```

## 2. 启动与进入内置控制台

- 直接进入控制台（不建议长期运行，仅演示）
```zsh
geth console
```
- 附加到已运行节点
```zsh
geth attach               # 通过默认 IPC
# 或指定 HTTP：geth attach http://127.0.0.1:8545
```
- 常用命令（在控制台里执行）
```javascript
eth.blockNumber
eth.syncing
net.peerCount
```

## 3. 同步到主网/测试网（需要共识客户端）
以太坊已切换 PoS，执行层(EL=geth)必须与共识层(CL=Beacon 客户端)配对运行。

- 生成 JWT 秘钥（EL/CL 共享）
```zsh
openssl rand -hex 32 | tr -d "\n" > ./jwt.hex
```
- 启动 Geth（以 Sepolia 测试网为例）
```zsh
geth \
  --sepolia \
  --syncmode snap \
  --datadir ~/ethereum-sepolia \
  --http --http.addr 127.0.0.1 --http.port 8545 \
  --http.api eth,net,web3,txpool \
  --authrpc.addr 127.0.0.1 --authrpc.port 8551 \
  --authrpc.jwtsecret ./jwt.hex
```
- 启动共识客户端（示例：Lighthouse Beacon）
```zsh
# 安装略，可用 brew install lighthouse 或参阅官方文档
lighthouse bn \
  --network sepolia \
  --execution-endpoint http://127.0.0.1:8551 \
  --jwt-secrets ./jwt.hex \
  --checkpoint-sync-url https://sepolia.checkpoint.sync.provider  # 可选，加速首次同步
```
- 验证同步
```zsh
geth attach --exec 'eth.syncing'
geth attach --exec 'eth.blockNumber'
```
提示：若长时间没有共识更新，检查 CL 是否运行、网络是否一致、JWT 是否匹配、时间是否同步（NTP）。

## 4. 使用公共 RPC（无需本地同步）
如果只需开发/调用，可直接使用公共 RPC（Infura/Alchemy/自建网关）：
- 在本项目 `.env` 设置 `ETH_NODE_URL`，运行 `go run .` 使用 `task1.go/task2.go` 中示例代码与链交互。

## 5. 私有链搭建（Clique PoA）
无需共识客户端，适合本地多节点开发与演示。

### 5.1 准备账户与目录
```zsh
mkdir -p ~/.eth-private/node1 ~/.eth-private/node2
geth account new --datadir ~/.eth-private/node1   # 记录地址 A1
geth account new --datadir ~/.eth-private/node2   # 记录地址 A2
```

### 5.2 生成创世文件（genesis.json）
将 A1、A2 替换为你的地址（0x 开头，不带校验签名）。`extraData` 由 32 字节 vanity + 签名者列表 + 65 字节尾部零组成。推荐使用 puppeth 交互式生成；这里给出一个模板（需将中间的签名者地址替换为 A1、A2）：
```json
{
  "config": {
    "chainId": 2025,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0,
    "clique": { "period": 2, "epoch": 30000 }
  },
  "alloc": {
    "A1": { "balance": "0x3635C9ADC5DEA00000" },
    "A2": { "balance": "0x3635C9ADC5DEA00000" }
  },
  "coinbase": "0x0000000000000000000000000000000000000000",
  "difficulty": "0x1",
  "gasLimit": "0x47b760",
  "extradata": "0x0000000000000000000000000000000000000000000000000000000000000000A1A2............................................................................................................................................................."
}
```
> 注意：`extradata` 的拼装容易出错，建议用 `puppeth` 生成 Clique 配置：
```zsh
puppeth   # 按提示选择 Clique、加入签名者地址、生成 genesis.json
```

### 5.3 初始化并启动节点
```zsh
# 初始化
geth --datadir ~/.eth-private/node1 init genesis.json
geth --datadir ~/.eth-private/node2 init genesis.json

# 启动节点 1（作为引导节点与签名者）
geth \
  --datadir ~/.eth-private/node1 \
  --networkid 2025 \
  --port 30303 \
  --http --http.addr 127.0.0.1 --http.port 8546 \
  --http.api eth,net,web3,txpool,clique \
  --allow-insecure-unlock \
  --unlock A1 --password /path/to/password.txt \
  --mine

# 获取节点 1 的 enode（控制台）
geth attach --exec 'admin.nodeInfo.enode' --datadir ~/.eth-private/node1

# 启动节点 2，连接节点 1
geth \
  --datadir ~/.eth-private/node2 \
  --networkid 2025 \
  --port 30304 \
  --http --http.addr 127.0.0.1 --http.port 8547 \
  --http.api eth,net,web3,txpool,clique \
  --bootnodes "enode://<node1_enode>@127.0.0.1:30303" \
  --allow-insecure-unlock \
  --unlock A2 --password /path/to/password.txt \
  --mine
```
验证：
```zsh
geth attach http://127.0.0.1:8546 --exec 'eth.blockNumber'
geth attach http://127.0.0.1:8547 --exec 'eth.blockNumber'
```
两个节点的高度应持续增长（Clique 签名出块）。

### 5.4 发送交易（本地私链）
- 在控制台解锁账户并发送：
```javascript
personal.sendTransaction({from: A1, to: A2, value: web3.toWei(0.01, "ether")}, "<password>")
```
- 或复用本仓库 Go 代码：将 `ETH_NODE_URL` 指向本地 HTTP（如 `http://127.0.0.1:8546`），在 `main.go` 取消注释 `transfer()` 或合约演示。

## 6. 常见问题
- `eth.blockNumber` 一直为 0：
  - 公网：确认共识客户端正常运行、与 geth 网络一致、JWT 文件一致、peer 数量充足。
  - 私链：确认已 `--mine` 且签名者账户已解锁、至少有一个对等节点。
- 控制台 `miner.start()` 报错：PoS 下不再支持；在 Clique 私链使用 `--mine` 启动出块。
- 安全：HTTP 仅监听本机（127.0.0.1），不要对公网暴露；私钥仅用于测试环境。

---
如需将节点作为系统服务运行或导出快照/检查点，请说明目标环境再补充脚本。