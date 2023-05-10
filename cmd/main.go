package main

import (
	"log"
	"vk_chat_bot/pkg/clients/vk_client"
	"vk_chat_bot/pkg/consumer/event-consumer"
	"vk_chat_bot/pkg/events/vk"
	"vk_chat_bot/pkg/storage/files"
)

const (
	host      = "api.vk.com"
	basePath  = "storage"
	batchSize = 100
)

func main() {

	client := vk_client.New(host, token)
	eventsProcessor := vk.NewProcessor(client, files.New(basePath))

	log.Printf("Listen and serve")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal()
	}
}
