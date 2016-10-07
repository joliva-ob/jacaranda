package main


import (

	"os/exec"
	"bytes"

	"github.com/tucnak/telebot"
)


func execCommandLine( words []string, message telebot.Message ) string {

	log.Infof("/exec %v requested from Chat ID: %v", words[1], message.Chat.ID)

	var cmd *exec.Cmd
	switch len(words) {
	case 2:
		cmd = exec.Command(words[1])
	case 3:
		cmd = exec.Command(words[1], words[2])
	case 4:
		cmd = exec.Command(words[1], words[2], words[3])
	case 5:
		cmd = exec.Command(words[1], words[2], words[3], words[4])
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		log.Error(err.Error())
		bot.SendMessage(message.Chat, err.Error(), nil)
	}
	err = cmd.Wait()

	bot.SendMessage(message.Chat, out.String(), nil)
	log.Infof("Command %v finished. %v", words[1], err.Error())

	return out.String()
}
