# Используем официальный образ Go для сборки приложения
FROM golang:1.18-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем все файлы проекта в рабочую директорию
COPY . .

# Инициализируем модуль Go и загружаем зависимости
RUN go mod init music-streamer
RUN go mod tidy

# Собираем бинарный файл приложения
RUN go build -o music-streamer .

# Используем минимальный образ для запуска приложения
FROM alpine:latest

# Устанавливаем необходимые зависимости
RUN apk --no-cache add sox

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /root/

# Копируем собранный бинарный файл из предыдущего этапа
COPY --from=builder /app/music-streamer .

# Устанавливаем путь к директории с музыкой
ENV MUSIC_DIR=/path/to/mounted/cloud/disk

# Открываем порт 8080 для доступа к приложению
EXPOSE 8080

# Запускаем приложение
CMD ["./music-streamer"]
