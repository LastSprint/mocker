FROM golang:1.13 	

RUN git clone https://github.com/LastSprint/mocker.git /app

WORKDIR /app/

RUN git checkout v3.0.0

RUN go build

COPY config.json config.json
