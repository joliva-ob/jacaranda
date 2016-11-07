package main


import (

	"time"
	"strings"

	"github.com/tucnak/telebot"


)


var bot *telebot.Bot

const (
	HELP = "/help"
	START = "/start"
	STOP = "/stop"
	LIST = "/list"
	STATUS = "/status"
	EXEC = "/exec"
)



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

	messageChan := make(chan telebot.Message)
	bot.Listen(messageChan, 1*time.Second)


	for {
		select {
		case message := <- messageChan:
			processMessage( message )
		}
	}

}



/*
 Processing the chat messages captured by the bot and
 answer over the chat.
 */
func processMessage( message telebot.Message )  {

	words := strings.Fields(message.Text)
	var rule *RuleType

	if len(words) > 0 {

		if len(words) > 1 {
			rule = GetAlert(words[1])
		}

		switch words[0] {
		case HELP:
			bot.SendMessage(message.Chat, "version 1.1.6\nBot commands available are:\n/help\n/list\n/start {alert_name}\n/stop {alert_name}\n/status", nil)
			log.Info("/help requested from Chat ID: %v", message.Chat.ID)
		case LIST:
			alist := GetAlerts()
			bot.SendMessage(message.Chat, "Alert rules available are:\n" + alist, nil)
			log.Infof("/list requested from Chat ID: %v", message.Chat.ID)
		case START:
			processAndNotifyWatchdogChange(message, rule, START)
		case STOP:
			processAndNotifyWatchdogChange(message, rule, STOP)
		case STATUS:
			getCurrentStatus(message)
//		case EXEC:
//			execCommandLine(words,message)
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


