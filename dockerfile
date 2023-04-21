FROM golang:1.19.3-alpine3.17 AS builder
COPY . /src/
WORKDIR /src
RUN GOPROXY="https://proxy.golang.com.cn,direct" GO111MODULE=on go build

FROM alpine:latest
COPY --from=builder /src/ohmygin /bin/ohmygin
ENV NACOS_IP="127.0.0.1"
ENV PORT=8848
ENV GRPC_PORT=9848
WORKDIR /bin
EXPOSE 1234
ENTRYPOINT  ["/bin/ohmygin"]
CMD ["$NACOS_IP","$PORT","$GRPC_PORT"]

