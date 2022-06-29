FROM golang:1.15.2-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /data-insert-background

FROM scratch

COPY --from=builder /data-insert-background /data-insert-background

CMD ["/data-insert-background", "mongodb://localhost:27017/?replicaSet=dojo", "5"]
