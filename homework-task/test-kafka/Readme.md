# **以太坊-Kafka-MySQL 数据管道学习总结**

## **1\. 总体目标**

搭建一个实时数据管道，用于监听以太坊链上的新区块事件，通过 Kafka 消息队列进行传输，最终由消费者程序将数据存入 MySQL 数据库。

## **2\. 环境与工具准备**

### **2.1. WSL2 环境配置**

* **问题**: 在 WSL (Ubuntu) 中使用 sudo apt install 时，出现大量网络连接失败和超时的错误。  
* **原因**: apt 默认的软件源服务器在国外，国内访问不稳定。  
* **解决方案**:  
  1. 备份原始源列表：  
     sudo cp /etc/apt/sources.list /etc/apt/sources.list.bak

  2. 编辑源列表文件 sudo nano /etc/apt/sources.list，注释掉所有内容，并替换为国内镜像源（如阿里云）：  
     deb \[http://mirrors.aliyun.com/ubuntu/\](http://mirrors.aliyun.com/ubuntu/) noble main restricted universe multiverse  
     deb \[http://mirrors.aliyun.com/ubuntu/\](http://mirrors.aliyun.com/ubuntu/) noble-updates main restricted universe multiverse  
     deb \[http://mirrors.aliyun.com/ubuntu/\](http://mirrors.aliyun.com/ubuntu/) noble-backports main restricted universe multiverse  
     deb \[http://mirrors.aliyun.com/ubuntu/\](http://mirrors.aliyun.com/ubuntu/) noble-security main restricted universe multiverse

  3. 刷新软件源列表：  
     sudo apt update

### **2.2. Kafka 安装 (KRaft 模式)**

* **依赖**: 安装 OpenJDK。  
  sudo apt install default-jdk  
  java \-version \# 验证安装

* **下载**:  
  * **问题**: 官网 downloads.apache.org 的直接链接可能因版本归档而失效 (404 Not Found)。  
  * **解决方案**: 使用官方镜像分发网络 (dlcdn.apache.org) 的链接。  
    wget \[https://dlcdn.apache.org/kafka/3.7.1/kafka\_2.13-3.7.1.tgz\](https://dlcdn.apache.org/kafka/3.7.1/kafka\_2.13-3.7.1.tgz)  
    tar \-xzf kafka\_2.13-3.7.1.tgz  
    cd kafka\_2.13-3.7.1

* **KRaft 模式启动流程 (无需 ZooKeeper)**:  
  1. **生成集群 ID**:  
     KAFKA\_CLUSTER\_ID="$(bin/kafka-storage.sh random-uuid)"

  2. **格式化存储目录**: 使用 KRaft 配置文件和集群 ID 初始化。  
     bin/kafka-storage.sh format \-t $KAFKA\_CLUSTER\_ID \-c config/kraft/server.properties

  3. **启动 Kafka 服务器**:  
     bin/kafka-server-start.sh config/kraft/server.properties

* **Topic 管理**:  
  1. **创建多分区 Topic**:  
     \# 创建一个名为 eth-events，有3个分区的 Topic  
     bin/kafka-topics.sh \--create \--topic eth-events \--bootstrap-server localhost:9092 \--partitions 3 \--replication-factor 1

  2. **删除 Topic** (用于重建):  
     bin/kafka-topics.sh \--delete \--topic eth-events \--bootstrap-server localhost:9092

## **3\. Go 生产者 (eth-producer)**

### **3.1. 项目设置**

mkdir eth-producer && cd eth-producer  
go mod init eth-producer  
go get \[github.com/ethereum/go-ethereum\](https://github.com/ethereum/go-ethereum)  
go get \[github.com/segmentio/kafka-go\](https://github.com/segmentio/kafka-go)

### **3.2. 核心代码逻辑**

1. **定义消息结构体**:  
   type BlockEvent struct {  
       BlockNumber uint64 \`json:"blockNumber"\`  
       BlockHash   string \`json:"blockHash"\`  
       Timestamp   uint64 \`json:"timestamp"\`  
       TxCount     int    \`json:"txCount"\`  
   }

2. **连接以太坊节点**:  
   * 必须使用 WebSocket (wss://) 协议以支持实时事件订阅。  
   * client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/YOUR\_PROJECT\_ID")  
3. **配置 Kafka 生产者**:  
   * kafka.Writer 用于发送消息。  
   * **关键点**: 设置 Balancer 为 \&kafka.Hash{}。这会根据消息的 Key 将消息分发到不同的分区，保证了相同 Key 的消息落在同一分区。

kafkaWriter := \&kafka.Writer{  
    Addr:     kafka.TCP("localhost:9092"),  
    Topic:    "eth-events",  
    Balancer: \&kafka.Hash{},  
}

4. **订阅与发送**:  
   * 使用 client.SubscribeNewHead 订阅新区块头事件。  
   * 在一个 goroutine 中使用 select 监听新区块、订阅错误和程序退出信号。  
   * **流程**: 收到区块头 \-\> 获取完整区块信息 \-\> 构造 BlockEvent \-\> json.Marshal 序列化 \-\> 调用 kafkaWriter.WriteMessages 发送。  
   * **消息 Key**: 使用区块哈希 block.Hash().Hex() 作为消息的 Key，用于分区路由。

## **4\. Go 消费者 (eth-consumer)**

### **4.1. 项目设置**

mkdir eth-consumer && cd eth-consumer  
go mod init eth-consumer  
go get \[github.com/segmentio/kafka-go\](https://github.com/segmentio/kafka-go)  
go get gorm.io/gorm  
go get gorm.io/driver/mysql

### **4.2. 核心代码逻辑**

1. **定义 GORM 模型**:  
   * 复用 BlockEvent 结构体，并添加 GORM 标签。  
   * gorm:"primaryKey": 定义主键。  
   * gorm:"uniqueIndex": 为 BlockNumber 添加唯一索引，防止数据重复入库。

type BlockEvent struct {  
    ID          uint      \`gorm:"primaryKey"\`  
    BlockNumber uint64    \`gorm:"uniqueIndex"\`  
    // ... 其他字段  
}

2. **连接 MySQL 数据库**:  
   * 使用 GORM 和 MySQL 驱动。  
   * DSN 字符串格式: "user:pass@tcp(host:port)/dbname?charset=utf8mb4\&parseTime=True\&loc=Local"  
3. **自动迁移 Schema**:  
   * db.AutoMigrate(\&BlockEvent{})  
   * 程序启动时自动创建或更新 block\_events 表结构。  
4. **配置 Kafka 消费者**:  
   * kafka.NewReader 用于读取消息。  
   * **关键点**: GroupID 字段。所有具有相同 GroupID 的消费者实例组成一个消费者组。Kafka 会自动将 Topic 的分区分配给组内的成员，实现负载均衡和容错。

kafkaReader := kafka.NewReader(kafka.ReaderConfig{  
    Brokers:  \[\]string{"localhost:9092"},  
    Topic:    "eth-events",  
    GroupID:  "eth-consumer-group",  
})

5. **消费与存储**:  
   * 在循环中调用 kafkaReader.FetchMessage(ctx) 阻塞式地获取消息。  
   * **流程**: 获取消息 \-\> json.Unmarshal 反序列化到 BlockEvent 结构体 \-\> 调用 db.Create(\&event) 写入数据库。  
   * **关键点**: kafkaReader.CommitMessages(ctx, m)。在消息处理完毕（无论成功或失败）后，必须提交 offset。这会通知 Kafka 该消息已被消费，防止程序重启后重复消费。

## **5\. 并行消费与扩展**

* **原理**: Kafka 的一个分区在同一时间只能被一个消费者组里的一个消费者实例消费。  
* **实现**:  
  1. Topic 必须创建为多分区（例如3个）。  
  2. 生产者使用 kafka.Hash 等策略将消息分发到不同分区。  
  3. 消费者代码**无需任何修改**。  
  4. 启动多个（例如3个）使用**相同 GroupID** 的消费者进程。Kafka 会自动将分区（0, 1, 2）分配给这三个进程，从而实现并行处理。