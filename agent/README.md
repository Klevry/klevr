# klevr-agent
## How to use
* Help
```
#] ./klevr -h
Usage of /tmp/go-build009650043/b001/exe/main:
  -ip string
    	local IP address for networking (default "192.168.2.100")
  -provider string
    	[baremetal|aws] - Service Provider for Host build up
  -user string
    	Account key from Klevr service
  -webconsole string
    	Klevr webconsole(server) address (URL or IP, Optional: Port) for connect (default "localhost:8080")
```

 * Using localhost: `./klevr -provider=baremetal -user=ralf`
 
 * Using seperated host:  `./klevr -provider=baremetal -user=ralf -webconsole=[WEBCONSOL_ADDR] `
