FROM ubuntu:focal
LABEL version=0.1

ADD ./conf/* /conf/
ADD ./klevr-manager /
ADD ./wait-for-it.sh /
CMD [ "./klevr-manager" ]
EXPOSE 8090

