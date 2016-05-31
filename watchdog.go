package main


import (

	"encoding/json"
	"time"
	"strings"
	"strconv"

	"github.com/mattbaird/elastigo/lib"

)

const (
	BELOW  =  "below"
	ABOVE  =  "above"
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

		rule := config.Rules[i]
		go watchdogRoutine( rule )

		log.Infof("Watchdog started for rule: %v", rule.Alert_name)
	}

}



/*
 Goroutine to keep listening a metric and rise a message
 in case of threashold were reached.
 get the query directly from consolemonit.metric index and eval the rule into elk (avg, elapsed time, count, whatever...)
 as the aggregation to count an absolute number (aggregations.count)
 http://pre.consolemonit1.oneboxtickets.com:9200/_plugin/marvel/sense/index.html
 */
func watchdogRoutine( rule RuleType ) {

	// Open connection to elasticsearch
	elk_host := rule.Elk_value
	elk_conn = elastigo.NewConn()
	elk_conn.Domain = elk_host
	args := make(map[string]interface{})
	args["size"] = 1
	args["from"] = 0
	log.Infof("Watchdog [%s]--> Elasticsearch connected to host: %v", rule.Alert_name, rule.Elk_host)


	// TODO control goroutine life cycle by a channel and let the bot be able to handle it
	// TODO control time between two consecutive alertes, do not spamming!
	for {

		// retrieve data from index
		lte := time.Now().UnixNano() / int64(time.Millisecond)
		duration := int64(rule.Time_frame_sec) * int64(time.Millisecond)
		gte := lte - duration
		rule.Elk_filter = strings.Replace(rule.Elk_filter, "$lte", strconv.FormatInt(lte, 10), -1)
		rule.Elk_filter = strings.Replace(rule.Elk_filter, "$gte", strconv.FormatInt(gte, 10), -1)

		out, err := elk_conn.Search(rule.Elk_index, "", args, rule.Elk_filter)
		if out.Hits.Total >= rule.Min_items {

//			log.Debugf("Total hits: %v aggregations: %v", out.Hits.Total, string(out.Aggregations))
			var res = new (ElkAggregationsResponse)
			if err := json.Unmarshal(out.Aggregations, &res); err != nil {
				log.Error(err)
			}

			evaluateResponse( res, rule )

		}
		if err != nil {
			log.Errorf("[%v] Error occurred while trying to retrieve elasticsearch data: %v", rule.Alert_name, err)
		}

		// Check rule every N seconds
		time.Sleep(time.Duration(rule.Time_frame_sec * 1000) * time.Millisecond)
	}

}




/*
 Evaluate the raise condition and send the message if needed.
 */
func evaluateResponse( res *ElkAggregationsResponse, rule RuleType ) {

	isRaised := false


	if rule.Raise_Condition == BELOW {

		if res.Count.Value < rule.Threshold {
			log.Debugf("count: %v threshold: %v", res.Count.Value, rule.Threshold)
			isRaised = true
		}

	} else if rule.Raise_Condition == ABOVE {

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

		log.Info(alert_message)
	}

}