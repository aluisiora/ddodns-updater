FROM golang:1.22.5-alpine3.20 AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o ./bin/ddodns-updater-linux-x64 *.go

FROM golang:1.22.5-alpine3.20 

ARG UID=998
ARG GID=998

RUN addgroup -S app -g ${GID} \
    && adduser -S -G app -u ${UID} app

WORKDIR /usr/src/app

COPY --from=build --chown=app:app /usr/src/app/bin/ddodns-updater-linux-x64 /usr/local/bin/ddodns-updater-linux-x64

USER app

CMD ["ddodns-updater-linux-x64"]
