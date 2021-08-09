#!/bin/bash

db_host="klevr-db"


Init_db(){
	/init-db.sh $${DB_MGMT_USER_NAME} $${DB_MGMT_USER_PASSWORD} klevr-db 3306 klevr /conf/klevr-manager-db.sql . $${DB_APP_USER} $${DB_APP_PASSWORD}
}

App_run(){
	#/klevr-manager -c /conf/klevr-manager.yml
	/klevr-manager -c /conf/klevr-manager-compose.yml
}

Check_db(){
        set -eu
        echo "Checking DB connection ..."

        i=0
        while [ $i -lt 10 ];do
                nc -z $db_host 3306 && break
                echo "$i: Waiting for DB 5 second ..."
                let i=$i+1
                sleep 5
        done

        if [ $i -eq 10 ];then
                echo "DB connection refused, terminating ..."
                exit 1
        else
                echo "DB is up ..."
        fi
}

Init_db
Check_db
App_run
