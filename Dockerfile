FROM golang:1.15.5-alpine
COPY . /go/app-demo
WORKDIR /go/app-demo
RUN go build -o adp-demo
EXPOSE 18080

CMD ["./adp-demo"]