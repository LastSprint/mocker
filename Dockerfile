FROM golang:1.15

COPY src /app/

WORKDIR /app/

RUN go build -v

CMD ./mocker