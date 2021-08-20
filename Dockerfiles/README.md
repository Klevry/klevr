# Dockerfiles of Docker image build for Klevr eco. system
 * ```sudo docker run -d -p 18800:18800 -v /tmp/status:/info klevry/beacon:latest ```
 * ```sudo docker run -d klevry/libvirt:latest```
 * ```sudo docker run -e API_SERVER=http://xxx.xxx.xxx.xxx:8500 klevry/webconsole:latest```
 * ```sudo docker -e K_API_KEY="1231" -e K_PLATFORM="1231" -e K_MANAGER_URL="http://192.168.2.100:8090" -e K_ZONE_ID="13123"```
