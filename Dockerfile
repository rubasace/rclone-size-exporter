FROM golang:1.14-alpine as build

ENV GO111MODULE=on


WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY main.go .

RUN go build

FROM rclone/rclone

ENV PORT=8080

WORKDIR /app

COPY --from=build /app/rclone-size-exporter /app/rclone-size-exporter

ENTRYPOINT ["/app/rclone-size-exporter"]