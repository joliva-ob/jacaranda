package main

import (

	"net/http"
	"os"
	"fmt"

	"github.com/gorilla/mux"

)


/**
 * Main command to load configuration by given environment argument
 * and start application server to listen the exposed endpoints and
 * provide the requested resources operations
 *
 * Mandatory parameters are path (/tmp...) and environment (dev, qa, pre, pro...)
 */
func main() {

	// Load configuration to start application
	checkParams( os.Args )
	var filename = os.Getenv("CONF_PATH") + "/" + os.Getenv("ENV") + ".yml"
	config = LoadConfiguration(filename)
	InitializeTelegramBot() // Create and initialize the bot
	ec, i := registerToEureka()
	go sendHeartBeatToEureka(ec, i)
	log.Infof("jacaranda started with environment: %s and listening in port: %v\n", os.Getenv("ENV"), config.Server_port)


	// Starting server on given port number and listen for a chat conversation
	go ListenQueryChatMessages()


	// Starting alerts whatchdog
	startAlertsWatchdogs()


	// Create the router to handle requests
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/jacaranda/1.0/chats", ChatsController)
	router.HandleFunc("/jacaranda/1.0/sendMessage", SendMessagesController)
	router.HandleFunc("/jacaranda/1.0/info", InfoController)
	router.HandleFunc("/jacaranda/1.0/health", HealthController)
	log.Fatal( http.ListenAndServe(":" + config.Server_port, router) )

}




// Check the arguments to launch the application
// and provide specifications if needed.
func  checkParams(  args []string ) {

	if os.Getenv("CONF_PATH") == "" || os.Getenv("ENV") == "" {

		if len(args) < 2 {

			fmt.Println("ERROR: invalid parameters!")
			fmt.Println("Usage:")
			fmt.Println("./jacaranda [path-to-config-files] [environment] or export CONF_PATH export ENV")
			os.Exit(0)
		}
	}

}


