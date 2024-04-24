package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	log.Println("Starting app...")

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"client.id":         "cli",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	topic := "msg"
	i := 1
sending:
	for {
		select {
		case <-ctx.Done():
			log.Println("Closing Kafka connection...")
			producer.Close()
			break sending
		default:
		}
		value := fmt.Sprintf("msg %d", i)
		log.Printf("Sending message: %s", value)
		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(value)}
		err := producer.Produce(msg, nil)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * 1)
		i++
	}

	log.Println("Bye bye...")
}
