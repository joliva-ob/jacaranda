package main


import (

	"encoding/json"
	"time"
	"strings"
	"strconv"
	"errors"

	"github.com/mattbaird/elastigo/lib"
	"github.com/tucnak/telebot"
)

const (
	BELOW  =  "below"
	ABOVE  =  "above"
	ENABLED = "enabled"
	DISABLED = "disabled"
	CHECK = "check"
	EVALUATE = "evaluate"
)


var (

	elk_conn *elastigo.Conn
	statusChan = make(chan string)

)



type ElkAggregationsResponse struct {

	Count struct {
		      Value float64 `json:"value"`
	      } `json:"count"`
}



/*
 Starts a go routine for each rule from configuration file
 and keep listening metrics to rise an alert message if
 threashold is reached.
 */
func startAlertsWatchdogs() {

	for i := 0; i<len(config.Rules); i++ {

		go watchdogRoutine( &config.Rules[i], statusChan )

		log.Infof("Watchdog started for rule: %v", (config.Rules[i]).Alert_name)
	}

}



/*
 Change status from a given alert rule
 */
func processAndNotifyWatchdogChange( message telebot.Message, rule *RuleType, action string ) error {

	if rule != nil {

		switch action {
		case START:
			rule.Alert_status = ENABLED
		case STOP:
			rule.Alert_status = DISABLED
		}

		// OK message
		bot.SendMessage(message.Chat, "Alert " + rule.Alert_name + " is now " + rule.Alert_status, nil)
		log.Infof("/%v %v requested from Chat ID: %v is now %v", action, rule.Alert_name, message.Chat.ID, rule.Alert_status)

	} else {
		// ERROR message
		bot.SendMessage(message.Chat, "Error stopping rule, does not exist.", nil)
		return errors.New("Error starting rule, does not exist.")
	}

	return nil
}



/*
 Goroutine to keep listening a metric and rise a message
 in case of a threshold were reached.
 get the query directly from consolemonit.metric index and eval the rule into elk (avg, elapsed time, count, whatever...)
 as the aggregation to count an absolute number (aggregations.count)
 http://pre.consolemonit1.oneboxtickets.com:9200/_plugin/marvel/sense/index.html
 */
func watchdogRoutine( rule *RuleType, statusChan chan string) {

	// Open connection to elasticsearch and keep it
	elk_host := rule.Elk_host
	elk_conn = elastigo.NewConn()
	elk_conn.Domain = elk_host
	ticker := time.Tick(time.Duration(rule.Check_time_sec * 1000) * time.Millisecond)
	var statusAction string
	log.Infof("Watchdog [%s]--> Elasticsearch connected to host: %v", rule.Alert_name, rule.Elk_host)

	for {
		select {
		case <- ticker:
			if rule.Alert_status == ENABLED && isTimeWindow(rule.Time_window_utc) {
				processRule( rule, elk_conn, EVALUATE )
			}
		case statusAction = <- statusChan:
			currentValue := processRule(rule, elk_conn, statusAction )
			log.Debugf("Watchdog rule %v is %v with current value of: %v", rule.Alert_name, rule.Alert_status, strconv.FormatFloat(currentValue, 'f', 0, 64))
		}
	}

}




func isTimeWindow( timeWindow string ) bool {

	isBetween := false

	now := time.Now()
	s := strings.Split(timeWindow, "-")
	from, _ := strconv.Atoi(s[0])
	to, _ := strconv.Atoi(s[1])

	if from <= now.Hour() && now.Hour() <= to {
		isBetween = true
	}

	return isBetween
}



/*
 Generic process to get data from a metric and evaluate func init() {
 on the rules defined by configuration.
 }
 */
func  processRule( rule *RuleType, elk_conn *elastigo.Conn, action string  ) float64 {

	// retrieve data from index
	args := make(map[string]interface{})
	args["size"] = 1
	args["from"] = 0
	args["timeout"] = rule.Elk_timeout
	lte := time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
	duration := int64(rule.Time_frame_sec) * 1000
	gte := lte - duration
	query := rule.Elk_filter
	query = strings.Replace(query, "$lte", strconv.FormatInt(lte, 10), -1)
	query = strings.Replace(query, "$gte", strconv.FormatInt(gte, 10), -1)
//	log.Debugf("RuleName: %v --> gte: %v lte: %v query: ", rule.Alert_name, strconv.FormatInt(gte, 10), strconv.FormatInt(lte, 10), query)

	// Query elasticsearch
	out, err := elk_conn.Search(rule.Elk_index, "", args, query)
	if err != nil {
		log.Errorf("[%v] Error occurred while trying to retrieve elasticsearch data: %v", rule.Alert_name, err)
		return 0
	}

	if out.Hits.Total >= rule.Min_items {

		res := new (ElkAggregationsResponse)
		if err := json.Unmarshal(out.Aggregations, &res); err != nil {
			log.Error(err)
		}

//		log.Debugf("RuleName: %v --> out: %v res addr: %v out addr: %v", rule.Alert_name, res.Count.Value, &res, &out)
		switch action {
		case EVALUATE:
			evaluateResponse( res, rule )
		case CHECK:
			return res.Count.Value
		}

	}

	return 0
}



/*
 Evaluate the raise condition and send the message if needed.
 */
func evaluateResponse( res *ElkAggregationsResponse, rule *RuleType ) {

	isRaised := false


	switch rule.Raise_Condition {
	case BELOW:
		if res.Count.Value < rule.Threshold {
			isRaised = true
		}
	case ABOVE:
		if res.Count.Value >= rule.Threshold {
			isRaised = true
		}

	}


	if isRaised {

		alert_message := "Alert: rule " + rule.Alert_name + " " + strconv.FormatFloat(res.Count.Value, 'f', 0, 64) + " is " + rule.Raise_Condition + " than threshold " + strconv.FormatFloat(rule.Threshold, 'f', 0, 64)
		err := sendTelegramMessage( rule.Telegram_chat_id, alert_message )
		if err != nil {
			log.Error(err)
		}

		log.Warning(alert_message)
	}

}



/*
 Retrieve the current monitoring kpi values
 */
func getCurrentStatus( message telebot.Message ) error {

	currentStatus := "Current status is:\n"

	for i := 0; i<len(config.Rules); i++ {

		statusChan <- CHECK
		value := processRule( &config.Rules[i], elk_conn, CHECK)
		alertName := config.Rules[i].Alert_name
		alertFrame := config.Rules[i].Time_frame_sec
		currentStatus = currentStatus + alertName + "\t" + strconv.FormatFloat(value, 'f', 0, 64) + " in " + strconv.FormatInt(alertFrame, 10) + " seconds.\n"

	}

	bot.SendMessage(message.Chat, currentStatus, nil)
	log.Infof("/status requested from Chat ID: %v", message.Chat.ID)

	return nil
}