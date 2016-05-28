package main


import (

	"flag"

	"github.com/mattbaird/elastigo/lib"

)


var (

	elk_conn *elastigo.Conn

)

/*
 Starts a go routine for each rule from configuration file
 and keep listening metrics to rise an alert message if
 threashold is reached.
 */
func startAlertsWatchdogs() {

	for i := 0; i<len(config.Rules); i++ {

		rule := config.Rules[i]
		go watchdogRoutine( rule )

		log.Info("Watchdog started for rule: %v", rule.Alert_name)
	}

}



/*
 Goroutine to keep listening a metric and rise a message
 in case of threashold were reached.
 */
func watchdogRoutine( rule RuleType ) {

	// Open connection to elasticsearch
	elk_host := flag.String(rule.Elk_name, rule.Elk_value, rule.Elk_usage)
	elk_conn = elastigo.NewConn()
	flag.Parse()
	elk_conn.Domain = *elk_host
	args := make(map[string]interface{})
	args["size"] = 1 // the metric is forced to an absolute number
	args["from"] = 0
	log.Infof("Watchdog [%s]--> Elasticsearch connected to host: %v", rule.Alert_name, rule.Elk_name)


	for {

		// retrieve data from index
		out, err := elk_conn.Search(rule.Elk_index, "", args, rule.Elk_filter)
		if out.Hits.Len() > 0 {

			// evaluate the raise condition
			/*
			 TODO get the query directly from consolemonit.metric index and eval the rule into elk (avg, elapsed time, count, whatever...)
			 TODO start from GET .kibana/visualization/_search and get the searchSourceJSON
			 TODO http://pre.consolemonit1.oneboxtickets.com:9200/_plugin/marvel/sense/index.html
			  */

		}
		if err != nil {
			log.Errorf("{%v} Error occurred while trying to retrieve elasticsearch data: %v", rule.Alert_name, err)
		}



		// send alert message if needed

	}

}
