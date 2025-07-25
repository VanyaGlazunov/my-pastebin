# Этап 1: Сборка
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем зависимости и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
# CGO_ENABLED=0 нужен для статической линковки
# -ldflags="-w -s" уменьшает размер бинарного файла
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/pastebin-server ./cmd/main.go

# Этап 2: Запуск
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированный бинарник из этапа сборки
COPY --from=builder /app/pastebin-server .
# Копируем .env файл, если он нужен внутри контейнера (хотя лучше передавать через docker-compose)
# COPY .env . 

EXPOSE 8080

# Запускаем приложение
CMD ["./pastebin-server"]