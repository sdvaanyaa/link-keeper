package main

import (
	"flag"
	"log"
	"saveBot/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken())

	// fetcher = fetcher.New(tgClient)

	// processor = processor.New(tgClient)

	// consumer.Start()
	// получает события (fetcher) и обрабатывает их (processor)
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
