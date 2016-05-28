#
# Dockerfile for Jacaranda alerting project
#

FROM golang:onbuild

MAINTAINER Joan Oliva

RUN mkdir /jacaranda
RUN mkdir /jacaranda/bin
RUN mkdir /jacaranda/cfg
RUN mkdir /jacaranda/logs

ADD *.yml /jacaranda/cfg/

ENV CONF_PATH /jacaranda/cfg
ENV ENV pre

EXPOSE 8000