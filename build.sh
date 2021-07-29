#!/bin/bash	

cd cmd/klevr-agent	
make build	
cd ../..	

cd cmd/klevr-manager	
make build	
cd ../..	

cd console
make pre 
cd ..

echo "Binary build complete"
