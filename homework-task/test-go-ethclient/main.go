package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"testGoEthclient/contract"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/IYvnOjyIUkrmjNsCRU8XR")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	background := context.Background()

	blance1, _ := client.BalanceAt(background, common.HexToAddress("0xBCe923aa9e09e76983061E49400555D407F827Ef"), nil)
	fmt.Println("account1余额", blance1)
	blance2, _ := client.BalanceAt(background, common.HexToAddress("0xF40FEBfDf7e411953104cbBA14A0c77E5dDfcC3a"), nil)
	fmt.Println("account2余额", blance2)

	// 获取hash
	blockNumber := big.NewInt(8964804)
	block, err := client.BlockByNumber(background, blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	blockByHash, err := client.BlockByHash(background, common.HexToHash("0x7aba2405a6407c8efd9b6079f80dc359ba0afa84ec1667f3d1fab43bd995207d"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("查询区块：")
	fmt.Println(block.Number())
	fmt.Println(blockByHash.Number())
	fmt.Println(blockByHash.GasUsed())

	//查询交易
	tx1 := block.Transaction(common.HexToHash("0x7a9d6b47eac4207fa84787e962051fc6e997249677e70ac8037aece0caf60f56"))

	tx2, isPending, err := client.TransactionByHash(background, common.HexToHash("0x7a9d6b47eac4207fa84787e962051fc6e997249677e70ac8037aece0caf60f56"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("交易：")
	fmt.Println(tx1.ChainId())
	fmt.Println(tx1.Hash())
	fmt.Println(tx2.ChainId())
	fmt.Println(tx2.Hash())
	fmt.Println(isPending)

	//转账操作
	fmt.Println("转账操作")
	nonce, err := client.PendingNonceAt(background, common.HexToAddress("0xBCe923aa9e09e76983061E49400555D407F827Ef"))
	if err != nil {
		log.Printf("获取nonce失败: %v\n", err)
		return
	}
	value := big.NewInt(10)
	gasLimit := uint64(21000) // ETH 转账的 gas limit 通常是 21000
	gasPrice, _ := client.SuggestGasPrice(background)

	toAddress := common.HexToAddress("0xF40FEBfDf7e411953104cbBA14A0c77E5dDfcC3a")
	// 创建交易
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	chainID, err := client.ChainID(background)
	if err != nil {
		log.Printf("获取Chain ID失败: %v\n", err)
		return
	}

	privateKey, err := crypto.HexToECDSA("6f7c629d3469c434d68a80c8485329808678c6d09b77f389b4314fb3a2a0e902")
	if err != nil {
		log.Fatalf("解析私钥失败: %v", err)
	}

	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	// err = client.SendTransaction(background, signedTx)
	// if err != nil {
	// 	log.Printf("发送交易失败: %v\n", err)
	// 	return
	// }
	//ETH转账交易已发送, Hash: 0x048b2cdb144bbf3f8c0904bdc9d2775ad7112db2e72baf3f6aef925b6910d76c
	fmt.Printf("ETH转账交易已发送, Hash: %s\n", signedTx.Hash().Hex())
	theTxHash := common.HexToHash("0x048b2cdb144bbf3f8c0904bdc9d2775ad7112db2e72baf3f6aef925b6910d76c")
	_, isPending, _ = client.TransactionByHash(background, theTxHash)

	//receipt, err := client.TransactionReceipt(background, signedTx.Hash())
	if isPending {
		log.Println("交易pending")
	}
	receipt, err := client.TransactionReceipt(background, theTxHash)
	if err != nil {
		log.Printf("查询收据失败: %v\n", err)
		return
	}
	fmt.Printf("收据状态 (1=成功, 0=失败): %d\n", receipt.Status)
	fmt.Printf("消耗的 Gas: %d\n", receipt.GasUsed)
	fmt.Printf("包含的日志/事件数量: %d\n", len(receipt.Logs))
	fmt.Println("转账操作结束")

	//EstimateGas 步骤
	// 1. 解析ABI字符串
	parsedABI, err := abi.JSON(strings.NewReader(contract.ContractMetaData.ABI))
	if err != nil {
		log.Fatalf("解析ABI失败: %v", err)
	}

	// 2. 准备构造函数所需的参数
	versionString := "1.0"

	// 3. ABI编码构造函数参数
	// 对于构造函数，Pack的第一个参数（方法名）留空
	encodedArgs, err := parsedABI.Pack("", versionString)
	if err != nil {
		log.Fatalf("编码构造函数参数失败: %v", err)
	}

	// 4. 合并字节码和编码后的参数
	contractBytecode := common.FromHex(contract.ContractMetaData.Bin)
	data := append(contractBytecode, encodedArgs...)

	if len(data) == 0 {
		log.Fatal("合约字节码为空，请检查 'contractBin' 常量")
	}
	msg := ethereum.CallMsg{
		From: common.HexToAddress("0xBCe923aa9e09e76983061E49400555D407F827Ef"),
		To:   nil, // To 为 nil 表示是合约创建
		Data: data,
	}

	estimatedGas, err := client.EstimateGas(background, msg)
	if err != nil {
		log.Fatalf("计算estimatedGas失败: %v\n", err)
	}
	fmt.Println("计算合约的estimatedGas：", estimatedGas)
	// contractAddress, storageContract := deployContract(client, &background, estimatedGas)

	contractAddress := common.HexToAddress("0x73917A74CeF6Ad755bDc6Bb712532A182C3086B1")

	storageContract, _ := contract.NewContract(contractAddress, client)
	if err != nil {
		log.Fatal("加载合约失败", err)
	}

	version, err := storageContract.Version(nil)
	if err != nil {
		log.Fatalf("Failed to retrieve version: %v", err)
	}
	fmt.Printf("Contract Version: %s\n", version)

	wsClient, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/IYvnOjyIUkrmjNsCRU8XR")
	if err != nil {
		log.Fatal(err)
	}
	defer wsClient.Close()
	sub := subscribe(wsClient, &background, contractAddress)
	defer sub.Unsubscribe()

	callContract(client, &background, storageContract)

	var key [32]byte
	copy(key[:], []byte("demo_save_key"))
	result, err := storageContract.Items(nil, key)
	if err != nil {
		log.Fatalf("Failed to retrieve item: %v", err)
	}
	fmt.Printf("Value for key 'demo_save_key': %s\n", string(result[:]))

	cancelCtx, cancel := context.WithTimeout(background, 30*time.Second)
	defer cancel()
	fmt.Println("订阅区块链事件开始", err)
	blockSub := subscibeBlock(wsClient, &cancelCtx)
	defer blockSub.Unsubscribe()

	<-cancelCtx.Done()
	err = cancelCtx.Err()
	fmt.Println("订阅区块链事件结束", err)
}

func subscibeBlock(wsClient *ethclient.Client, background *context.Context) ethereum.Subscription {
	// 创建一个用于接收区块头的通道
	// 注意通道的类型是 *types.Header
	headers := make(chan *types.Header)

	// 订阅新的区块头事件
	sub, err := wsClient.SubscribeNewHead(*background, headers)
	if err != nil {
		log.Fatalf("无法订阅新区块头事件: %v", err)
	}
	fmt.Println("已成功订阅新区块事件，正在等待新区块产生...")

	// 4. 在循环中监听和处理事件
	go func() {
		for {
			select {
			case err := <-sub.Err():
				if err != nil {
					log.Fatalf("订阅发生错误: %v", err)
				}

			case header := <-headers:
				// 当有新区块产生时，这里的代码会执行
				log.Println("--- 发现新区块 ---")
				log.Printf("区块号: %s\n", header.Number.String())
				log.Printf("区块哈希: %s\n", header.Hash().Hex())
				log.Printf("时间戳: %d\n", header.Time)
				log.Printf("Gas 使用量: %d\n", header.GasUsed)
				log.Printf("叔块哈希: %s\n", header.UncleHash.Hex())

				// 如果需要获取区块内的完整交易信息，需要额外发起一次请求
				block, err := wsClient.BlockByHash(*background, header.Hash())
				if err != nil {
					log.Printf("无法根据哈希获取完整区块信息: %v", err)
				} else {
					log.Printf("区块中交易数量: %d\n", len(block.Transactions()))
					// 打印前几个交易的哈希
					for i, tx := range block.Transactions() {
						if i >= 3 {
							log.Println("...等更多交易")
							break
						}
						log.Printf("  -> 交易 #%d 哈希: %s\n", i+1, tx.Hash().Hex())
					}
				}
			}
		}
	}()
	return sub

}

func subscribe(client *ethclient.Client, background *context.Context, contractAddress common.Address) ethereum.Subscription {

	eventSignature := []byte("ItemSet(bytes32,bytes32)")
	hash := crypto.Keccak256Hash(eventSignature)
	topic0 := hash.Hex()
	fmt.Println("Event Signature Hash (Topic 0):", topic0)

	// 创建过滤查询
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics: [][]common.Hash{
			{hash}, // 只关心 ItemSet 事件
		},
	}

	// 创建一个用于接收日志的通道
	logs := make(chan types.Log)

	// 开始订阅
	sub, err := client.SubscribeFilterLogs(*background, query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	fmt.Println("Subscribed to contract events. Waiting for new logs...")

	contractAbi, _ := abi.JSON(strings.NewReader(contract.ContractMetaData.ABI))
	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatalf("Subscription error: %v", err)
			case vLog := <-logs:
				// --- 关键的安全检查 ---
				// 检查1: Topic[0] 是不是 ItemSet 的签名？
				// 检查2: Topics 数组的长度是否足够？(对于 ItemSet 应该是2)
				// 这两步确保了我们只处理正确的事件，并能安全地访问 Topics[1]

				// 当有新事件时，这里的代码会执行
				fmt.Printf("Received new log: BlockNumber: %d, TxHash: %s\n", vLog.BlockNumber, vLog.TxHash.Hex())

				var itemSetEvent contract.ContractItemSet
				// `Topics` 包含了事件签名和被索引(indexed)的参数
				// `Data` 包含了未被索引的参数
				// `UnpackIntoInterface` 会将 `Data` 字段解包到结构体中
				err := contractAbi.UnpackIntoInterface(&itemSetEvent, "ItemSet", vLog.Data)
				if err != nil {
					log.Printf("解包事件数据失败: %v", err)
					continue
				}

				// 从 Topics 中手动填充被索引的字段
				// 根据合约，`key` 是第一个 `indexed` 参数，所以它在 `Topics[1]`
				// `Topics[0]` 是事件签名哈希
				itemSetEvent.Key = vLog.Topics[1]
				// --- 现在可以使用强类型的事件数据了 ---
				log.Println("事件 'ItemSet' 已解析:")
				// bytes32 类型是 [32]byte 数组，我们通常以十六进制字符串形式查看它
				keyHex := hex.EncodeToString(itemSetEvent.Key[:])
				valueHex := hex.EncodeToString(itemSetEvent.Value[:])
				log.Printf("  Key (Hex): 0x%s\n", keyHex)
				log.Printf("  Value (Hex): 0x%s\n", valueHex)

				keyString := string(bytes.TrimRight(itemSetEvent.Key[:], "\x00"))
				valueString := string(bytes.TrimRight(itemSetEvent.Value[:], "\x00"))
				log.Printf("  Key (String): %s\n", keyString)
				log.Printf("  Value (String): %s\n", valueString)
			}
		}
	}()

	return sub
}

func callContract(client *ethclient.Client, background *context.Context, storeContract *contract.Contract) {
	// 准备数据
	var key [32]byte
	var value [32]byte
	copy(key[:], []byte("demo_save_key"))
	copy(value[:], []byte("中文测试value"))

	// 初始化交易opt实例
	privateKey, _ := crypto.HexToECDSA("6f7c629d3469c434d68a80c8485329808678c6d09b77f389b4314fb3a2a0e902")
	chainID, err := client.ChainID(*background)
	if err != nil {
		log.Fatal(err)
	}
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	// 1. 调用合约方法，发送交易
	tx, err := storeContract.SetItem(opt, key, value)
	if err != nil {
		log.Fatalf("调用合约方法失败: %v", err)
	}
	fmt.Println("交易已发送，哈希:", tx.Hash().Hex())

	// 2. 等待交易被打包并获取收据
	fmt.Println("等待交易被矿工打包...")
	receipt, err := bind.WaitMined(*background, client, tx)
	if err != nil {
		log.Fatalf("等待交易打包失败: %v", err)
	}

	// 3. 从收据中读取GasUsed
	fmt.Println("交易成功打包!")
	fmt.Println("实际使用的Gas:", receipt.GasUsed)

	// 你还可以检查交易状态
	if receipt.Status == 1 {
		fmt.Println("交易状态: 成功 (Success)")
	} else {
		fmt.Println("交易状态: 失败 (Reverted)")
	}
}

func deployContract(client *ethclient.Client, background *context.Context, estimatedGas uint64) (common.Address, *contract.Contract) {
	fmt.Println("部署合约")

	privateKey, _ := crypto.HexToECDSA("6f7c629d3469c434d68a80c8485329808678c6d09b77f389b4314fb3a2a0e902")
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(*background, fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.ChainID(*background)
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)   // in wei
	auth.GasLimit = estimatedGas // in units
	gasPrice, _ := client.SuggestGasPrice(*background)
	auth.GasPrice = gasPrice

	input := "1.0"
	address, tx, instance, err := contract.DeployContract(auth, client, input)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("部署地址Hex", address.Hex())
	fmt.Println(tx.Gas())

	receipt, err := bind.WaitMined(*background, client, tx)
	if err != nil {
		log.Fatalf("等待交易打包失败: %v", err)
	}
	if receipt.Status == 1 {
		fmt.Println("部署合约交易状态: 成功 (Success)")
	} else {
		fmt.Println("部署合约交易状态: 失败 (Reverted)")
	}
	return address, instance
}
