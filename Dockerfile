FROM golang:latest as builder
WORKDIR /go/src/rcl-assistant/assistant
RUN go get gopkg.in/yaml.v2
COPY *.go .
COPY messages/* ./messages/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rcl-assistant .

FROM alpine:latest

EXPOSE 8080

COPY --from=builder /go/src/rcl-assistant/assistant /usr/local/bin
CMD ["/usr/local/bin/rcl-assistant"]