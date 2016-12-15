package main

import (


	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

)



func TestChangePodDoubleChecktime(t *testing.T) {

	res, err := testSendHttpRequest("PUT", "http://10.1.22.50:30920/refreshtime?-1", nil, nil)

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
	//q.Add("", "-1")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}