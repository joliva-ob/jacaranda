# jacaranda 
Golang alert api for Telegram messaging over elasticsearch index metrics exposed to kibana, tho goal is to have a chat where one of the guests will be a bot.
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
+ docker build --file=Dockerfile -t hub.docker.com/r/joanoj/jacaranda .
+ docker push hub.docker.com/r/joanoj/jacaranda



## TODO list
+ get goroutine status from the goroutine, not from the global status list
+ re-start an alert previously disabled manually after 24 hour.
+ register to eureka and its heartbeat
+ Dockerfile + captain for Onebox Jenkins
+ Double check /info and /health for onebox requirements
+ Double check environment vars for jenkins + kubernetes


## DONE list
1.0.0
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

1.1.0
+ added /status bot command

1.1.1
+ bug fixing

1.1.2
+ fixed precision from alerting values

1.1.3
+ exec command-line commands

1.1.4
+ command line bug fixing

1.1.4.1
+ command line disabled