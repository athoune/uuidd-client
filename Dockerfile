FROM golang as src

COPY ./uuidd/ /usr/src/uuidd/
COPY main.go go.* /usr/src/
WORKDIR /usr/src

RUN ls -l && go build -o /usr/local/bin/uuid .

FROM debian:11

COPY --from=src /usr/local/bin/uuid /usr/local/bin/uuid
