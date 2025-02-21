package event_consumer

import (
	"context"
	"log"
	"saveBot/events"
	"time"
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
		gotEvents, err := c.fetcher.Fetch(context.Background(), c.batchSize)
		if err != nil {
			log.Printf("error fetching events: %v", err)

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(context.Background(), gotEvents); err != nil {
			log.Printf("error handling events: %s", err)
		}

	}
}

func (c Consumer) handleEvents(ctx context.Context, events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(ctx, event); err != nil {
			log.Printf("can't handle event: %v", err)

			continue
		}
	}
	return nil
}
