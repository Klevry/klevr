FROM ubuntu:focal
<<<<<<< HEAD
LABEL version=0.3

RUN apt-get update -y && apt-get install -y ssh
#RUN mkdir -p root/.ssh/ && \
#    cd root/.ssh && \
#    echo $'Host * \n\
#        StrictHostKeyChecking no \n\
#        UserKnownHostsFile /dev/null \n '\
#    > config

=======
LABEL version=0.1.4
>>>>>>> 3285330a9d39ab7dcabd3163d9231cdc8b4a76e2
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
