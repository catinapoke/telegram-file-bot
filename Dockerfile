FROM golang

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build -o /filebot -v ./...

EXPOSE 2000

CMD [ "/filebot"]