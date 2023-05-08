package main

import (
	"log"
	"vk_chat_bot/pkg/clients/vk_client"
	"vk_chat_bot/pkg/config"
	event_consumer "vk_chat_bot/pkg/consumer/event-consumer"
	"vk_chat_bot/pkg/events/vk"
	"vk_chat_bot/pkg/storage/files"
)

const (
	host      = "api.vk.com"
	basePath  = "storage"
	batchSize = 100
)

func main() {

	cfg := config.MustToken()

	client := vk_client.New(host, cfg.VkBotToken)
	eventsProcessor := vk.NewProcessor(client, files.New(basePath))

	log.Printf("Listen and serve")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal()
	}
}
