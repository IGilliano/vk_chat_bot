package event_consumer

import (
	"log"
	"time"
	"vk_chat_bot/pkg/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		events, err := c.fetcher.Fetch()
		if err != nil {
			log.Printf("Errror! Consumer: %s", err.Error())
			continue
		}

		if len(events) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err = c.handleEvents(events); err != nil {
			log.Print(err)
			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("Got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("Cant handle event: %s", err.Error())

			continue
		}

	}

	return nil
}
