package main


import (

	"time"
	"strings"
	"strconv"
	"github.com/tucnak/telebot"
)


var bot *telebot.Bot

const (
	HELP = "/help"
	START = "/start"
	STOP = "/stop"
	LIST = "/list"
	STATUS = "/status"
	POD_DOUBLECHECK = "/pod-doublecheck"
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
			bot.SendMessage(message.Chat, "version release/4.0.16\nBot commands available are:\n/help\n/list\n/start {alert_name}\n/stop {alert_name}\n/status\n/pod-doublecheck {>0: new_refresh_time_sec | <=0: is to disable | status: is to get current status}", nil)
			log.Info("/help requested from Chat ID: %v", message.Chat.ID)
		case LIST:
			alist := GetAlerts()
			bot.SendMessage(message.Chat, "Alert rules available are:\n"+alist, nil)
			log.Infof("/list requested from Chat ID: %v", message.Chat.ID)
		case START:
			processAndNotifyWatchdogChange(message, rule, START)
		case STOP:
			processAndNotifyWatchdogChange(message, rule, STOP)
		case STATUS:
			getCurrentStatus(message)
			//		case EXEC:
			//			execCommandLine(words,message)
		case POD_DOUBLECHECK:
			processPodDoublecheck(words[1], &message)
		}
	}
}


func processPodDoublecheck(param string, message *telebot.Message) {

	if param != "" {
		newRefreshtime, err := getNewRefreshtime(param)
		if err != nil {
			processPodDoubleCheckStatus(param, message)
		} else {
			processNewPodDoublecheckRefreshtime(newRefreshtime, message)
		}
	} else {
		bot.SendMessage(message.Chat, "Param not allowed: "+param+". Valid ones are {status/int}", nil)
	}
}



func getNewRefreshtime(strTime string) (int, error) {
	newtime, err := strconv.Atoi(strTime)
	if err != nil {
		return -99, err
	} else {
		return newtime, nil
	}
}



/*
 Send a telegram message to the given chat Id
 */
func sendTelegramMessage( chatId int64, text string ) error {

	var options telebot.SendOptions
	options.ParseMode = "Markdown"
	var chat telebot.Chat
	chat.ID = chatId

	err := bot.SendMessage(chat, text, &options)

	return err
}


