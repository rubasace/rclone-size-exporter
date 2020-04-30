FROM golang:1.14-alpine as build

ENV GO111MODULE=on


WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY main.go .

RUN go build

FROM rclone/rclone

WORKDIR /app

COPY --from=build /app/rclone-size-exporter /app/rclone-size-exporter

ENV PORT 8080
ENV RCLONE_CONFIG /config/rclone/rclone.conf

ENTRYPOINT ["/app/rclone-size-exporter"]