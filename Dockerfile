FROM golang:1.20-alpine

RUN apk add --no-cache sox

WORKDIR /app

COPY . .

RUN go build -o music-streamer .

CMD ["./music-streamer"]