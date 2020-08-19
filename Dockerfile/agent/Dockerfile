FROM ubuntu:focal
LABEL version=0.3

RUN apt-get update -y && apt-get install -y ssh
#RUN mkdir -p root/.ssh/ && \
#    cd root/.ssh && \
#    echo $'Host * \n\
#        StrictHostKeyChecking no \n\
#        UserKnownHostsFile /dev/null \n '\
#    > config

COPY ./klevr  /
COPY ./entrypoint.sh  /

# Set environment variable
ENV K_API_KEY ""
ENV K_PLATFORM ""
ENV K_MANAGER_URL ""
ENV K_ZONE_ID ""

# Entrypoint for startup
ENTRYPOINT /entrypoint.sh

EXPOSE  8080

