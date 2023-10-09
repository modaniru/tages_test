FROM golang:latest

RUN go version
ENV GOPATH=/

WORKDIR /app
COPY ./ ./

RUN go build -o tages-test ./cmd/main.go

CMD ["./tages-test"]