package main


import (

	"time"

	"github.com/tucnak/telebot"


)


var bot *telebot.Bot


/*
 Initialize the bot for the given token loaded from the
 configuration file.
 */
func InitializeTelegramBot() bool {

	var err error
	bot, err = telebot.NewBot(config.Telegram_bot_token)
	if err != nil {
		return false
	}

	log.Infof("Telegram Bot initialized with token: %v", config.Telegram_bot_token)

	return true
}


/*
 Goroutine to keep listening a chat channel for any command
 the bot could answer
 */
func ListenQueryChatMessages() {

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {

		if message.Text == "/hi" {

			bot.SendMessage(message.Chat, "Hello , " + message.Sender.FirstName + "!", nil)
			log.Info("/hi requested from Chat ID: %v", message.Chat.ID)
		}

	}

}



/*
 Send a telegram message to the given chat Id
 */
func sendTelegramMessage( chatId int64, text string ) error {

	var chat telebot.Chat
	chat.ID = chatId
	err := bot.SendMessage(chat, text, nil)

	return err
}


