FROM golang:1.26.3-alpine3.23 as builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o /main

FROM alpine:latest AS final

COPY --from=builder /main .

EXPOSE 5000

CMD ["/main"]