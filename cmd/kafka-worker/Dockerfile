#снова выполняем сборку в два этапа, чтобы в финальном образе не было лишних зависимостей и исходников, а был только бинарник
FROM golang:1.25.0-bookworm AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/kafka-worker/exe /app/cmd/kafka-worker

#2 этап - финальный образ, в котором будет только бинарник и ничего лишнего
FROM alpine:3.24
WORKDIR /app

COPY --from=builder /app/cmd/kafka-worker/exe /app
CMD ["/app/exe"]
