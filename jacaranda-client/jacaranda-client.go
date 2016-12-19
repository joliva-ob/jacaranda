package jacaranda_client



import (
	"net/http"
)



// move it into jacaranda project as a jacaranda-client.go library
func SendTelegramMessage( url string, text string, chatId string ) (*http.Response, error) {

	client := &http.Client{}
	//url := "http://jacaranda:8001/jacaranda/1.0/sendMessage"
	//url := "http://10.1.2.173:30002/jacaranda/1.0/sendMessage"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bear 1736cc7f-7c60-4576-b851-b7b3630cfeab")
	q := req.URL.Query()
	q.Add("chat_id", chatId)
	q.Add("text",text)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
