#!/bin/sh

if [ ! -d "lib/github.com/grd/stat" ]; then
	git clone git@github.com:grd/stat.git lib/github.com/grd/stat
fi

mkdir -p code.google.com/p
cd code.google.com/p
ln -fs ../../lib/goauth2
mkdir -p google-api-go-client
cd google-api-go-client
ln -fs ../../../../google-api-go-client/googleapi
