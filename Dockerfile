FROM golang:1.13.2 AS builder
WORKDIR /go/src/github.com/jjshanks/header-proxy/
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch
COPY --from=builder /go/src/github.com/jjshanks/header-proxy/app .
ENTRYPOINT ["./app"]