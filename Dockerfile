FROM golang:1.22.0

LABEL maintainer="Adam Barcan <abarcan@gmail.com>"

RUN \
    DEBIAN_FRONTEND=noninteractive apt-get update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get upgrade -y \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential

WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/
COPY vendor/ /app/vendor/
COPY cmd/ /app/cmd/
COPY internal/ /app/internal/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o $GOPATH/bin ./...

EXPOSE 8080

CMD ["http-server"]
