package telegram

import (
	"context"
	"errors"
	"saveBot/clients/telegram"
	"saveBot/events"
	"saveBot/lib/errwrap"
	"saveBot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errwrap.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, upd := range updates {
		res = append(res, event(upd))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func (p *Processor) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.ProcessMessage(ctx, event)
	default:
		return errwrap.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) ProcessMessage(ctx context.Context, event events.Event) (err error) {
	defer func() { err = errwrap.WrapIfErr("can't process message", err) }()

	metaData, err := meta(event)
	if err != nil {
		return err
	}

	if err := p.doCmd(ctx, event.Text, metaData.ChatID, metaData.Username); err != nil {
		return err
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, ErrUnknownMetaType
	}
	return res, nil
}
