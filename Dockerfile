FROM golang as src

COPY src /usr/src
WORKDIR /usr/src

RUN go build -o /usr/local/bin/uuid .

FROM debian:11

COPY --from=src /usr/local/bin/uuid /usr/local/bin/uuid
