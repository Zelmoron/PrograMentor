FROM golang:1.23.4-alpine

WORKDIR /app

# Установка нужных пакетов и docker CLI
RUN apk add --no-cache bash curl tar
RUN curl -fsSL https://download.docker.com/linux/static/stable/x86_64/docker-24.0.6.tgz | \
    tar -xzvf - --strip-components=1 -C /usr/local/bin docker/docker

COPY . .

RUN go mod tidy
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]