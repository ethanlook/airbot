#!/usr/bin/env bash

cd `dirname $0`
set -eux

exec /bins/viam-carto --appimage-extract-and-run $@
