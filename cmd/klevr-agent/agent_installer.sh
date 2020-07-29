#!/bin/bash

echo " Now, Downloading the klevr agent, Please wait..."
curl -L https://github.com/ralfyang/klevr/blob/master/pkg/agent/klevr?raw=true -o ~/klevr
chmod 755 ~/klevr

