FROM docker.io/library/golang:1.21.3-alpine3.18 as build


COPY . /opt/app

WORKDIR /opt/app

RUN go mod download && go mod verify

RUN go build -v -o /opt/app/telegram-bot cmd/main.go


FROM registry.access.redhat.com/ubi9/ubi:9.2

COPY --from=build /opt/app/telegram-bot /opt/telegram-bot

RUN chmod +x /opt/telegram-bot

ENTRYPOINT ["/opt/telegram-bot"]
