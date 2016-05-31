package main


import (

	"net/http"
	"encoding/json"
	"time"
	"strconv"

)



// Message response struct
type MessageResponseType struct {

	Status string `json:"status"`

}



/**
 * Messages resource endpoint
 */
func SendMessagesController(w http.ResponseWriter, request *http.Request) {

	sentStatus := STATUS_OK
	uuid := GetUuid()
	log.Infof( "{%v} /sendMessage request %v received from: %v", uuid, request.URL, getIP(w, request) )
	start := time.Now()

	// Check authorization
	if !Authorize( request.Header.Get("Authorization") ) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Warningf("/sendMessage error status 401 unauthorized.")
		return
	}

	// GET request params
	chatId, err := strconv.ParseInt(request.URL.Query().Get(CHAT_ID),10,64)
	text = request.URL.Query().Get(TEXT)

	// Send the message to the given chat id
	err = sendTelegramMessage( chatId, text )
	if err != nil {
		sentStatus = STATUS_ERROR
	}

	// Set json response struct
	var messageresponse MessageResponseType
	messageresponse.Status = sentStatus

	messagejson, error := json.Marshal(messageresponse)
	if error != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Errorf("/sendMessage error status 204 no content.")
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Set response body
	w.Write(messagejson)

	elapsed := time.Since(start)
	log.Infof( "{%v} /sendMessage response status 200 in %v", uuid, elapsed )

}
