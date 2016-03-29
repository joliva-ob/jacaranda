package main


import (

	"net/http"
	"encoding/json"
	"time"

)


const (

	JOAN_CHAT_ID int64 = 146665083
	ONEBOX_ALERT_CHAT_ID int64 = -126617985
	GROUP_TEST_CHAT_ID int64 = -117924915

)


// Chats response struct
type ChatsResponseType struct {

	Chats []*ChatResponseType `json:"chats"`
}

type ChatResponseType struct {

	Id int64 `json:"id"`
	Description string `json:"description"`
}



/**
 * Chats resource endpoint
 */
func ChatsController(w http.ResponseWriter, request *http.Request) {

	uuid := GetUuid()
	log.Infof( "{%v} /chats request %v received from: %v", uuid, request.URL, getIP(w, request) )
	start := time.Now()

	// Check authorization
	if !Authorize( request.Header.Get("Authorization") ) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Warningf("/chats error status 401 unauthorized.")
		return
	}

	// Set json response struct
	var chats []*ChatResponseType
	chat := new(ChatResponseType)
	chat.Id = JOAN_CHAT_ID
	chat.Description = "Joan Oliva TEST Telegram Chat"
	chats = append(chats, chat)
	chat = new(ChatResponseType)
	chat.Id = ONEBOX_ALERT_CHAT_ID
	chat.Description = "Onebox Alert chat"
	chats = append(chats, chat)
	chat = new(ChatResponseType)
	chat.Id = GROUP_TEST_CHAT_ID
	chat.Description = "Group Test chat"
	chats = append(chats, chat)

	chatsResponse := new(ChatsResponseType)
	chatsResponse.Chats = chats

	chatsjson, err := json.Marshal(chatsResponse)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		log.Errorf("/chats error status 204 no content.")
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Set response body
	w.Write(chatsjson)

	elapsed := time.Since(start)
	log.Infof( "{%v} /chats response status 200 in %v", uuid, elapsed )

}
