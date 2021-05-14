FROM golang:1.16

WORKDIR /usr/src/app

COPY . .

RUN go get .
RUN go build -o app .

CMD [ "./app" ]
