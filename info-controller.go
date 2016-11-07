package main


import (

	"net/http"
	"time"
	"encoding/json"

)


// Info response struct
type InfoResponseType struct {

	Version string `json:"version"`
}



/**
 * Info resource endpoint
 */
func InfoController(w http.ResponseWriter, request *http.Request) {

	uuid := GetUuid()
	log.Infof( "{%v} /info request %v received from: %v", uuid, request.URL, getIP(w, request) )
	start := time.Now()

	// Set json response struct
	var inforesponse InfoResponseType
	inforesponse.Version = "1.1.6"
	// TODO fill the version, release and git branch
	infojson, _ := json.Marshal(inforesponse)

	// Set response headers and body
	w.Header().Set("Content-Type", "application/json")
	w.Write(infojson)

	elapsed := time.Since(start)
	log.Infof( "{%v} /info response status 200 in %v", uuid, elapsed )

}
