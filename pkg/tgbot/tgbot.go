package tgbot

import (
	"fmt"
	"log"
	"runtime/debug"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type TgInfoBot struct {
	logger *zap.Logger

	tgbot *tgbotapi.BotAPI
}

func NewTgInfoBot(logger *zap.Logger, tgbot *tgbotapi.BotAPI) *TgInfoBot {
	return &TgInfoBot{logger: logger, tgbot: tgbot}
}

func (self *TgInfoBot) Listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 120
	updates := self.tgbot.GetUpdatesChan(u)

	fmt.Println("bot is started .....")
	for update := range updates {
		log.Println("incoming update", update)
		go self.Handle(update)
	}
}

func (self *TgInfoBot) Handle(update tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("PANIC:", err, "stack:", string(debug.Stack()))
		}
	}()

	if update.Message != nil {
		if !self.tgbot.IsMessageToMe(*update.Message) {
			return
		}

		if update.Message.IsCommand() {
			cmd := update.Message.Command()
			if cmd == "current" {
				fmt.Print("%#s", update)
				// self.HandleAnalyzeCommand(update)
			}
		}
	}
}

func NewMessage(chatId int64, replyId int, txt string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, txt)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if replyId != 0 {
		msg.ReplyToMessageID = replyId
	}
	return msg
}

func (self *TgInfoBot) reply(update tgbotapi.Update, msg string) {
	self.tgbot.Send(NewMessage(update.Message.Chat.ID, update.Message.MessageID, msg))
}

func (self *TgInfoBot) HandleAnalyzeCommand(update tgbotapi.Update) {
	self.reply(update, "Смотрим ! Немного терпения ...")

}
