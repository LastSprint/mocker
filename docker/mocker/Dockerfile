FROM golang:latest 	

RUN git clone https://github.com/LastSprint/mocker.git /app
WORKDIR /app/src/mocker

ENV GOPATH=/app/

RUN go get github.com/sirupsen/logrus

RUN go build .

COPY config.json config.json