package main

import (
	"flag"
	"log"
	"proj/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.New(tgBotHost, token())
	//fetcher - собиратель - отправляет запросы API телеги, чтобы получить новые события
	//fetcher = fetcher.New()
	//обработка сообщений и выполнение каких-либо действий
	//processor = processor.New()

	//consumer - потребитель - получение и обрпботка событий при помощи fetcher и processor
	// consumer.Start()
}

func token() string {
	t := flag.String(
		"t-bot-t",
		"",
		"token to acess tg bot",
	)
	flag.Parse()

	if *t == "" {
		log.Fatal("no token")
	}
	return *t
}
