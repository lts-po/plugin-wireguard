# plugin-wireguard

this will be integrated with this conatiner: https://github.com/spr-networks/super/tree/main/wireguard

testing the setup:

```sh
docker build -t plugin-wireguard --build-arg TARGETARCH=amd64 .
docker run -v $PWD/../state/wireguard:/state/api -v $PWD/../configs:/configs plugin-wireguard
```

verify plugin is working:
```sh
export SOCK=$PWD/../state/wireguard/wireguard_plugin
sudo chmod a+w $SOCK
curl -s --unix-socket $SOCK http://localhost/peers
```

get a client config:
```sh
KEY=$(wg genkey)
PUBKEY=$(echo $KEY | wg pubkey)
curl -s --unix-socket $SOCK http://localhost/peer -X PUT --data "{\"PublicKey\": \"${PUBKEY}\"}" | jq . > wg.conf
cat wg.conf | sed "s/<PRIVATE KEY>/$KEY/g" | tee wg.conf
```

if no PublicKey is specifed one will be generated:
```sh
curl -s --unix-socket $SOCK http://localhost/peer -X PUT --data "{}"
```
