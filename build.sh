#!/bin/bash	

cd cmd/klevr-agent	
make build	
cd ../..	

cd cmd/klevr-manager	
make build	
cd ../..	

cd console
make build
cd ..

echo "Binary build complete"
