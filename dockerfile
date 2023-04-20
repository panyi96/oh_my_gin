FROM golang:1.19.3-alpine3.17 AS builder
COPY . /src/
WORKDIR /src
RUN GOPROXY="https://proxy.golang.com.cn,direct" GO111MODULE=on go build

FROM alpine:latest
COPY --from=builder /src/ohmygin /bin/ohmygin
WORKDIR /bin
EXPOSE 1234
ENTRYPOINT  ["/bin/ohmygin"]
CMD []

