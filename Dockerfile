FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -mod=readonly -o /server cmd/api/main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /server /server

EXPOSE 8080

# Run
CMD ["/server"]
