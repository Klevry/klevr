FROM ubuntu:focal
LABEL version=0.1.4

RUN apt-get update -y && \
    apt-get install -y ssh && \
    apt-get -y clean && \
    rm -rf /var/lib/apt/lists/*

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
