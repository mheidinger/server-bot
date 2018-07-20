package main

import (
	"log"
	"sync"
	"time"

	"github.com/mheidinger/server-bot/checkers"
	telebot "gopkg.in/tucnak/telebot.v2"
)

// StartBot creates and starts the telegram bot; Blocking while bot runs!
func StartBot(telegramToken, botSecret string, results map[string]*checkers.CheckResult, mutex *sync.Mutex) {
	poller := &telebot.LongPoller{Timeout: 10 * time.Second}
	authPoller := telebot.NewMiddlewarePoller(poller, func(upd *telebot.Update) bool {
		return true
	})

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  telegramToken,
		Poller: authPoller,
	})

	if err != nil {
		log.Fatalf("Error setting up the telegram bot: %v", err)
	}

	bot.Handle("/start", func(m *telebot.Message) {
		bot.Send(m.Sender, "Welcome to the server-bot! 🎉\nFirst unlock the bot with the correct password and then try /help for all commands 😁")
	})

	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		if m.Text == botSecret {
			bot.Send(m.Sender, "Correct password!")
		} else {
			bot.Send(m.Sender, "⛔ Wrong password! ⛔")
		}
	})

	bot.Start()
}
