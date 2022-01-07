#!/bin/sh

echo `which dind-nurse` | entr -nr `which dind-nurse` $@
