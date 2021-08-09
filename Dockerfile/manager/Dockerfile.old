FROM alpine:3.13.2
LABEL version=0.4

RUN apk update && apk add curl bash vim libc6-compat

COPY ./conf/* /conf/
COPY ./klevr-manager /
COPY ./entrypoint.sh /
ENTRYPOINT ./entrypoint.sh
EXPOSE 8090
