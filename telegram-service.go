package main


import (

	"github.com/tucnak/telebot"

)


var bot *telebot.Bot


func InitializeTelegramBot() bool {

	var err error
	bot, err = telebot.NewBot("182636765:AAFHu8FhAK1KPgad3hfqov8JNgqIZvWzLF0")
	if err != nil {
		return false
	}

	return true
}





