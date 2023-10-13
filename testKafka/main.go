package main

import (
	"fmt"
	"github.com/IBM/sarama"
)

func main() {
	config := sarama.NewConfig()
	brokers := []string{"localhost:9092"} // Kafka broker地址
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		fmt.Println("连接Kafka失败：", err)
		return
	}
	defer client.Close()

	fmt.Println("连接Kafka成功！")
}
