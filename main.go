package main

import (
	"flag"
	"log"
	tgClient "proj/clients/telegram"

	event_consumer "proj/consumer/event-consumer"
	"proj/events/telegram"
	files "proj/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	//fetcher - собиратель - отправляет запросы API телеги, чтобы получить новые события
	//fetcher = fetcher.New()
	//обработка сообщений и выполнение каких-либо действий
	//processor = processor.New()

	//consumer - потребитель - получение и обрпботка событий при помощи fetcher и processor
	// consumer.Start()

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, token()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func token() string {
	t := flag.String(
		"tg-bot-token",
		"",
		"token to acess tg bot",
	)
	flag.Parse()

	if *t == "" {
		log.Fatal("no token")
	}
	return *t
}
