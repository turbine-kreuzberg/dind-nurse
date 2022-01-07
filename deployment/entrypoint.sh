#!/bin/sh

if [ "$DEBUG" = "true" ]; then
  echo `which dind-nurse` | dlv --listen=:$DEBUG_PORT --headless=true --api-version=2 --accept-multiclient exec `which dind-nurse` $@
else
  echo `which dind-nurse` | entr -nr `which dind-nurse` $@
fi