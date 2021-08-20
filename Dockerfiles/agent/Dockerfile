FROM alpine:3.13.2
LABEL version=0.1.4

RUN apk add --no-cache --update openssh-client jq

### klevr go Build Opt : env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o
COPY ./klevr  /
COPY ./entrypoint.sh  /

# Set environment variable
ENV K_API_KEY ""
ENV K_PLATFORM ""
ENV K_MANAGER_URL ""
ENV K_ZONE_ID ""

# Entrypoint for startup
ENTRYPOINT /entrypoint.sh

EXPOSE 18800
