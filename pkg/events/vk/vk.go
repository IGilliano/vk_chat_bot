package vk

import (
	"errors"
	"vk_chat_bot/pkg/clients/vk_client"
	"vk_chat_bot/pkg/events"
	"vk_chat_bot/pkg/storage"
)

type Processor struct {
	vk      *vk_client.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	FromId int `json:"from_id"`
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMeta      = errors.New("unknown meta type")
)

func NewProcessor(client *vk_client.Client, storage storage.Storage) *Processor {
	return &Processor{
		vk:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch() ([]events.Event, error) {
	updates, err := p.vk.LongPollServer()
	if err != nil {
		return nil, err
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	return res, nil

}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case "message_new":
		return p.processMessage(event)
	default:
		return ErrUnknownEventType
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return err
	}

	if err = p.doCmd(event.Text, meta.FromId); err != nil {
		return err
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, ErrUnknownMeta
	}

	return res, nil
}

func event(upd vk_client.Update) events.Event {
	upd.Type = fetchType(upd)

	res := events.Event{
		Type: upd.Type,
		Text: fetchText(upd),
	}

	if res.Type == "message_new" {
		res.Meta = Meta{
			FromId: upd.Object.Message.FromId,
		}

	}

	return res

}

func fetchType(upd vk_client.Update) string {
	return upd.Type

}

func fetchText(upd vk_client.Update) string {
	return upd.Object.Message.Text

}
