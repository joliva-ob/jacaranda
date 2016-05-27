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
	var filename = os.Args[1] + "/" + os.Args[2] + ".yml"
	config = LoadConfiguration(filename)
	log = GetLog()
	InitializeTelegramBot() // Create and initialize the bot
	log.Infof("alertigo started with environment: %s and listening in port: %v\n", os.Args[2], config.Server_port)


	// Starting server on given port number and listen for a chat conversation
	go ListenQueryChatMessages()


	// Starting alerts whatchdog
	startAlertsWatchdogs()


	// Create the router to handle requests
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/alertigo/1.0/chats", ChatsController)
	router.HandleFunc("/alertigo/1.0/sendMessage", SendMessagesController)
	router.HandleFunc("/alertigo/1.0/info", InfoController)
	router.HandleFunc("/alertigo/1.0/health", HealthController)
	log.Fatal( http.ListenAndServe(":" + config.Server_port, router) )

}




// Check the arguments to launch the application
// and provide specifications if needed.
func  checkParams(  args []string ) {

	if len(args) < 2 {

		fmt.Println("ERROR: invalid arguments number!")
		fmt.Println("Usage:")
		fmt.Println("./alertigo [path-to-config-files] [environment]")
		os.Exit(0)
	}

}


