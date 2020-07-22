# klevr-agent
## Kickstart
```
curl -sL bit.ly/klevry |bash  && ./klevr -apiKey=\"{apiKey}\" -platform={platform} -manager=\"{managerUrl}\" -zoneId={zoneId}
```
## Primary host election
* Image click to Youtube
   * [![Image click to Youtube](https://github.com/Klevry/klevr/blob/master/assets/primary_election_s.png)](https://youtu.be/hyMaVsCcgbA)

## How to use
* Help
```
#] ./klevr -h
Usage of ./klevr:
  -group string
    	Group name will be a uniq company name or team name
  -id string
    	Account ID from Klevr service
  -ip string
    	local IP address for networking (default "192.168.15.50")
  -platform string
    	[baremetal|aws] - Service Platform for Host build up
  -webconsole string
    	Klevr webconsole(server) address (URL or IP, Optional: Port) for connect (default "localhost:8080")
  -zone string
    	zone will be a [Dev/Stg/Prod] (default "dev-zone")
```

 * Using localhost: `./klevr -platform=baremetal -id=ralf -group="[COMPANY/TEAM]"`

 * Using seperated host:  `./klevr -platform=baremetal -id=ralf -webconsole=[WEBCONSOL_ADDR] `
