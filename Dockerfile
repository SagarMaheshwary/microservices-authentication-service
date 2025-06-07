FROM golang:1.23 AS builder

WORKDIR /app

ARG MODE=production

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# create execulable in case of production mode
RUN if [ "$MODE" = "production" ]; then CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o . ./cmd/server/main.go; fi

FROM alpine:3.21 AS production

WORKDIR /app

COPY --from=builder /app/main /app/main

CMD [ "./main" ]

FROM builder AS development

WORKDIR /app

COPY --from=builder /app /app

RUN go install github.com/air-verse/air@v1.52.3

CMD ["air", "-c", ".air.toml"]
