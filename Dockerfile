FROM golang

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
# RUN go build  ./...
RUN go build -v -o /usr/local/bin/bot ./main.go 

EXPOSE 2000

CMD [ "bot"]