package main

import (

	"net/http"
	"net"

	"github.com/satori/go.uuid"
)


const (

	CHAT_ID = "chat_id"
	TEXT = "text"

	STATUS_UP = "UP"
	STATUS_DOWN = "DOWN"
	STATUS_OK = "OK"
	STATUS_ERROR = "ERROR"
	STATUS_SENT = "SENT"

)


// Parameters response struct
type ParametersResponseType struct {

	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	EventId int `json:event_id`
	Page int `json:"page"`
	TraceId string `json:"trace_id"`
}


// Global vars and default values
var chatId int
var text string



// https://blog.golang.org/context/userip/userip.go
// Funtion to retrieve the sender IP from request
// or from forwared headers instead
func getIP(w http.ResponseWriter, req *http.Request) string {

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Debugf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return req.RemoteAddr
	}

	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr Header.Get is case-insensitive
	forward := req.Header.Get("X-Forwarded-For")
	return forward
}



// Generate a universal unique identifier UUID
func GetUuid() string {

	// Creating UUID Version 4
	uuid1 := uuid.NewV4()

	return uuid1.String()
}
