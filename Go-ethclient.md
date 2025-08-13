# 将合约生成GO代码
使用 nodejs，安装 solc 工具：

```go
npm install -g solc
```

使用命令，编译合约代码，会在当目录下生成一个编译好的二进制字节码文件 store_sol_Store.bin：

```go
solcjs --bin Store.sol
```

使用命令，生成合约 abi 文件，会在当目录下生成 store_sol_Store.abi 文件：

```go
solcjs --abi Store.sol
```

abigin 工具可以使用下面的命令安装：

```go
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

使用 abigen 工具根据这两个生成 bin 文件和 abi 文件，生成 go 代码：

```go
abigen --bin=Store_sol_Store.bin --abi=Store_sol_Store.abi --pkg=store --out=store.go
```

# 部署合约
```go
    auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
    if err != nil {
        log.Fatal(err)
    }
    auth.Nonce = big.NewInt(int64(nonce))
    auth.Value = big.NewInt(0)     // in wei
    auth.GasLimit = uint64(300000) // in units
    auth.GasPrice = gasPrice

    input := "1.0"
    address, tx, instance, err := store.DeployStore(auth, client, input)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(address.Hex())
    fmt.Println(tx.Hash().Hex())
```

# **Go-Ethereum 核心功能摘要**

本文档是对所提供 Go 代码中 go-ethereum 库用法的精简总结，旨在快速掌握与以太坊区块链交互的关键函数和模式。

## **1\. 客户端连接**

与以太坊网络的所有交互都始于客户端连接。

* **RPC 连接 (HTTP/HTTPS)**：用于标准的请求-响应式调用，如查询余额、发送交易等。  
  Go  
  client, err := ethclient.Dial("https://...")

* **实时订阅连接 (WebSocket/WSS)**：用于建立持久连接以接收实时事件推送，如新区块或合约日志。  
  Go  
  wsClient, err := ethclient.Dial("wss://...")

## **2\. 链上数据查询（只读操作）**

这些函数用于从区块链读取数据，通常是免费且即时返回的。

* **查询账户余额**：client.BalanceAt 获取指定地址在特定区块（nil 代表最新）的 ETH 余额。  
  Go  
  balance, \_ := client.BalanceAt(ctx, address, nil)

* **查询区块信息**：通过区块号 (client.BlockByNumber) 或哈希 (client.BlockByHash) 获取完整的区块数据。  
  Go  
  block, err := client.BlockByNumber(ctx, big.NewInt(12345))

* **查询交易信息**：通过交易哈希 (client.TransactionByHash) 获取交易详情。isPending 布尔值表示交易是否仍在交易池中。  
  Go  
  tx, isPending, err := client.TransactionByHash(ctx, txHash)

* **获取 Nonce**：client.PendingNonceAt 获取账户下一个可用的交易序号，这是发送新交易的必需参数。  
  Go  
  nonce, err := client.PendingNonceAt(ctx, fromAddress)

* **建议 Gas 价格**：client.SuggestGasPrice 从节点获取一个市场建议的 Gas 价格。  
  Go  
  gasPrice, \_ := client.SuggestGasPrice(ctx)

* **查询交易收据**：client.TransactionReceipt 是获取交易执行结果的唯一方式，返回一个包含状态、Gas 使用量和事件日志的收据。

| 字段名 | 数据类型 | 描述 |
| :---- | :---- | :---- |
| Status | uint64 | 1 代表成功，0 代表失败（回滚）。 |
| GasUsed | uint64 | 交易实际消耗的 Gas。 |
| Logs | \*types.Log | 合约发出的事件日志列表。 |
| ContractAddress | common.Address | 如果是合约部署交易，则为新合约地址。 |

## **3\. 发送交易（写入操作）**

改变区块链状态的操作都需要通过发送签名交易来完成。

1. **构建交易**：使用 types.NewTransaction 创建一个交易对象，包含 Nonce、目标地址、金额、Gas 上限、Gas 价格和数据。对于原生 ETH 转账，数据字段为 nil。  
2. **签名交易**：  
   * 使用 crypto.HexToECDSA 从十六进制字符串加载私钥。  
   * 使用 client.ChainID 获取当前链 ID 以防止重放攻击。  
   * types.SignTx 配合 types.NewEIP155Signer(chainID) 使用私钥对交易进行签名。  
3. **广播交易**：client.SendTransaction 将签名后的交易广播到以太坊网络。

## **4\. 智能合约交互**

go-ethereum 提供了强大的工具来简化与智能合约的交互，通常与 abigen 工具生成的 Go 代码配合使用。

### **4.1 合约部署**

* **底层方式**：手动解析 ABI (abi.JSON)，编码构造函数参数 (parsedABI.Pack)，然后将字节码和参数合并作为交易的 Data 字段发送。client.EstimateGas 可用于估算部署成本。  
* **abigen 方式（推荐）**：  
  1. 使用 bind.NewKeyedTransactorWithChainID 创建一个包含认证信息的 \*bind.TransactOpts 对象。  
  2. 调用 abigen 生成的 DeployContract 函数，传入 auth 对象、客户端实例和构造函数参数。

### **4.2 与已部署合约交互**

1. **加载合约实例**：使用 abigen 生成的 NewContract 函数，传入合约地址和客户端实例，获取一个类型安全的合约对象。  
   Go  
   storageContract, \_ := contract.NewContract(contractAddress, client)

2. **调用只读函数 (Call)**：直接调用 abigen 生成的对应方法。这类调用是免费和同步的。  
   Go  
   version, err := storageContract.Version(nil) // \`nil\` 使用默认 CallOpts

3. **调用写入函数 (Transaction)**：  
   * 创建一个 \*bind.TransactOpts 对象。  
   * 调用 abigen 生成的对应方法，传入 TransactOpts 和函数参数。该调用会返回一个 \*types.Transaction 对象。  
   * 使用 bind.WaitMined 辅助函数来阻塞并等待交易被打包，最终返回交易收据。

Go  
tx, err := storeContract.SetItem(opts, key, value)  
receipt, err := bind.WaitMined(ctx, client, tx)

## **5\. 事件订阅（WebSocket）**

通过 WebSocket 连接可以实时监听链上事件。

### **5.1 订阅新区块**

使用 wsClient.SubscribeNewHead 订阅新区块头的产生。事件会通过一个 chan \*types.Header 通道被接收。

### **5.2 订阅合约日志**

1. **构建过滤查询**：创建一个 ethereum.FilterQuery 实例，指定要监听的 Addresses（合约地址）和 Topics。Topics 通常是事件签名的 Keccak256 哈希。  
2. **开始订阅**：client.SubscribeFilterLogs 使用查询条件创建一个订阅，并通过一个 chan types.Log 通道返回匹配的日志。  
3. **解码日志**：  
   * 使用 contractAbi.UnpackIntoInterface 来解码日志的 Data 字段（包含未索引的参数）。  
   * 手动从 vLog.Topics 数组中提取被 indexed 关键字标记的参数。