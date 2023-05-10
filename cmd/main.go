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
	token     = "vk1.a.LVvYNHuPbKLve39PzdK2ROjvxdftBCQgeQ5KEjd2Sh9lFdA_643XNTfoDLKcQGrDFdAs9L-E7AsE-Jpi85tVAjGLwDx0xdaiSCOst5K9eQ5rUFVtrRhbJESnJvmnYrumwCHyfLgdRTRnVs3uaO9zUTHYJg6wPuQyEPywT92I6EUJAM8vYnyG7KsBJbAZVAkdVa1KcGLDkXjKtV4w5b7voA"
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
