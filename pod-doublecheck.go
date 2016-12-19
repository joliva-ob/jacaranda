package main

import (

	"net/http"
	"io/ioutil"
	"github.com/tucnak/telebot"
	"strconv"
)


type MonitoringType struct {
	Name string
	Value int64
	Threshold int64
	Alert bool
}


var Monitoring = make(map[string]*MonitoringType)


// get and send the status by requesting it to pod-doublecheck microservice
func processPodDoubleCheckStatus( param string, message *telebot.Message) {

	switch param {
	case "status":

		url := config.Pod_doublecheck_url+"monitoring"
		res, err := sendHttpRequest("GET", url, nil, nil)

		if err != nil {
			bot.SendMessage(message.Chat, "Error retrieving status from pod-doublecheck monitoring: "+err.Error(), nil)
			log.Infof("Error retrieving status from pod-doublecheck monitoring: %v", err.Error())
		} else {
			jsonData, _ := ioutil.ReadAll(res.Body)
			bot.SendMessage(message.Chat, "Pod-Doublecheck current status is: "+string(jsonData), nil)
			log.Infof("Pod-Doublecheck current status is: %v", res.Body)
		}
	}
}



func processNewPodDoublecheckRefreshtime(newRefreshtime int, message *telebot.Message) {

	url := config.Pod_doublecheck_url+"refreshtime?"
	params := make(map[string]string)
	params["time"] = strconv.Itoa(newRefreshtime)
	res, err := sendHttpRequest("PUT", url, params, nil)

	if err != nil {
		bot.SendMessage(message.Chat, "Error changing pod-doublecheck refreshtime to "+strconv.Itoa(newRefreshtime)+" -> "+err.Error(), nil)
		log.Infof("Error changing pod-doublecheck refreshtime to %v: %v", newRefreshtime, res.Status)
	} else {
		if newRefreshtime <=0 {
			bot.SendMessage(message.Chat, "Pod-Doublecheck is now stopped.", nil)
			log.Infof("Pod-Doublecheck is now stopped.")
		} else {
			bot.SendMessage(message.Chat, "Pod-Doublecheck refresh time changed to: "+strconv.Itoa(newRefreshtime), nil)
			log.Infof("Pod-Doublecheck refresh time changed to: %v", newRefreshtime)
		}

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