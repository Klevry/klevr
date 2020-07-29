FROM ubuntu:focal
LABEL version=0.1

COPY ./conf/* /conf/
COPY ./klevr-manager /
COPY ./wait-for-it.sh /
CMD [ "./klevr-manager" ]
EXPOSE 8090

