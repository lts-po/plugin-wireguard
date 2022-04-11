#!/bin/bash
. /configs/base/config.sh
# haxxy way to put the pubkey in env for go instead of parsing the conf
export LANIP=$LANIP
#export WIREGUARD_PORT=$WIREGUARD_PORT
/wireguard_plugin
