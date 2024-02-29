FROM golang:1.22.0

LABEL maintainer="Adam Barcan <abarcan@gmail.com>"

RUN \
    DEBIAN_FRONTEND=noninteractive apt-get update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get upgrade -y \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential

ENV APP_HOME /usr/src/app/

WORKDIR $APP_HOME

COPY . $APP_HOME

EXPOSE 8080

CMD ["go", "run", "main.go"]
