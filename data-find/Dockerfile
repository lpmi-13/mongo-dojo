FROM golang:1.15.2-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /find

FROM scratch

COPY --from=builder /find /find

CMD ["/find", "mongodb://localhost:27017/?replicaSet=dojo", "5"]
