#!/bin/bash

docker build -t waterfall-engine .

docker run \
    -v ./data:/tmp/data \
    -e BD_APP_STORE_INMEMORY="true" \
    -e BD_APP_ENGINE_PREFERREDRETURN_HURDLEPERCENTAGE="0.08" \
    -e BD_APP_ENGINE_CATCHUP_CATCHUPPERCENTAGE="1.0" \
    -e BD_APP_ENGINE_CATCHUP_CARRIEDINTERESTPERCENTAGE="0.2" \
    -e BD_APP_ENGINE_FINALSPLIT_LPPERCENTAGE="0.8" \
    -e BD_APP_ENGINE_FINALSPLIT_GPPERCENTAGE="0.2" \
    -p 8080:8080 \
    waterfall-engine serve --data /tmp/data