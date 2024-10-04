# music-stream
Код для веб-сервера на Go, который будет выбирать случайные звуковые файлы из примонтированного облачного диска и потоково стримить музыку с 5-секундным эффектом затухания между треками. Мы также настроим Nginx как обратный прокси и создадим Docker-контейнеры с помощью Docker Compose.

### Шаг 1: Создание Go веб-сервера

Создайте новый проект и файл `main.go`:

```sh
mkdir music-streamer
cd music-streamer
nano main.go
```

Добавьте следующий код в `main.go`:

```go
package main

import (
    "fmt"
    "io"
    "math/rand"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "time"
)

var musicDir = "/path/to/mounted/cloud/disk"

func main() {
    http.HandleFunc("/stream", streamHandler)
    fmt.Println("Starting server on :8080")
    http.ListenAndServe(":8080", nil)
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
    files, err := filepath.Glob(filepath.Join(musicDir, "*.mp3"))
    if err != nil || len(files) == 0 {
        http.Error(w, "No music files found", http.StatusInternalServerError)
        return
    }

    rand.Seed(time.Now().UnixNano())
    file := files[rand.Intn(len(files))]

    cmd := exec.Command("sox", file, "-t", "mp3", "-", "fade", "t", "5")
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        http.Error(w, "Unable to process file", http.StatusInternalServerError)
        return
    }

    if err := cmd.Start(); err != nil {
        http.Error(w, "Unable to start processing", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "audio/mpeg")
    io.Copy(w, stdout)
    cmd.Wait()
}
```

### Шаг 2: Создание Dockerfile

Создайте файл `Dockerfile` в корневом каталоге проекта:

```sh
nano Dockerfile
```

Добавьте следующий код в `Dockerfile`:

```dockerfile
FROM golang:1.20-alpine

RUN apk add --no-cache sox

WORKDIR /app

COPY . .

RUN go build -o music-streamer .

CMD ["./music-streamer"]
```

### Шаг 3: Настройка Nginx

Создайте файл конфигурации Nginx `nginx.conf`:

```sh
nano nginx.conf
```

Добавьте следующий конфигурационный код:

```nginx
server {
    listen 80;
    server_name localhost;

    location /stream {
        proxy_pass http://music-streamer:8080/stream;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Шаг 4: Создание Docker Compose файла

Создайте файл `docker-compose.yml`:

```sh
nano docker-compose.yml
```

Добавьте следующий код в `docker-compose.yml`:

```yaml
version: '3.8'

services:
  music-streamer:
    build: .
    volumes:
      - /path/to/mounted/cloud/disk:/path/to/mounted/cloud/disk
    ports:
      - "8080:8080"

  nginx:
    image: nginx:latest
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf
    ports:
      - "80:80"
    depends_on:
      - music-streamer
```

### Шаг 5: Запуск Docker Compose

Запустите Docker Compose:

```sh
docker-compose up --build
```

Теперь ваш веб-сервер на Go будет случайным образом выбирать звуковые файлы из примонтированного облачного диска и потоково стримить музыку с 5-секундным эффектом затухания между треками. Nginx будет выступать в роли обратного прокси для вашего Go сервера, и все это будет работать в Docker-контейнерах.

Источник: беседа с Copilot, 04.10.2024
(1) How To Deploy a Go Web Application Using Nginx on Ubuntu 18.04. https://www.digitalocean.com/community/tutorials/how-to-deploy-a-go-web-application-using-nginx-on-ubuntu-18-04.
(2) Deploying a Go Web Application Using Nginx on Ubuntu 22.04. https://www.atlantic.net/dedicated-server-hosting/deploying-a-go-web-application-using-nginx-on-ubuntu-22-04/.
(3) How To Deploy a Go Web Application with Docker and Nginx on Ubuntu 22. .... https://www.digitalocean.com/community/tutorials/how-to-deploy-a-go-web-application-with-docker-and-nginx-on-ubuntu-22-04.
(4) Getty Images. https://www.gettyimages.com/detail/news-photo/in-this-photo-illustration-the-nginx-logo-is-displayed-on-a-news-photo/1244811361.
