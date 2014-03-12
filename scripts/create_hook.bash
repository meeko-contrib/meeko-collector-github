#!/bin/bash

error() {
	echo "$1" 1>&2
	exit 1
}

[ "$#" != 1 ] && error "Usage: create_hook.bash <service>"
[ -z "$GH_OWNER" ] && error "GH_OWNER is not set"
[ -z "$GH_REPO" ] && error "GH_REPO is not set"

curl -POST \
     -d "@$1.json" \
     -H "Content-Type: application/json" \
     -u "$GH_OWNER" \
     "https://api.github.com/repos/$GH_OWNER/$GH_REPO/hooks"
