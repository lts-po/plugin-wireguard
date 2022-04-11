#!/bin/bash
. /configs/base/config.sh
# haxxy way to put the pubkey in env for go instead of parsing the conf
export LANIP=$LANIP
export WIREGUARD_NETWORK=$WIREGUARD_NETWORK
export WIREGUARD_PORT=$WIREGUARD_PORT
WIREGUARD_PRIVKEY=$(cat /configs/wireguard/wg0.conf|grep -m1 ^PrivateKey | sed 's/\s//g'|sed 's/PrivateKey=//g')
export WIREGUARD_PUBKEY=$(echo "${WIREGUARD_PRIVKEY}" | wg pubkey)
/wireguard_plugin
