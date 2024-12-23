FROM golang:1.23

WORKDIR /app

COPY src/ ./

EXPOSE 8080

CMD ["./ragserver"]
