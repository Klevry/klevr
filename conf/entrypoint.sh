#!/bin/bash

db_host="klevr-db"


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
		App_run
        fi
}

Check_db


