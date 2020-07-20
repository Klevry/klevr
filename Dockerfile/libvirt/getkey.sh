#!/bin/bash
key1="http://pubkey.nexclipper.io:9876/pubkey"
key2="http://zxz.kr/nexc_pubkey"
#if [ ! -e ~/.ssh/config ]; then
#	mkdir ~/.ssh/
#	cp -Rfvp /tmp/ssh-config ~/.ssh/config
#fi
pubkeyget_IP(){
        pubkey=$(curl -ksfL $key1 --connect-timeout 2)
        if [[ $pubkey != "" ]]; then
                echo "${pubkey}"
        else
                pubkey=$(curl -ksfL $key2 --connect-timeout 2)
                echo "${pubkey}"
        fi
}
pubkeyget_IP
