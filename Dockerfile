FROM golang:1.25 AS builder

WORKDIR /app

# Сначала зависимости для кеша
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/tracker ./cmd/tracker

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /app/bin/tracker /app/tracker
EXPOSE 8080

ENV HTTP_PORT=8080
CMD ["/app/tracker"]
