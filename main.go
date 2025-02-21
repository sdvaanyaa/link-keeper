package main

import (
	"flag"
	"log"
	tgClient "saveBot/clients/telegram"
	"saveBot/consumer/event_consumer"
	"saveBot/events/telegram"
	"saveBot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	//sqliteStoragePath = "data/sqlite/storage.db"
	batchSize = 100
)

func main() {
	s := files.New(storagePath)
	//s, err := sqlite.New(sqliteStoragePath)
	//if err != nil {
	//	log.Fatalf("can't connect to storage: %s", err)
	//}
	//
	//if err := s.Init(context.TODO()); err != nil {
	//	log.Fatalf("can't init storage: %s", err)
	//}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped")
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("token is required")
	}

	return *token
}
