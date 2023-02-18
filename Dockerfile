FROM golang:1.20-alpine as build
WORKDIR /usr/src
COPY . .
ENV CGO_ENABLED 0
RUN go version && go build -o goapp -v -ldflags="-s -w" .

# hadolint ignore=DL3006
FROM gcr.io/distroless/static-debian11
LABEL maintainer="bat@sbz.fr"
ARG uid
USER ${uid:-65534}
COPY --from=build /usr/src/goapp /goapp
EXPOSE 8080
ENTRYPOINT ["/goapp"]