package main

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	topic := "my_topic"
	brokers := "b-1.staging-kafka-cluster.tu3o28.c4.kafka.eu-west-3.amazonaws.com:9094"
	producer, err := sarama.NewSyncProducer([]string{brokers}, nil)
	if err != nil {
		fmt.Println("error %w", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	app := fiber.New()

	app.Use(
		logger.New(),
	)

	app.Get("/send", func(c *fiber.Ctx) error {
		// publish sync
		msg := &sarama.ProducerMessage {
			Topic: topic,
			Value: sarama.StringEncoder(c.Query("msg")),
		}
		p, o, err := producer.SendMessage(msg)
		if err != nil {
			fmt.Println("Error publish: ", err.Error())
		}

		fmt.Println("Partition: ", p)
		fmt.Println("Offset: ", o)

		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
