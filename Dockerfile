FROM golang:alpine

RUN go version
ENV GOPATH=/

WORKDIR /app
COPY ./ ./

RUN go build ./cmd/main.go

CMD ["./main"]