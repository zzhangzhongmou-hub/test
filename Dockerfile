FROM golang:1.25-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/configs/config.yaml ./Configs/config.yaml
COPY --from=builder /app/main .

COPY --from=builder /app/static ./static

ENV TZ=Asia/Shanghai
EXPOSE 8080

CMD ["./main"]