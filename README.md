# go_advanced_task1 使用说明

本项目基于 Go 与 go-ethereum（Geth 库），完成两类任务：
- Task1：区块与交易基础能力
  - 通过 RPC 查询指定区块信息（`QueryBlockInfoAt`）
  - 发送一笔原生 ETH 转账（`transfer`）
- Task2：合约部署与交互能力
  - 部署示例合约 Counter（`DeployCounter`）
  - 调用合约方法（`InteractWithCounter`）：读取 `count`，执行 `increment` 并再读取

代码入口在 `main.go`，根据需要手动取消注释对应函数后运行。

## 目录结构
- `task1.go`：区块查询与 ETH 转账
- `task2.go`：合约部署与交互（使用已生成的 Go 绑定 `counter/`）
- `counter/`：通过 abigen 生成的 Counter 合约 Go 绑定
- `counter.sol`、`counter_sol_Counter.abi`、`counter_sol_Counter.bin`：示例合约与 ABI/字节码
- `architecture-design.md`、`理论分析.md`：架构与理论说明

## 前置条件
- Go（建议 1.20+）
- 一个可用的以太坊 RPC 地址（建议测试网，如 Sepolia）
- 一个测试网私钥（用于发起交易/部署合约）

## 环境变量
程序通过环境变量或根目录下的 `.env` 文件读取配置（`main.go` 会自动加载 `.env`）。

必须：
- `ETH_NODE_URL`：以太坊节点 RPC，如 `https://sepolia.infura.io/v3/<YOUR_KEY>`
- `PRIVATE_KEY`：发起交易的私钥，16 进制，不要带 `0x` 前缀（仅用于测试网）

可选（Task2 交互用）：
- `COUNTER_ADDRESS`：已部署的 Counter 合约地址（0x 开头）

macOS（zsh）示例 `.env`：
```
ETH_NODE_URL=https://sepolia.infura.io/v3/YOUR_PROJECT_ID
PRIVATE_KEY=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
# 部署完成后再填：
# COUNTER_ADDRESS=0xYourDeployedCounterAddress
```

## 安装依赖
项目已包含 `go.mod`/`go.sum`，首次运行会自动拉取依赖。
```
go run .
```
或预先拉取：
```
go mod tidy
```

## 运行方式
在 `main.go` 中按需取消注释对应函数，然后运行：
```
go run .
```
或构建后运行：
```
go build -o app && ./app
```

### Task1-1：查询区块信息（QueryBlockInfoAt）
1. 打开 `main.go`，取消注释：
   ```go
   // QueryBlockInfoAt()
   ```
2. 确保已设置 `ETH_NODE_URL`。
3. 运行：
   ```
   go run .
   ```
4. 终端将输出区块号、哈希、时间戳与交易数量（当前示例查询区块 `5671744`）。

如需查询其他高度，可修改 `task1.go` 中的：
```go
blockNumber := big.NewInt(5671744)
```

### Task1-2：发送 ETH 转账（transfer）
1. 在 `main.go`，取消注释：
   ```go
   // transfer()
   ```
2. 设置：
   - `ETH_NODE_URL` 为目标网络 RPC（建议测试网）
   - `PRIVATE_KEY` 为有测试 ETH 的账户（无 `0x` 前缀）
3. 根据需要修改收款地址/金额（`task1.go`）：
   ```go
   toAddress := common.HexToAddress("0xRecipientAddress...")
   var value int64 = 1000000000000000 // 0.001 ETH（单位 wei）
   ```
4. 运行：
   ```
   go run .
   ```
5. 终端输出交易哈希，可在区块浏览器查询状态（选择与 `ETH_NODE_URL` 对应的网络）。

注意：当前使用 `SuggestGasPrice` 与 EIP-155 签名（`NewEIP155Signer`），在主流公链/测试网上可正常广播。

### Task2-1：部署合约（DeployCounter）
1. 在 `main.go`，取消注释：
   ```go
   // DeployCounter()
   ```
2. 设置 `ETH_NODE_URL`、`PRIVATE_KEY`（测试网账户需有少量测试 ETH）。
3. 运行：
   ```
   go run .
   ```
4. 终端将打印：
   - 合约地址（`Address:`）
   - 部署交易哈希（`Transaction hash:`）
5. 将地址复制到 `.env`：
   ```
   COUNTER_ADDRESS=0x...
   ```

说明：本仓库已包含合约 Go 绑定（`counter/`）。如需重新生成，请使用 abigen（可选）。

### Task2-2：交互合约（InteractWithCounter）
1. 在 `.env` 设置 `COUNTER_ADDRESS` 为上一步部署的地址。
2. 在 `main.go`，取消注释：
   ```go
   // InteractWithCounter()
   ```
3. 运行：
   ```
   go run .
   ```
4. 期望输出：
   - 调用前 `count` 值
   - `increment` 交易哈希
   - 交易上链后再次读取的 `count` 值（应 +1）

## 常见问题（FAQ）
- 报错 `ETH_NODE_URL is not set`：未正确设置 `.env` 或环境变量。
- 报错私钥格式：`PRIVATE_KEY` 需为 64 位 16 进制，不含 `0x`。
- 交易长期 pending：确认测试网是否拥堵、账户是否有余额、Gas 设置是否合理。
- 交互时报错合约地址无效：检查 `COUNTER_ADDRESS` 网络是否与 `ETH_NODE_URL` 一致，地址是否正确。
- 本地 Geth：可将 `ETH_NODE_URL` 指向本地节点 RPC（确保开放 HTTP-RPC 并连到目标网络）。

## 安全提示
- 切勿在公共仓库提交真实私钥。
- 始终在测试网演练，确认逻辑后再切换主网。

---
如需进一步自动化（命令行参数切换任务、网络选择、确认提示等），可在 `main.go` 中加入 flag 解析并按参数调用对应函数。