FROM golang:1.18-alpine as compiler

RUN apk update
RUN apk upgrade -U
RUN apk add gcc gcompat libgcc libwebp  musl-dev

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=1 go build -o /discord-bot

FROM alpine
WORKDIR /

RUN apk update
RUN apk upgrade -U
RUN apk add gcc gcompat libgcc libwebp  musl-dev
RUN apk add git
COPY --from=compiler /discord-bot /

ENTRYPOINT ["/discord-bot"]