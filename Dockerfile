FROM --platform=linux/amd64 golang:1.21.0

EXPOSE 8080

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

CMD ["/payload-relay"]
