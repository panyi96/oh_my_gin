FROM golang:1.19.3-alpine3.17
COPY . /src/
WORKDIR /src
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"
RUN go build -o ohmygin .
EXPOSE 8888
ENTRYPOINT  ["./ohmygin"]

