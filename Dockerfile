FROM golang:1.23-alpine

WORKDIR /opt/wrauth

COPY *.go go.mod ./
RUN go get github.com/Skaytacium/wrauth
RUN mkdir /config
RUN CGO_ENABLED=0 GOOS=linux go build -o /wrauth

EXPOSE 9092
VOLUME /config

CMD ["/wrauth", "-config", "/config/config.yaml", "-db", "/config/db.yaml"]
