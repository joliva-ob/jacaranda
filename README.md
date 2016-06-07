# jacaranda 
Golang alerting api for Telegram messaging, tho goal is to have a chat where one of the guests will be a bot.
Hence, this bot will be able to send alert messages to the group and also be able to be listening the conversation 
in order to answer some technical questions. (ie: /status, and anything starting by a slash).


Run unit tests
+ go test

Compiled with runtime with: 
+ GOOS=windows GOARCH=386 go build -o jacaranda.exe
+ GOOS=linux GOARCH=386 go build -o jacaranda.linux
+ GOOS=darwin GOARCH=386 go build -o jacaranda

Build Docker image with
+ cp /source_cfg_files/*env* .
+ docker build -f docker/Dockerfile . -tag jacaranda 
+ docker run --publish 8000:8000 --name jacaranda --rm jacaranda --restart=always jacaranda 

Kubernetes
+ docker build --file=Dockerfile -t docker-registry.oneboxtickets.com/oneboxtm/jacaranda .
+ docker push docker-registry.oneboxtickets.com/oneboxtm/jacaranda



## TODO list


## DONE list
+ load configuration from yml file
+ send message to chat
+ listening telegram bot
+ start alerts from configuration file
+ + control the go routines to stop, start from config
+ change sleep for ticker + channel
+ only process alert between time_window config
+ start/stop from endpoint/bot
+ dockerization
+ check unit tests coverage