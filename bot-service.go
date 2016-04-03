package main


import (

	"time"
	"fmt"

	"github.com/tucnak/telebot"

)


func ListenQueryChatMessages() {

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {

		fmt.Println("request message received: %v", message.Text)

		if message.Text == "/hi" {

			bot.SendMessage(message.Chat, "Hello , " + message.Sender.FirstName + "!", nil)
			fmt.Println("/hi requested from Chat ID: %v", message.Chat.ID)
		}

	}

}
