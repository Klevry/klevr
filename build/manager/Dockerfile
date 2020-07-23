FROM ubuntu:focal
COPY ./cmd/klevr-manager/klevr-manager /
COPY ./conf/* /conf/
COPY ./wait-for-it.sh /
CMD [ "./klevr-manager" ]
EXPOSE 8090