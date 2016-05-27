package main








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

	// retrieve data from index

	// evaluate the raise condition

	// send alert message if needed

}
