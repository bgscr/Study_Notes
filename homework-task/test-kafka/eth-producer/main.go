package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/segmentio/kafka-go"
)

// BlockEvent 是我们要发送到 Kafka 的消息结构体
// 它包含了区块的一些关键信息
type BlockEvent struct {
	BlockNumber uint64 `json:"blockNumber"`
	BlockHash   string `json:"blockHash"`
	Timestamp   uint64 `json:"timestamp"`
	TxCount     int    `json:"txCount"`
}

func main() {
	// --- 1. 连接到以太坊节点 ---
	// 这里我们使用一个公共的以太坊节点 WebSocket 地址。
	// WebSocket (wss) 对于实时订阅事件来说是必需的。
	// 你也可以换成你自己的节点地址，例如 Infura 或 Alchemy 提供的。
	client, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/IYvnOjyIUkrmjNsCRU8XR")
	if err != nil {
		log.Fatalf("无法连接到以太坊节点: %v", err)
	}
	defer client.Close()
	log.Println("成功连接到以太坊节点")

	// --- 2. 配置 Kafka 生产者 (Writer) ---
	// 我们连接到本地运行的 Kafka，并指定要写入的 topic。
	kafkaWriter := &kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"),
		Topic: "eth-events",
		// 使用 Hash 均衡器。它将根据消息的 Key 来决定发送到哪个分区。
		// 这保证了具有相同 Key 的消息会进入同一个分区。
		Balancer: &kafka.Hash{},
	}
	defer kafkaWriter.Close()
	log.Println("Kafka 生产者已配置完成")

	// --- 3. 订阅新的区块头事件 ---
	// 我们创建一个 channel 来接收新的区块头信息。
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalf("无法订阅新的区块头: %v", err)
	}
	defer sub.Unsubscribe()
	log.Println("已成功订阅以太坊新区块事件")

	// --- 4. 监听事件并处理 ---
	// 我们使用一个 goroutine 来异步处理订阅、错误和程序中断信号。
	// 这是 Go 中处理并发任务的常用模式。
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听操作系统中断信号 (例如 Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			// a. 接收到新的区块头
			case header := <-headers:
				log.Printf("接收到新区块: #%d\n", header.Number.Uint64())

				// 通过区块号获取完整的区块信息，因为 header 里信息不全
				block, err := client.BlockByNumber(ctx, header.Number)
				if err != nil {
					log.Printf("获取区块 #%d 失败: %v\n", header.Number.Uint64(), err)
					continue // 跳过这个区块
				}

				// b. 构造消息
				event := BlockEvent{
					BlockNumber: block.NumberU64(),
					BlockHash:   block.Hash().Hex(),
					Timestamp:   block.Time(),
					TxCount:     len(block.Transactions()),
				}

				// 将结构体序列化为 JSON 字节
				eventBytes, err := json.Marshal(event)
				if err != nil {
					log.Printf("序列化区块事件失败: %v\n", err)
					continue
				}

				// c. 将消息发送到 Kafka
				err = kafkaWriter.WriteMessages(ctx, kafka.Message{
					Key:   []byte(block.Hash().Hex()), // 使用区块哈希作为 Key
					Value: eventBytes,
				})

				if err != nil {
					log.Printf("发送消息到 Kafka 失败: %v\n", err)
				} else {
					log.Printf("成功发送区块 #%d 的消息到 Kafka\n", block.NumberU64())
				}

			// b. 订阅出错
			case err := <-sub.Err():
				log.Printf("订阅错误: %v\n", err)
				cancel() // 发生错误时取消上下文，准备退出
				return

			// c. 程序即将退出
			case <-ctx.Done():
				log.Println("上下文已取消，正在关闭...")
				return
			}
		}
	}()

	log.Println("生产者已启动，正在等待新的区块事件...")
	// 等待中断信号
	<-sigChan
	log.Println("接收到中断信号，正在优雅地关闭...")
	cancel() // 发送取消信号给 goroutine
}
