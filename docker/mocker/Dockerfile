FROM golang:latest 	

RUN git clone https://github.com/LastSprint/mocker.git /app

WORKDIR /app/

RUN go build

COPY config.json config.json
