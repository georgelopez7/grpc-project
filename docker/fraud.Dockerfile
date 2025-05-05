# BUILD STAGE
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o fraud ./cmd/fraud

# RUN STAGE
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/fraud .

EXPOSE 50052
ENTRYPOINT ["./fraud"]