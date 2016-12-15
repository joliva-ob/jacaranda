package main

import (


	"net/http"
	"testing"
	"strconv"
	"github.com/stretchr/testify/assert"

	"fmt"
)



func TestChangePodDoubleChecktime(t *testing.T) {

	go startServer(9090)

	params := make(map[string]string)
	params["time"] = strconv.Itoa(20)
	// res, err := testSendHttpRequest("PUT", "http://10.1.22.50:30920/refreshtime?", params, nil)
	res, err := testSendHttpRequest("PUT", "http://localhost:9090/refreshtime?", params, nil)

	if res != nil {
		assert.NotNil(t, res, "Http response is %v", res)
	} else {
		assert.Nil(t, res, "Http response time is nil")
	}

	if err != nil {
		assert.NotNil(t, err, "Error changing pod-doublecheck refreshtime: %v", err.Error())
	}
}




func testSendHttpRequest(method string, url string, params map[string]string, headers map[string]string) (*http.Response, error) {

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, nil)
	//	req.Header.Set("Authorization", "Bear 1736cc7f-7c60-4576-b851-b7b3630cfeab")
	q := req.URL.Query()
	q.Add("time", params["time"])
	req.URL.RawQuery = q.Encode()
	fmt.Printf("Url to send http request: %v", req.URL.RawQuery)

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}


func startServer(port int) {
	http.HandleFunc("/", fooHandler)
	http.ListenAndServe(":" + strconv.Itoa(port), nil)
}

func fooHandler(w http.ResponseWriter, r *http.Request){
	fmt.Printf("time param received is: " + r.URL.Query().Get("time"))
}