FROM alpine:3.13.2
LABEL version=0.3

RUN apk update && apk add curl bash vim

COPY ./conf/* /conf/
COPY ./klevr-manager /
COPY ./wait-for-it.sh /
COPY ./init-db.sh /
CMD [ "./klevr-manager" ]
EXPOSE 8090

