from golang:1.15.2-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./

RUN go build -o /insert-data

CMD ["/insert-data", "500000"]
