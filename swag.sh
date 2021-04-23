#!/bin/bash	

GOPATH=`go env | grep GOPATH | sed -n 's/^GOPATH=//p' | sed -n 's/"//gp'`
cd pkg/manager
# go mod vendor
echo ${GOPATH}
${GOPATH}/bin/swag init -g server.go --parseDependency --parseInternal true