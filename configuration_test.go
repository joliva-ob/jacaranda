package main


import (

	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"

)




func TestLoadConfiguration(t *testing.T) {

	var filename = os.Getenv("CONF_PATH") + "/" + os.Getenv("ENV") + ".yml"
	config = LoadConfiguration(filename)

	assert.NotNil(t, config, "Configuration cannot be null.")

}


func TestGetAlerts(t *testing.T) {

	alerts := GetAlerts()

	assert.NotZero(t, len(alerts), "There is one alert loaded at least.")

}



func TestGetAlert(t *testing.T) {

	alert := GetAlert( "no-sales" )

	assert.NotNil(t, alert, "Error, we expected to have loaded a no-sales alert at least.")

}