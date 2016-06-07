package main


import (

	"encoding/json"
	"time"
	"strings"
	"strconv"

	"github.com/mattbaird/elastigo/lib"

	"errors"
)

const (
	BELOW  =  "below"
	ABOVE  =  "above"
	ENABLED = "enabled"
	DISABLED = "disabled"
)


var (

	elk_conn *elastigo.Conn

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

		go watchdogRoutine( &config.Rules[i] )

		log.Infof("Watchdog started for rule: %v", (config.Rules[i]).Alert_name)
	}

}



/*
 Change status from a given alert rule
 */
func ManageWatchdog( rule *RuleType, action string ) error {

	if rule != nil {
		switch action {
		case START:
			rule.Alert_status = ENABLED
		case STOP:
			rule.Alert_status = DISABLED
		}
	} else {
		return errors.New("Error starting rule, does not exist.")
	}

	return nil
}



/*
 Goroutine to keep listening a metric and rise a message
 in case of threashold were reached.
 get the query directly from consolemonit.metric index and eval the rule into elk (avg, elapsed time, count, whatever...)
 as the aggregation to count an absolute number (aggregations.count)
 http://pre.consolemonit1.oneboxtickets.com:9200/_plugin/marvel/sense/index.html
 */
func watchdogRoutine( rule *RuleType ) {

	// Open connection to elasticsearch and keep it
	elk_host := rule.Elk_host
	elk_conn = elastigo.NewConn()
	elk_conn.Domain = elk_host
	ticker := time.Tick(time.Duration(rule.Check_time_sec * 1000) * time.Millisecond)
	log.Infof("Watchdog [%s]--> Elasticsearch connected to host: %v", rule.Alert_name, rule.Elk_host)

	for {
		select {
		case <- ticker:
			if rule.Alert_status == ENABLED && isTimeWindow(rule.Time_window) {
				processRule( rule, elk_conn )
			}
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
func  processRule( rule *RuleType, elk_conn *elastigo.Conn  )  {

	// retrieve data from index
	args := make(map[string]interface{})
	args["size"] = 1
	args["from"] = 0
	lte := time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
	duration := int64(rule.Time_frame_sec) * 1000
	gte := lte - duration
	rule.Elk_filter = strings.Replace(rule.Elk_filter, "$lte", strconv.FormatInt(lte, 10), -1)
	rule.Elk_filter = strings.Replace(rule.Elk_filter, "$gte", strconv.FormatInt(gte, 10), -1)
//	log.Debugf("RuleName: %v --> gte: %v lte: %v query: ", rule.Alert_name, strconv.FormatInt(gte, 10), strconv.FormatInt(lte, 10), rule.Elk_filter)

	// Query elasticsearch
	out, err := elk_conn.Search(rule.Elk_index, "", args, rule.Elk_filter)
//	log.Debugf("RuleName: %v --> out: %v", rule.Alert_name, out.String(), string(out.RawJSON[:]))
	if out.Hits.Total >= rule.Min_items {

		var res = new (ElkAggregationsResponse)
		if err := json.Unmarshal(out.Aggregations, &res); err != nil {
			log.Error(err)
		}

		evaluateResponse( res, rule )

	}
	if err != nil {
		log.Errorf("[%v] Error occurred while trying to retrieve elasticsearch data: %v", rule.Alert_name, err)
	}

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

		alert_message := "Alert: rule " + rule.Alert_name + " " + strconv.FormatFloat(res.Count.Value, 'f', 6, 64) + " is " + rule.Raise_Condition + " than threshold " + strconv.FormatFloat(rule.Threshold, 'f', 6, 64)
		err := sendTelegramMessage( rule.Telegram_chat_id, alert_message )
		if err != nil {
			log.Error(err)
		}

		log.Warning(alert_message)
	}

}