#!/bin/bash

## Replace the API_SERVER by environment variable

#api_target=$KLEVR_MANAGER_IP

## Test stemping
echo "$api_target" > /app/api_server_addr.txt

sed -i  "s#%%KLEVR_API_SERVER_IP_MARKUP%%#$KLEVR_MANAGER_IP#g" /app/build/static/js/*.js

nginx -g 'daemon off;'
