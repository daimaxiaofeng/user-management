FROM golang:latest

ENV GOPROXY=https://goproxy.cn

WORKDIR /app
COPY . .

RUN go build

EXPOSE 2024

CMD ["./user-management"]
