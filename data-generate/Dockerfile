FROM golang:1.15.2-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /generate-data

FROM scratch

COPY --from=builder /generate-data /generate-data

# the /data-generate/ directory is used as a bind mount in the docker-compose file
# to keep the files generated after container exit
CMD ["/generate-data", "5000000", "data-generate"]
