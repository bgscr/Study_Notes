package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// BlockEvent 对应了生产者发送的消息结构，
// 同时也用作我们的 GORM 模型，将直接映射到数据库的表。
type BlockEvent struct {
	ID          uint64 `gorm:"primaryKey"`  // GORM 会自动使用 ID 作为主键
	BlockNumber uint64 `gorm:"uniqueIndex"` // 为区块号创建唯一索引，防止重复记录
	BlockHash   string `gorm:"size:66"`
	Timestamp   uint64
	TxCount     int
	CreatedAt   time.Time // GORM 会自动管理创建时间
}

func main() {
	// --- 1. 连接到 MySQL 数据库 ---
	// 格式: "用户名:密码@tcp(地址:端口)/数据库名?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:123456@tcp(127.0.0.1:3306)/blockchain_data?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到 MySQL 数据库: %v", err)
	}
	log.Println("成功连接到 MySQL 数据库")

	// --- 2. 自动迁移 schema ---
	// GORM 的 AutoMigrate 功能会自动检查 `block_events` 表是否存在，
	// 如果不存在，它会根据 BlockEvent 结构体的定义创建这张表。
	log.Println("正在迁移数据库 schema...")
	err = db.AutoMigrate(&BlockEvent{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移成功")

	// --- 3. 配置 Kafka 消费者 (Reader) ---
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "eth-events",
		GroupID:  "eth-consumer-group", // 消费者组 ID，Kafka 用它来跟踪消费进度
		MinBytes: 10e3,                 // 10KB
		MaxBytes: 10e6,                 // 10MB
	})
	defer kafkaReader.Close()
	log.Println("Kafka 消费者已配置完成")

	// --- 4. 监听和处理消息 ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("上下文已取消，正在关闭消费者...")
				return
			default:
				// a. 从 Kafka 读取消息 (这是一个阻塞操作)
				m, err := kafkaReader.FetchMessage(ctx)
				if err != nil {
					// 如果是上下文取消导致的错误，则正常退出
					if ctx.Err() != nil {
						return
					}
					log.Printf("读取消息失败: %v\n", err)
					continue
				}
				log.Printf("从分区 %d 读取到消息: Key=%s\n", m.Partition, string(m.Key))

				// b. 反序列化消息
				var event BlockEvent
				err = json.Unmarshal(m.Value, &event)
				if err != nil {
					log.Printf("反序列化消息失败: %v\n", err)
					// 即使消息格式错误，我们依然需要提交 offset，以防卡住
					kafkaReader.CommitMessages(ctx, m)
					continue
				}

				// c. 将数据存入数据库
				// 使用 Create 方法插入一条新记录
				result := db.Create(&event)
				if result.Error != nil {
					log.Printf("写入数据库失败 (区块 #%d): %v\n", event.BlockNumber, result.Error)
				} else {
					log.Printf("成功将区块 #%d 的数据存入数据库，影响行数: %d\n", event.BlockNumber, result.RowsAffected)
				}

				// d. 提交 offset
				// 这一步非常关键，它告诉 Kafka 这条消息已经被成功处理了，
				// 下次消费者重启时不会再重复消费。
				if err := kafkaReader.CommitMessages(ctx, m); err != nil {
					log.Printf("提交 offset 失败: %v", err)
				}
			}
		}
	}()

	log.Println("消费者已启动，正在等待消息...")
	<-sigChan
	log.Println("接收到中断信号，正在优雅地关闭...")
	cancel()
	// 等待 goroutine 正常退出
	time.Sleep(2 * time.Second)
	log.Println("程序已关闭")

}
