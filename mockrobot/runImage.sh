#!/usr/bin/env bash
docker run --rm -it -v "$PWD":/airbot -w /airbot mockrobot "$@"
