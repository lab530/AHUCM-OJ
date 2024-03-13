#!/usr/bin/env bash

# set -x

pip install -r ./requirements.txt
python specify.py

if [ -z "$WSL_DISTRO_NAME" ]; then
    cd ../docker && docker-compose up
else
    cd ../docker && sudo docker-compose up
fi

