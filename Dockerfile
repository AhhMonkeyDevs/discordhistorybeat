FROM golang:1.15.8

WORKDIR /go/src/github.com/AhhMonkeyDevs/discordhistorybeat

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /usr/share/discordbeat
COPY --from=0 /go/src/github.com/AhhMonkeyDevs/discordhistorybeat .
STOPSIGNAL SIGINT
ENTRYPOINT ["./discordhistorybeat"]

