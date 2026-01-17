FROM golang:1.23-alpine AS build

WORKDIR /opt/wrauth

COPY *.go go.mod ./
RUN go get github.com/Skaytacium/wrauth
RUN CGO_ENABLED=0 GOOS=linux go build -o /wrauth

FROM alpine

COPY --from=build /wrauth /

RUN mkdir /config
VOLUME /config

EXPOSE 9092

CMD ["/wrauth", "-config", "/config/config.yaml", "-db", "/config/db.yaml"]
