package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	log.Println("Starting app...")

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "go",
		"auto.offset.reset": "smallest"})
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	topic := "msg"
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Printf("Cannot subscribe to topic [%s]: %v\n", topic, err)
	}
receiving:
	for {
		select {
		case <-ctx.Done():
			log.Println("Closing Kafka connection...")
			consumer.Close()
			break receiving
		default:
		}
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v\n", err)
			continue
		}
		log.Printf("Received message: %s\n", msg.Value)
	}

	log.Println("Bye bye...")
}
