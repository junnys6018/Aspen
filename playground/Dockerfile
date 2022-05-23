FROM golang:1.17

WORKDIR /usr/src/playground

ENV PLAYGROUND_PORT="8080"
ENV PLAYGROUND_TIMEOUT_DURATION="4000ms"
ENV PLAYGROUND_ASPEN_PATH="/usr/src/playground/aspen"

COPY aspen ./
COPY main.go ./

RUN go build -o playground -v main.go

EXPOSE ${PLAYGROUND_PORT}

CMD ["./playground"]