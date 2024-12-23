FROM golang:1.23

WORKDIR /app

RUN ls

COPY src/ragserver ./

EXPOSE 8080

CMD ["./ragserver"]
