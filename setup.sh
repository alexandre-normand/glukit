#!/bin/sh

if [ ! -d "github.com/alexandre-normand/glukit/lib/github.com/grd/stat" ]; then
	mkdir -p "github.com/alexandre-normand/glukit/lib/github.com/grd"
	git clone git@github.com:grd/stat.git lib/github.com/grd/stat
fi

mkdir -p code.google.com/p
cd code.google.com/p
ln -fs ../../lib/goauth2
mkdir -p google-api-go-client/googleapi/transport
cd google-api-go-client/googleapi
curl -O https://google-api-go-client.googlecode.com/hg/googleapi/googleapi.go
cd transport 
curl -O https://google-api-go-client.googlecode.com/hg/googleapi/transport/apikey.go
