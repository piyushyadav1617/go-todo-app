FROM golang:1.23.4-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN  go mod download
COPY . .

RUN go build -v -o ./go-todo-app .

EXPOSE 8080
CMD ["./go-todo-app"]
