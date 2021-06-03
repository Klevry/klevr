#!/usr/bin/env bash

USERNAME=$1
PASSWORD=$2
HOST=$3
PORT=$4
SCHEME=$5
SQL_PATH=$6
EXPORT_PATH=$7
APP_USER=$8
APP_PASSWD=$9

echo USERNAME=${USERNAME}
echo PASSWORD=${PASSWORD}
echo HOST=${HOST}
echo PORT=${PORT}
echo SCHEME=${SCHEME}
echo SQL_PATH=${SQL_PATH}
echo EXPORT_PATH=${EXPORT_PATH}
echo APP_USER=${APP_USER}
echo APP_PASSWD=${APP_PASSWD}

EXPORT_FILE=${EXPORT_PATH}/${SCHEME}_$(date +%Y%m%d%H%M).sql

CMD_PRE="mysql --user=${USERNAME} --password=${PASSWORD} --host=${HOST} --port=${PORT}"


apk update
apk add mariadb-client


EXISTS=$(${CMD_PRE} --execute "show databases" | grep "${SCHEME}")


if [[ ${EXISTS} != "" && ${EXPORT_PATH} != "" ]] ; then
	echo "=============== start export for backup scheme ==============="
	mkdir ${EXPORT_PATH}
	CMD=$(mysqldump --user=${USERNAME} --password=${PASSWORD} --host=${HOST} --port=${PORT} -e --single-transaction -c ${SCHEME} > ${EXPORT_FILE})
	echo "=============== complete export for backup scheme ==============="
fi


if [[ ${EXISTS} != "" ]] ; then
	cat > ${SQL_PATH}.execute <<- EOM
		USE \`${SCHEME}\`;
	EOM
	
	cat ${SQL_PATH}.modify >> ${SQL_PATH}.execute
	
	SQL_PATH="${SQL_PATH}.execute"
else
	cat > ${SQL_PATH}.execute <<- EOM
		CREATE DATABASE IF NOT EXISTS \`${SCHEME}\` DEFAULT CHARACTER SET utf8;
		USE \`${SCHEME}\`;
		CREATE USER IF NOT EXISTS \`${APP_USER}\`@\`%\` IDENTIFIED BY '${APP_PASSWD}';
		GRANT ALL PRIVILEGES ON \`${SCHEME}\`.* to \`${APP_USER}\`@\`%\`;
	EOM
	
	cat ${SQL_PATH}.create >> ${SQL_PATH}.execute
	
	SQL_PATH="${SQL_PATH}.execute"
fi

cat ${SQL_PATH} > ${EXPORT_PATH}/${SCHEME}.execute.sql
echo SQL_PATH=${SQL_PATH}
echo ${EXPORT_PATH}/execute.sql


if [ -s "${EXPORT_FILE}" ] ; then
    echo "=============== start import scheme ==============="
    CMD=$(${CMD_PRE} -f --execute "source ${SQL_PATH}")
    echo "=============== complete import scheme ==============="
elif [[ ${EXISTS} == "" || ${EXPORT_PATH} == "" ]] ; then
	echo "=============== start import scheme ==============="
    CMD=$(${CMD_PRE} -f --execute "source ${SQL_PATH}")
    echo "=============== complete import scheme ==============="
else
	echo "Import failed due to schema export failed."
	exit 1
fi

