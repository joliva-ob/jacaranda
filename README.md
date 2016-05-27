# jacaranda 
Golang alerting api for Telegram messaging, tho goal is to have a chat where one of the guests will be a bot.
Hence, this bot will be able to send alert messages to the group and also be able to be listening the conversation 
in order to answer some technical questions. (ie: /status, and anything starting by a slash).

Compiled with runtime with: 
+ GOOS=windows GOARCH=386 go build -o jacaranda.exe jacaranda.go
+ GOOS=linux GOARCH=386 go build -o jacaranda.linux jacaranda.go
+ GOOS=darwin GOARCH=386 go build -o jacaranda jacaranda.go

Build Docker image with
+ cp /source_cfg_files/*env* .
+ docker build -f docker/Dockerfile . -tag jacaranda 
+ docker run --publish 8000:8000 --name jacaranda --rm jacaranda --restart=always jacaranda 



## TODO list
+ make bot to be listenig the chat and answer basic questions (/status, /who is down, ...)
+ manage chats and people from chat
+ unit tests
+ dockerize


## Optional TODO list



## DONE list
+ load configuration from yml file
+ register to eureka
+ send message to chat
