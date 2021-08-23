#!/bin/bash

echo " Now, Downloading the klevr agent, Please wait..."
curl -sL https://github.com/Klevry/klevr/raw/master/Dockerfiles/agent/klevr -o ~/klevr
chmod 755 ~/klevr

## klevr -apiKey=${K_API_KEY} -platform=${K_PLATFORM} -manager=${K_MANAGER_URL} -zoneId=${K_ZONE_ID}
#nohup ./klevr -apiKey=121234123 -platform=linux_laptop -manager="http://localhost:8090" -zoneId=123124123 >> /tmp/klevr_agent.log 2>&1 &
#sh -c "nohup ./klevr -apiKey=$1 -platform=$2 -manager=$3 -zoneId=$4 >> /tmp/klevr_agent.log 2>&1 &"
