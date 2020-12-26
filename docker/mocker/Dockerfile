FROM golang:1.13 	

RUN git clone https://github.com/LastSprint/mocker.git /app

WORKDIR /app/

RUN git checkout v4.1.1

RUN go build

CMD "./mocker"