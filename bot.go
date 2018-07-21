package main

import (
	"bytes"
	"strings"
	"sync"
	"time"

	"github.com/mheidinger/server-bot/checkers"
	clog "gopkg.in/clog.v1"
	telebot "gopkg.in/tucnak/telebot.v2"
)

// TelegramUsers contains all telegram user ids of authorized users
var TelegramUsers []int
var noAuthCommands = [...]string{"/start", "/help"}
var markdownOptions = &telebot.SendOptions{ParseMode: telebot.ModeMarkdown}

// StartBot creates and starts the telegram bot; Blocking while bot runs!
func StartBot(telegramToken, botSecret string, results map[string]*checkers.CheckResult, mutex *sync.Mutex, notificationChannel chan *checkers.CheckResult) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token: telegramToken,
	})

	poller := &telebot.LongPoller{Timeout: 10 * time.Second}
	authPoller := telebot.NewMiddlewarePoller(poller, func(upd *telebot.Update) bool {
		if upd.Message == nil {
			return true
		}

		foundUser := -1
		for _, user := range TelegramUsers {
			if user == upd.Message.Sender.ID {
				foundUser = user
				break
			}
		}

		isNoAuthCommand := false
		for _, command := range noAuthCommands {
			if strings.HasPrefix(upd.Message.Text, command) {
				isNoAuthCommand = true
				break
			}
		}

		if foundUser == -1 && !isNoAuthCommand {
			if upd.Message.Text == botSecret {
				addUser(upd.Message.Sender.ID)
				_, err := bot.Send(upd.Message.Sender, "Correct password! 🐬")
				if err != nil {
					clog.Error(0, "Failed to send correct password message: %v", err)
				}
				return false
			}

			_, err := bot.Send(upd.Message.Sender, "⛔ *You are not authorized!* ⛔\nEnter the correct password to gain access!", markdownOptions)
			if err != nil {
				clog.Error(0, "Failed to send unauthorized message: %v", err)
			}
			return false
		}

		return true
	})

	bot.Poller = authPoller

	if err != nil {
		clog.Fatal(0, "Error setting up the telegram bot: %v", err)
	}

	bot.Handle("/start", func(m *telebot.Message) {
		bot.Send(m.Sender, "Welcome to the server-bot! 🎉\nFirst unlock the bot with the correct password and then try /help for all commands 🐬")
	})

	bot.Handle("/overview", func(m *telebot.Message) {
		buffer := bytes.NewBufferString("Last results of your services:\n")

		mutex.Lock()
		for _, res := range results {
			buffer.WriteString(res.Service.Name)
			buffer.WriteString(":\n")

			buffer.WriteString("\tStatus: ")
			if res.Success {
				buffer.WriteString("✔️\n")
			} else {
				buffer.WriteString("❌\n")
			}

			buffer.WriteString("\tChecker: ")
			buffer.WriteString(res.Service.CheckerName)
			buffer.WriteString("\n")

			buffer.WriteString("\tLast Check: ")
			buffer.WriteString(res.TimeStamp.String())
			buffer.WriteString("\n")

			buffer.WriteString("======================\n")
		}
		mutex.Unlock()

		_, err := bot.Send(m.Sender, buffer.String())
		if err != nil {
			clog.Error(0, "Failed to send overview message: %v", err)
		}
	})

	bot.Handle(telebot.OnText, func(m *telebot.Message) {
		_, err := bot.Send(m.Sender, "Unknown command 😱\nTry /help to list the best features 🐬")
		if err != nil {
			clog.Error(0, "Failed to send unknown command message: %v", err)
		}
	})

	go func() {
		for result := range notificationChannel {
			buffer := bytes.NewBufferString("Check for ")
			buffer.WriteString(result.Service.Name)
			if result.Success {
				buffer.WriteString(" succeeded ✔️\n")
			} else {
				buffer.WriteString(" failed ❌\n")
			}
			buffer.WriteString("Get more info on /overview 🐬")
			message := buffer.String()

			for _, user := range TelegramUsers {
				_, err := bot.Send(&telebot.User{ID: user}, message)
				if err != nil {
					clog.Error(0, "Failed to send notification message to %v: %v", user, err)
				}
			}
		}
	}()

	bot.Start()

}

func addUser(user int) {
	TelegramUsers = append(TelegramUsers, user)
	WriteConfig()
}
