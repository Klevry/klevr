#!/bin/bash

cd /go/src/app

cd cmd/klevr-agent
make build
cd ../..

cd cmd/klevr-manager
make build
cd ../..

echo "Binary build complete"