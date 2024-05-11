FROM golang:1.22 AS build

WORKDIR /build

ARG DIST_TARGET

COPY . .

RUN go install github.com/goreleaser/goreleaser@latest
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$DIST_TARGET goreleaser build --rm-dist --snapshot --verbose --single-target --output .

FROM alpine:3.14

WORKDIR /app


COPY --from=build ./aws-authenticator .

ARG DYNAMODB_REGION=eu-central-2
ARG DYNAMODB_TABLE_NAME=authenticator
ARG APP_ROUTER=LOGIN

ENV AWS_DYNAMODB_REGION=$DYNAMODB_REGION
ENV AWS_DYNAMODB_TABLE_NAME=$DYNAMODB_TABLE_NAME
ENV AWS_ROUTER=$APP_ROUTER

RUN chmod +x aws-authenticator

CMD ["/aws-authenticator"]