package main


import (

	"time"
	"os"

	"github.com/op/go-logging"
	"github.com/hudl/fargo"

)


// Global vars
// var config ConfigType
var log = logging.MustGetLogger("jacaranda")
var eurekaConn fargo.EurekaConnection


// Register and keep the eureka connection
func registerToEureka( ) (fargo.EurekaConnection, fargo.Instance) {

	eurekaConn = fargo.NewConn("http://pre.eureka1.oneboxtickets.com:8761")
	hostname, _ := os.Hostname()
	i := fargo.Instance{
		HostName:         hostname,
		Port:             8000,
		App:              "jacaranda",
		IPAddr:           hostname,
		SecureVipAddress: "10.1.121.199",
		VipAddress:       "jacaranda",
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.Amazon},
		Status:           fargo.UP,
		HealthCheckUrl:	  "http://" + hostname + ":" +config.Server_port+ "/jacaranda/1.0/health",
		StatusPageUrl:	  "http://" + hostname + ":" +config.Server_port+ "/jacaranda/1.0/health",
		HomePageUrl:      "http://" + hostname + ":" +config.Server_port+ "/jacaranda/1.0/health",
	}
	err := eurekaConn.RegisterInstance(&i)
	if err != nil {
		log.Error("%v", err)
	}

	return eurekaConn, i
}



// Go routine to keep registered into
// Eureka service discovery
func sendHeartBeatToEureka( ec fargo.EurekaConnection, i fargo.Instance ) {

	ticker := time.Tick(time.Duration(30 * 1000) * time.Millisecond)

	for {
		select {
		case <- ticker:
			ec.HeartBeatInstance(&i)
		}
	}
}
