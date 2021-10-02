FROM golang:1.15.2-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 go build -o /insert-data

FROM scratch

COPY --from=builder /insert-data /insert-data

CMD ["/insert-data", "500000"]
