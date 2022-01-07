#!/bin/sh

echo `which dind-nurse` | dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec `which dind-nurse` $@
