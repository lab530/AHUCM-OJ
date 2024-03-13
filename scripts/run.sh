#!/usr/bin/env bash

# set -x

pip install -r ./requirements.txt
python specify.py

cd ../docker && docker-compose up

