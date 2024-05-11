FROM golang:1.22 AS build

WORKDIR /build

ARG DIST_TARGET

COPY . .

RUN go install github.com/goreleaser/goreleaser@latest
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$DIST_TARGET goreleaser build --rm-dist --snapshot --verbose --single-target --output .

FROM debian:12

WORKDIR /app

COPY --from=build ./aws-authenticator .

RUN chmod +x aws-authenticator

CMD ["/aws-authenticator"]