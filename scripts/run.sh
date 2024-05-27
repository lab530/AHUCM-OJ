#!/usr/bin/env bash
# set -x
cd "$(dirname "$0")"

pip install -r ./requirements.txt
python specify.py

if [ ! -z "$WSL_DISTRO_NAME" ]; then
	(cd ../docker && docker-compose up) | tee -a docker-compose.log &
else
    (cd ../docker && sudo docker-compose up) | tee -a docker-compose.log &
fi

# (cd ../core && RUST_LOG=debug cargo run) | tee -a core.log &
#
# (cd ../backend && go mod tidy && go run main.go) | tee -a backend.log &
#
# (cd ../frontend && npm install && npm run dev) | tee -a frontend.log &

wait
