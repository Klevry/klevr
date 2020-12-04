#!/usr/bin/env bash

USERNAME=$1
PASSWORD=$2
HOST=$3
PORT=$4
SCHEME=$5
SQL_PATH=$6
EXPORT_PATH=$7

EXPORT_FILE=${EXPORT_PATH}/${SCHEME}_$(date +%Y%m%d%H%M).sql

CMD_PRE="mysql --user=${USERNAME} --password=${PASSWORD} --host=${HOST} --port=${PORT}"


apt-get update
apt-get -y install mariadb-client


EXISTS=$(${CMD_PRE} --execute "show databases" | grep "${SCHEME}")


if [[ ${EXISTS} != "" && ${EXPORT_PATH} != "" ]] ; then
	echo "=============== start export for backup scheme ==============="
	CMD=$(mysqldump --user=${USERNAME} --password=${PASSWORD} --host=${HOST} --port=${PORT} -e --single-transaction -c ${SCHEME} > ${EXPORT_FILE})
	echo "=============== complete export for backup scheme ==============="
fi


if [[ ${EXISTS} != "" ]] ; then
	SQL_PATH="${SQL_PATH}.modify"
else
	SQL_PATH="${SQL_PATH}.create"
fi


if [ -s "${EXPORT_FILE}" ] ; then
    echo "=============== start import scheme ==============="
    CMD=$(${CMD_PRE} --execute "source ${SQL_PATH}")
    echo "=============== complete import scheme ==============="
elif [[ ${EXISTS} == "" || ${EXPORT_PATH} == "" ]] ; then
	echo "=============== start import scheme ==============="
    CMD=$(${CMD_PRE} --execute "source ${SQL_PATH}")
    echo "=============== complete import scheme ==============="
else
	echo "Import failed due to schema export failed."
	exit 1
fi