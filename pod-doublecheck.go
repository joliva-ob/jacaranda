package main

import (

	"net/http"

	"github.com/tucnak/telebot"
	"strconv"
)


func processNewPodDoublecheckRefreshtime(newRefreshtime int, message *telebot.Message) {

	url := config.Pod_doublecheck_url+strconv.Itoa(newRefreshtime)
	params := make(map[string]string)
	params["time"] = strconv.Itoa(newRefreshtime)
	res, err := sendHttpRequest("PUT", url, params, nil)

	if err != nil {
		bot.SendMessage(message.Chat, "Error changing pod-doublecheck refreshtime to "+strconv.Itoa(newRefreshtime)+" -> "+err.Error(), nil)
		log.Infof("Error changing pod-doublecheck refreshtime to %v: %v", newRefreshtime, res.Status)
	} else {
		bot.SendMessage(message.Chat, "Pod-Doublecheck refresh time changed to: "+strconv.Itoa(newRefreshtime), nil)
		log.Infof("Pod-Doublecheck refresh time changed to: %v", newRefreshtime)
	}
}


func sendHttpRequest(method string, url string, params map[string]string, headers map[string]string) (*http.Response, error) {

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, nil)
//	req.Header.Set("Authorization", "Bear 1736cc7f-7c60-4576-b851-b7b3630cfeab")
	q := req.URL.Query()
	q.Add("time", params["time"])
	req.URL.RawQuery = q.Encode()
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}