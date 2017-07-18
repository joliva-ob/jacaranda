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
+ alert when timeout reached on getting ELK metrics
+ re-start an alert previously disabled manually after 24 hour.
+ register to eureka and its heartbeat
+ Dockerfile + captain for Onebox Jenkins
+ Double check /info and /health for onebox requirements
+ Double check environment vars for jenkins + kubernetes
+ Use vendors to ensure version dependencies


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

1.1.4.2
+ status debug from inside goroutine watchdog

1.1.4.3
+ added alarm for response time 3rd party integrations

1.1.5
+ Added timeout on elasticsearch queries
+ improved security on elasticsearch queries

1.1.6
+ Improved alerting in case of elk timeout

1.1.7
+ Added capabilities to talk with Pod-Doublecheck from Bot

1.1.18
+ Updated response-time limit to 400 ms.

release/4.0.8
+ configurations for new monitoring in kubernetes cluster

release/4.0.16
+ Added detailed instance info for response time alert
+ Added markdown options over telegram send messages