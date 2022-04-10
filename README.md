# plugin-wireguard


*NOTE* dev mode

this will be integrated with this conatiner: https://github.com/spr-networks/super/tree/main/wireguard

testing the setup:

```sh
docker build -t plugin-wireguard --build-arg TARGETARCH=amd64 .
docker run -v $PWD/../state/wireguard:/state/api -v $PWD/../configs:/configs plugin-wireguard
```

get a client config:
```sh
export SOCK=$PWD/../state/wireguard/wireguard_plugin
sudo chmod a+w $SOCK
curl --unix-socket $SOCK http://localhost/config | jq .
```
