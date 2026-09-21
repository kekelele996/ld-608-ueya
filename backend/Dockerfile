FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o server ./main.go
FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/server ./server
CMD ["./server"]
