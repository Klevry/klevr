## How to run docker of Klevr agent
* ex)
```
sudo docker run -d -p 18800:18800 -e K_API_KEY="1231" -e K_PLATFORM="1231" -e K_MANAGER_URL="192.168.2.100:8090" -e K_ZONE_ID="13123" --name=klevr_agent  klevry/klevr-agent:latest
```
