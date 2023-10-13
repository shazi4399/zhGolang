package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	brokerAddress = "localhost:9092"
	topic         = "hellokafka"
)

func main() {
	// 创建一个 Kafka 生产者
	producer := createKafkaProducer()

	// 启动一个协程，用于生产消息
	go produceMessages(producer)

	// 创建一个 Kafka 消费者
	consumer := createKafkaConsumer()

	// 启动一个协程，用于消费消息
	go consumeMessages(consumer)

	// 等待一段时间，让生产者和消费者有足够的时间运行
	time.Sleep(10 * time.Second)
}

func createKafkaProducer() *kafka.Writer {
	// 配置 Kafka 生产者
	config := kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	}

	// 创建 Kafka 生产者
	producer := kafka.NewWriter(config)

	return producer
}

func produceMessages(producer *kafka.Writer) {
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Message %d", i)

		// 发送消息到 Kafka
		err := producer.WriteMessages(context.Background(), kafka.Message{
			Value: []byte(message),
		})

		if err != nil {
			log.Printf("Failed to produce message: %v\n", err)
		} else {
			log.Printf("Produced message: %s\n", message)
		}

		time.Sleep(1 * time.Second)
	}

	// 关闭 Kafka 生产者
	err := producer.Close()
	if err != nil {
		log.Printf("Failed to close producer: %v\n", err)
	}
}

func createKafkaConsumer() *kafka.Reader {
	// 配置 Kafka 消费者
	config := kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		GroupID: "test-group",
	}

	// 创建 Kafka 消费者
	consumer := kafka.NewReader(config)

	return consumer
}

func consumeMessages(consumer *kafka.Reader) {
	for {
		message, err := consumer.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Failed to consume message: %v\n", err)
			continue
		}

		log.Printf("Consumed message: %s\n", string(message.Value))
	}
}
