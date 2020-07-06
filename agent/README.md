# klevr-agent
## Kickstart
```
curl -sL https://bit.ly/klevry |bash
```
## Primary host election
* Image click to Youtube
   * [![Image click to Youtube](https://github.com/ralfyang/klevr/blob/master/src/primary_election_s.png)](https://youtu.be/6-fV-ubTwXw)

## How to use
* Help
```
#] ./klevr -h
Usage of ./klevr:
  -id string
    	Account ID from Klevr service
  -ip string
    	local IP address for networking (default "192.168.1.21")
  -provider string
    	[baremetal|aws] - Service Provider for Host build up
  -webconsole string
    	Klevr webconsole(server) address (URL or IP, Optional: Port) for connect (default "localhost:8080")
```

 * Using localhost: `./klevr -provider=baremetal -id=ralf`
 
 * Using seperated host:  `./klevr -provider=baremetal -id=ralf -webconsole=[WEBCONSOL_ADDR] `
