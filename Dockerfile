#
# Dockerfile for Jacaranda alerting project
#

FROM golang

MAINTAINER Joan Oliva

RUN mkdir /jacaranda
RUN mkdir /jacaranda/bin
RUN mkdir /jacaranda/cfg
RUN mkdir /jacaranda/logs

ADD *.yml /jacaranda/cfg/

ENV CONF_PATH /jacaranda/cfg
ENV ENV pro

EXPOSE 8001
