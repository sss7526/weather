FROM golang:1.24.4-alpine3.22 AS builder

WORKDIR /app

COPY . .

RUN go build -ldflags="-s -w" -trimpath -o weatherdotexe .

FROM alpine:latest AS runtime

RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/weatherdotexe .

CMD [ "./weatherdotexe" ]