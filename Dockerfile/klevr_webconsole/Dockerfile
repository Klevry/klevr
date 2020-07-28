FROM ubuntu:focal
LABEL version=0.3
COPY ./klevr_webconsole  /
COPY ./entrypoint.sh  /
ENV API_SERVER "0.0.0.0"
ENTRYPOINT /entrypoint.sh
EXPOSE  8080
