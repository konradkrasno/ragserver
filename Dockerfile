FROM golang:1.23

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ ./
COPY config.yaml ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./ragserver

EXPOSE 8080

CMD ["./ragserver"]
