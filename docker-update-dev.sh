#!/usr/bin/env zsh

local IMG=code.benallen.dev/benallen/rsvp:dev

docker build -t $IMG .
docker login code.benallen.dev
docker push $IMG
