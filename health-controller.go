package main


import (

	"net/http"
	"time"
	"encoding/json"
)


// Health response struct
type HealthResponseType struct {

	Status string `json:"status"`
}



/**
 * Health resource endpoint
 */
func HealthController(w http.ResponseWriter, request *http.Request) {

	uuid := GetUuid()
	log.Infof( "{%v} /health request %v received from: %v", uuid, request.URL, getIP(w, request) )
	start := time.Now()

	// Set json response struct
	var healthresponse HealthResponseType
	healthresponse.Status = STATUS_UP
	// TODO fill the discovery and other resources statuses
	healthjson, _ := json.Marshal(healthresponse)

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Set response body
	w.Write(healthjson)

	elapsed := time.Since(start)
	log.Infof( "{%v} /health response status 200 in %v", uuid, elapsed )

}
