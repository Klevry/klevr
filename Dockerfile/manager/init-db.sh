#!/usr/bin/env bash

USERNAME=$1
PASSWORD=$2
HOST=$3
PORT=$4
SCHEME=$5
SQL_PATH=$6

apt-get update
apt-get -y install mysql-client

CMD_PRE="mysql --user=${USERNAME} --password=${PASSWORD} --host=${HOST} --port=${PORT}"

CMD=$(${CMD_PRE} --execute "show databases" | grep "${SCHEME}")

if [[ ${CMD} == "" ]] ; then
    echo "start import scheme"
    CMD=$(${CMD_PRE} --execute "source ${SQL_PATH}")
    echo "complete import scheme"
fi