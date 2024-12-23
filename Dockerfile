FROM golang:1.23

WORKDIR /app

COPY src/ragserver ./

EXPOSE 8080

CMD ["./ragserver"]
