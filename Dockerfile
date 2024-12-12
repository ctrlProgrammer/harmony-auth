FROM golang:1.21.6

WORKDIR /
COPY . .
ENTRYPOINT ["go", "run", "main.go"]