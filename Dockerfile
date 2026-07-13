# syntax=docker/dockerfile:1
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
      -ldflags="-X github.com/y-yu/kindle-clock-go/domain/build.gitCommitHash=$GIT_COMMIT_HASH -s -w" \
      -trimpath \
      -o /app/server main.go

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* &&  \
    update-ca-certificates

ENV DOSIS_FONT_PATH=/etc/Dosis.ttf
ENV ROBOTO_SLAB_FONT_PATH=/etc/RobotoSlab.ttf
ARG GOOGLE_FONTS_COMMIT=ec0464b978de222073645d6d3366f3fdf03376d8
RUN curl --fail --location --retry 3 --output "$ROBOTO_SLAB_FONT_PATH" "https://raw.githubusercontent.com/google/fonts/${GOOGLE_FONTS_COMMIT}/apache/robotoslab/RobotoSlab%5Bwght%5D.ttf" && \
    curl --fail --location --retry 3 --output "$DOSIS_FONT_PATH" "https://raw.githubusercontent.com/google/fonts/${GOOGLE_FONTS_COMMIT}/ofl/dosis/Dosis%5Bwght%5D.ttf" && \
    test -s "$ROBOTO_SLAB_FONT_PATH" && test -s "$DOSIS_FONT_PATH"

COPY --from=builder /app/server /bin/server
COPY --from=builder /app/etc/weather_icon /etc/weather_icon

EXPOSE 8080
ENTRYPOINT ["/bin/server"]
