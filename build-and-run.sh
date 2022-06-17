#!/bin/bash

docker build -t arduino-help-bot .
docker run --rm -it -v $(pwd)/_docs_example:/docs -e BOT_TOKEN="$1" -e BOT_ADMIN_ROLES="$2"  arduino-help-bot
