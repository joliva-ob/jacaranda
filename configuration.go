package main

import (
	"io/ioutil"
	"os"
	"fmt"
	"bytes"

	"gopkg.in/yaml.v2"
	"github.com/op/go-logging"
)


// Global vars
var config ConfigType
var log = logging.MustGetLogger("jacaranda")
var alertsMap map[string] *RuleType
var alertsList string



type RuleType struct{

	Alert_name       string
	Alert_status	 string
	Telegram_chat_id int64
	Elk_index        string
	Elk_host	 string
	Threshold        float64
	Raise_Condition  string
	Time_window_utc  string
	Time_frame_sec   int64
	Check_time_sec	 int64
	Min_items   	 int
	Elk_filter       string
}


// Instance configuration
type ConfigType struct {

	Server_port string
	Log_file string
	Log_format string
	Eureka_port int
	Eureka_ip_addr string
	Eureka_app_name string
	Telegram_bot_token string
	Rules[] RuleType
}



/**
 * Load configuration yaml file
 */
func LoadConfiguration(filename string) ConfigType {

	// Set config
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("--> Configuration loaded values: %#v\n", config)


	// Set logger
	format := logging.MustStringFormatter( config.Log_format )
	logbackend1 := logging.NewLogBackend(os.Stdout, "", 0)
	logbackend1Formatted := logging.NewBackendFormatter(logbackend1, format)
	f, err := os.OpenFile(config.Log_file, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		defer f.Close()
	}
	logbackend2 := logging.NewLogBackend(f, "", 0)
	logbackend2Formatted := logging.NewBackendFormatter(logbackend2, format)
	logging.SetBackend(logbackend1Formatted, logbackend2Formatted)

	return config
}




func GetAlerts() string {

	var buffer bytes.Buffer

	for i := 0; i<len(config.Rules); i++ {

		buffer.WriteString(config.Rules[i].Alert_name)
		buffer.WriteString(" - ")
		buffer.WriteString(config.Rules[i].Alert_status)
		buffer.WriteString("\n")
	}

	return buffer.String()
}




func GetAlert( alertName string ) *RuleType {

	if len(alertsMap) == 0 {

		alertsMap = make(map[string] *RuleType)

		for i := 0; i<len(config.Rules); i++ {
			alertsMap[config.Rules[i].Alert_name] = &config.Rules[i]
		}

	}

	return alertsMap[alertName]

}