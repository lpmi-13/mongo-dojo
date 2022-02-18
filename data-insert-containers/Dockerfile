FROM golang:1.15.2-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /insert-data

# I wanted to use a scratch container here, but keeping something bigger until it's working
FROM golang:1.15.2-alpine

COPY --from=builder /insert-data /insert-data

# and now we can connect directly to the primary via the DNS
CMD ["/insert-data", "50000", "mongodb://mongo1:27017"]
